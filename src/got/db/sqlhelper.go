/*
 SQL Helper is a helper method in total filtering.
*/
package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"gxl"
	"strings"
)

// ________________________________________

// Entities cache.
// TODO: thread safe? need lock?
var entities map[string]*Entity
var queryparserCache map[string]*QueryParser

func init() {
	entities = make(map[string]*Entity, 10)
	queryparserCache = make(map[string]*QueryParser)
}

func RegisterEntity(name string, entity *Entity) {
	if _, ok := entities[name]; ok {
		panic("DB: Register duplicated entities.")
	}
	entities[name] = entity
}

// --------------------------------------------------------------------------------

// ________________________________________________________________________________
// DAO Helper
type Entity struct {
	Table        string   // table name
	PK           string   // primary key field name
	Fields       []string // field names
	CreateFields []string // fields used in create things.
}

// TODO Cache queryParser here.
func (e *Entity) Create(queryName string) *QueryParser {
	parser := &QueryParser{
		e: e,
	}
	return parser
}

func (e *Entity) Select(fields ...string) *QueryParser {
	parser := &QueryParser{
		e:                 e,
		operation:         "select",
		useCustomerFields: true,
		fields:            fields,
	}
	if nil == fields || len(fields) == 0 {
		parser.useCustomerFields = false
	}
	return parser
}

func (e *Entity) Insert(fields ...string) *QueryParser {
	parser := &QueryParser{
		e:                 e,
		operation:         "insert",
		useCustomerFields: true,
		fields:            fields,
	}
	if nil == fields || len(fields) == 0 {
		parser.useCustomerFields = false
	}
	return parser
}

// TODO not used
func (e *Entity) NamedQuery(name string, createfunc func() *QueryParser) *QueryParser {
	cached, ok := queryparserCache[name]
	if !ok {
		cached = createfunc()
		queryparserCache[name] = cached
	}
	return cached

}

// ________________________________________________________________________________
// Query parser
//
type QueryParser struct {
	e         *Entity
	operation string        // select, insert, update, insertorupdate, delete
	fields    []string      // selected fields
	where     []interface{} // where 'id' = 1
	limit     *gxl.Int      // limit 4
	n         *gxl.Int      // limit 'limit','n'

	prepared          bool
	useCustomerFields bool

	sql    string // generated sql
	values []interface{}
}

func (p *QueryParser) Fields(fields ...string) *QueryParser {
	p.useCustomerFields = true
	p.fields = fields
	return p
}

// support map[string]interface{}
// support interface{}...
func (p *QueryParser) Where(conditions ...interface{}) *QueryParser {
	if len(conditions)%2 != 0 {
		panic("Wrong numnber of parameters!")
	}
	p.where = conditions
	return p
}

// pin sql and cache them
func (p *QueryParser) Prepare() *QueryParser {
	if p.prepared {
		return p
	}

	e := p.e
	var sql bytes.Buffer
	switch p.operation {
	case "select":
		sql.WriteString("SELECT ")
		if p.useCustomerFields {
			sql.WriteString(fieldString(p.fields))
		} else {
			sql.WriteString(fieldString(e.Fields))
		}

		// from
		sql.WriteString(" FROM ")
		sql.WriteString(e.Table)

		// add where condition, default only support and
		if p.where != nil && len(p.where) > 0 {
			sql.WriteString(" WHERE ")
			for i := 0; i < len(p.where); i = i + 2 {
				k, v := p.where[i].(string), p.where[i+1]
				sql.WriteString(fmt.Sprintf(" `%v` = ?", k))
				p.values = append(p.values, v)
				if i < len(p.where)-3 {
					sql.WriteString(" and ")
				}
			}
		}
		// TODO order by
		// TODO limit

	case "insert":
		sql.WriteString("insert into ")
		sql.WriteString(e.Table)
		sql.WriteString(" (")

		fields := e.CreateFields
		if p.useCustomerFields {
			fields = p.fields
		}
		sql.WriteString(fmt.Sprintf("`%v`", strings.Join(fields, "`,`")))
		sql.WriteString(" )")
		// values
		sql.WriteString(" values (")
		for i := 0; i < len(fields); i++ {
			if i > 0 {
				sql.WriteString(",")
			}
			sql.WriteString("?")
		}
		sql.WriteString(" )")

	}
	p.sql = sql.String()
	p.prepared = true
	return p
}

// param: use these value parameters to replace default value.
func (p *QueryParser) QueryOne(receiver func(*sql.Row) error) error {
	p.Prepare()

	// TODO use values to replace default one.
	debuglog("Exec", "QueryOne \"%v\" with parameters %v", p.sql, p.values)

	// 1. get connection
	conn, err := Connect()
	if Err(err) {
		return err
	}
	defer conn.Close()

	// 2. prepare sql
	stmt, err := conn.Prepare(p.sql)
	if Err(err) {
		return err
	}
	defer stmt.Close()

	// 3. execute
	row := stmt.QueryRow(p.values...)
	err = receiver(row) // callbacks to receive values.
	if Err(err) {
		return err
	}
	return nil
}

// query multi-results
func (p *QueryParser) Query(receiver func(*sql.Rows) (bool, error)) error {
	p.Prepare()

	// TODO use values to replace default one.
	debuglog("Exec", "Query \"%v\" with parameters %v", p.sql, p.values)

	// 1. get connection
	conn, err := Connect()
	if Err(err) {
		return err
	}
	defer conn.Close()

	// 2. prepare sql
	stmt, err := conn.Prepare(p.sql)
	if Err(err) {
		return err
	}
	defer stmt.Close()

	// 3. execute
	rows, err := stmt.Query(p.values...)
	if err != nil {
		return err
	}
	for rows.Next() {
		goon, err := receiver(rows) // callbacks to receive values.
		if Err(err) {
			return err
		}
		if !goon {
			break
		}
	}
	return nil
}

func (p *QueryParser) Exec(values ...interface{}) (sql.Result, error) {
	p.Prepare()

	debuglog("Exec", "Insert: \"%v\" with parameters %v", p.sql, p.values)

	conn, err := Connect()
	if Err(err) {
		return nil, err
	}
	defer conn.Close()

	// 1. prepare ql
	stmt, err := conn.Prepare(p.sql)
	defer stmt.Close()
	if Err(err) {
		return nil, err
	}

	// 3. execute
	res, err := stmt.Exec(values...)
	if Err(err) {
		return nil, err
	}
	return res, nil
}

// ________________________________________________________________________________
var logEnabled = true

func debuglog(method string, format string, params ...interface{}) {
	if logEnabled {
		fmt.Printf("[DB.%v] %v\n",
			method,
			fmt.Sprintf(format, params...),
		)
	}
}

// helper methods
func fieldString(fields []string) string {
	if fields == nil || len(fields) == 0 {
		return "*"
	}
	return fmt.Sprintf("`%v`",
		strings.Join(fields, "`, `"),
	)
}
