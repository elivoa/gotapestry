/*
 SQL Helper is a helper method in total filtering.

 Usage Examples:
  em.Select().Where("id", 5).QueryOne()
  em.Select("id","name").Where("type", "person").Query()
  ...
  em.Update().Where("id", 5).Exec(name, class, ...)
  em.Update().Exec(name, class, ..., id)
  em.Update("time").Where("id", 5, "person", 6).Exec(time)

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
	UpdateFields []string // fields used in create things.
}

// TODO Cache queryParser here.
func (e *Entity) Create(queryName string) *QueryParser {
	parser := &QueryParser{
		e: e,
	}
	return parser
}

func (e *Entity) Select(fields ...string) *QueryParser {
	return e.createQueryParser("select", fields...)
}

func (e *Entity) Insert(fields ...string) *QueryParser {
	return e.createQueryParser("insert", fields...)
}

func (e *Entity) Update(fields ...string) *QueryParser {
	return e.createQueryParser("update", fields...)
}

func (e *Entity) Delete() *QueryParser {
	return e.createQueryParser("delete")
}

func (e *Entity) createQueryParser(operation string, fields ...string) *QueryParser {
	parser := &QueryParser{
		e:         e,
		operation: operation,
		fields:    fields,
	}
	if nil != fields && len(fields) > 0 {
		parser.useCustomerFields = true
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
	e          *Entity
	operation  string       // select, insert, update, insertorupdate, delete
	fields     []string     // selected fields
	conditions []*condition // where 'id' = 1
	limit      *gxl.Int     // limit 4
	n          *gxl.Int     // limit 'limit','n'
	orderby    string

	prepared          bool
	useCustomerFields bool

	sql    string // generated sql
	values []interface{}
}

type condition struct {
	field  string
	values []interface{} // values, if only 1 value, use values[0]
	op     string        // and, or, andx, orx,
}

func (p *QueryParser) Fields(fields ...string) *QueryParser {
	p.useCustomerFields = true
	p.fields = fields
	return p
}

// support map[string]interface{}
// support interface{}...
// func (p *QueryParser) Where__the_old_where(conditions ...interface{}) *QueryParser {
// 	if len(conditions)%2 != 0 {
// 		panic("Wrong numnber of parameters!")
// 	}
// 	p.where = conditions
// 	return p
// }

func (p *QueryParser) Where(conditions ...interface{}) *QueryParser {
	p.conditions = []*condition{}
	if len(conditions) == 0 {
	} else if len(conditions) == 2 {
		p.conditions = append(p.conditions, &condition{
			field:  conditions[0].(string),
			values: []interface{}{conditions[1]},
			op:     "and",
		})
	} else {
		panic("Where clouse only accept 0 or 2 parameters.")
	}
	// TODO
	return p
}

func (p *QueryParser) And(field string, values ...interface{}) *QueryParser {
	p.conditions = append(p.conditions, &condition{field: field, values: values, op: "and"})
	return p
}

func (p *QueryParser) Or(field string, values ...interface{}) *QueryParser {
	p.conditions = append(p.conditions, &condition{field: field, values: values, op: "or"})
	return p
}

func (p *QueryParser) OrderBy(orderby string) *QueryParser {
	p.orderby = orderby
	return p
}

func (p *QueryParser) Limit(limit ...int) *QueryParser {
	if len(limit) >= 1 {
		p.limit = gxl.NewInt(limit[0])
	}
	if len(limit) >= 2 {
		p.n = gxl.NewInt(limit[1])
	}
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
		sql.WriteString(" FROM `")
		sql.WriteString(e.Table)
		sql.WriteString("`")

		// add where condition, default only support and
		if p.conditions != nil && len(p.conditions) > 0 {
			sql.WriteString(" WHERE ")
			p.values = appendWhereClouse(&sql, p.conditions...)
		}

		if p.orderby != "" {
			sql.WriteString(" order by ")
			sql.WriteString(p.orderby)
		}

		if p.limit != nil {
			sql.WriteString(" limit ")
			sql.WriteString(p.limit.String())
			if p.n != nil {
				sql.WriteString(",")
				sql.WriteString(p.n.String())
			}
		}
	case "insert":
		// em.Insert().Exec(name, class, ...)
		sql.WriteString("insert into `")
		sql.WriteString(e.Table)
		sql.WriteString("` (")

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

	case "update":
		// em.Update().Where("id", 5).Exec(name, class, ...)
		// em.Update().Exec(name, class, ..., id)
		sql.WriteString("update `")
		sql.WriteString(e.Table)
		sql.WriteString("` set ")

		fields := e.UpdateFields
		if p.useCustomerFields {
			fields = p.fields
		}
		for i := 0; i < len(fields); i++ {
			if i > 0 {
				sql.WriteString(",")
			}
			sql.WriteString(fmt.Sprintf("`%v`=?", fields[i]))
		}

		// where
		sql.WriteString(" WHERE ")
		if p.conditions == nil || len(p.conditions) == 0 {
			sql.WriteString(fmt.Sprintf(" `%v` = ?", e.PK))
		} else {
			p.values = appendWhereClouse(&sql, p.conditions...)
		}

	case "delete":
		// em.Delete().Where("id", 5).Exec()
		sql.WriteString("delete from `")
		sql.WriteString(e.Table)
		sql.WriteString("`")

		// where
		sql.WriteString(" WHERE ")
		if p.conditions == nil || len(p.conditions) == 0 {
			sql.WriteString(fmt.Sprintf(" `%v` = ?", e.PK))
		} else {
			p.values = appendWhereClouse(&sql, p.conditions...)
			// // TODO ... to be condinued....
			// for i := 0; i < len(p.conditions); i = i + 2 {
			// 	k, v := p.where[i].(string), p.where[i+1]
			// 	sql.WriteString(fmt.Sprintf(" `%v` = ?", k))
			// 	p.values = append(p.values, v)
			// 	if i < len(p.where)-3 {
			// 		sql.WriteString(" and ")
			// 	}
			// }
		}

	}
	p.sql = sql.String()
	p.prepared = true
	return p
}

// not
// param: use these value parameters to replace default value.
func (p *QueryParser) QueryOne(receiver func(*sql.Row) error) error {
	// query one will throw exceptions, so use query instead
	// TODO add limit support to QueryBuilder

	p.Prepare()

	// TODO use values to replace default one.
	debuglog("", "---------------------")
	debuglog("DB", "%v \"%v\" with parameters %v", p.operation, p.sql, p.values)

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
	if row != nil {
		err = receiver(row) // callbacks to receive values.
		if Err(err) {
			return err
		}
	}

	return nil
}

// query multi-results
func (p *QueryParser) Query(receiver func(*sql.Rows) (bool, error)) error {
	p.Prepare()

	// TODO use values to replace default one.
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(p.sql)
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

// exec command insert, update, delete
func (p *QueryParser) Exec(values ...interface{}) (sql.Result, error) {
	p.Prepare()

	debuglog("Exec", "\"%v\"", p.sql)

	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = Connect(); err != nil {
		return nil, err
	}
	defer conn.Close()

	if stmt, err = conn.Prepare(p.sql); err != nil {
		return nil, err
	}
	defer stmt.Close()

	// execute
	v := []interface{}{}
	v = append(v, values...)
	// for update command, use values as where condition.
	if p.values != nil && len(p.values) > 0 {
		v = append(v, p.values...)
	}

	debuglog("Exec", "with parameters %v", v)

	res, err := stmt.Exec(v...)
	if Err(err) {
		return nil, err
	}
	return res, nil
}

// ________________________________________________________________________________
var logEnabled = false

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

func appendWhereClouse(sql *bytes.Buffer, conditions ...*condition) []interface{} {
	values := []interface{}{}
	thefirst := true
	sql.WriteString(" ")
	for _, con := range conditions {
		lenvalue := len(con.values)
		switch con.op {
		case "and", "or":
			if !thefirst {
				if lenvalue > 1 {
					sql.WriteString(" and ")
				} else {
					sql.WriteString(" ")
					sql.WriteString(con.op)
					sql.WriteString(" ")
				}
			}
			if lenvalue == 1 {
				sql.WriteString(fmt.Sprintf("`%v`=?", con.field))
			} else if lenvalue > 1 {
				sql.WriteString("(")
				for idx, _ := range con.values {
					sql.WriteString(fmt.Sprintf("`%v`=?", con.field))
					if idx < lenvalue-1 {
						sql.WriteString(" ")
						sql.WriteString(con.op)
						sql.WriteString(" ")
					}
				}
				sql.WriteString(")")
			}
			values = append(values, con.values...)
		}
		thefirst = false
	}
	sql.WriteString(" ")
	return values
}
