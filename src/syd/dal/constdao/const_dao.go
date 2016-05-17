package constdao

import (
	"database/sql"
	"fmt"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"syd/model"
)

var logdebug = true

func init() {
	db.RegisterEntity("syd/const", em)
}

var core_fields = []string{"name", "key", "value", "doublevalue"}

var em = &db.Entity{
	Table:        "const",
	PK:           "id",
	Fields:       append(append([]string{"id"}, core_fields...), "time"),
	CreateFields: core_fields,
	UpdateFields: core_fields,
}

func EntityManager() *db.Entity {
	return em
}

// ________________________________________________________________________________
// Get product
//
//

// do not use any more
func GetIntValue(name, key string) (int64, error) {
	if ccc, err := GetOne(name, key); err != nil {
		return 0, err
	} else {
		return strconv.ParseInt(ccc.Value, 10, 64)
	}
}

// do not use any more
func Get2ndIntValue(name, key string) (int64, error) {
	if ccc, err := GetOne(name, key); err != nil {
		return 0, err
	} else {
		return int64(ccc.FloatValue), nil //  strconv.ParseInt(ccc.Value, 10, 64)
	}
}

// do not use any more
func Get2ndStringValue(name, key string) (string, error) {
	if ccc, err := GetOne(name, key); err != nil {
		return "", err
	} else {
		return fmt.Sprint(ccc.FloatValue), nil
	}
}

// do not use any more
func GetStringValue(name, key string) (string, error) {
	if ccc, err := Get(name, key); err != nil {
		return "", err
	} else {
		return ccc.Value, nil
	}
}

func GetById(id int64) (*model.Const, error) {
	return Get(em.PK, id)
}

func GetOne(name, key string) (*model.Const, error) {
	var query = em.Select().Where("name", name).And("key", key)
	if models, err := _list(query); err != nil {
		panic(err)
	} else {
		if nil != models && len(models) > 0 {
			if models[0] == nil {
				// return nil, errors.New("Const Not founc.")
				return nil, nil
			}
		}
		return models[0], nil
	}
}

// 搞什么鬼
func Get(field string, value interface{}) (*model.Const, error) {
	var query = em.Select().Where(field, value)
	if models, err := _list(query); err != nil {
		panic(err)
	} else {
		if nil != models && len(models) > 0 {
			if models[0] == nil {
				// return nil, errors.New("Const Not founc.")
				return nil, nil
			}
		}
		return models[0], nil
	}
}

func GetByNameKey(name, key string) (*model.Const, error) {
	var query = em.Select().Where("name", name).And("key", key)
	if models, err := _list(query); err != nil {
		panic(err)
	} else {
		if nil != models && len(models) > 0 {
			if models[0] != nil {
				return models[0], nil
			}
		}
		return nil, nil
	}
}

func Set(name string, key string, value interface{}, floatValue float64) error {
	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return err
	}
	defer conn.Close()

	sql := "insert into const(`name`, `key`, `value`, doublevalue) values(?,?,?,?) on duplicate key update `key`=?, `value`=?, doublevalue=?"
	if stmt, err = conn.Prepare(sql); err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, key, value, floatValue, key, value, floatValue)
	if err != nil {
		return err
	}
	return nil // res.RowsAffected()
}

func Update(name string, key string, value interface{}, floatValue float64, id int64) error {
	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return err
	}
	defer conn.Close()

	sql := "update const set `name`=?, `key`=?, `value`=?, doublevalue=? where id=?"
	if stmt, err = conn.Prepare(sql); err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, key, value, floatValue, id)
	if err != nil {
		return err
	}
	return nil // res.RowsAffected()
}

func DeleteById(id int64) (int64, error) {
	res, err := em.Delete().Exec(id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func Delete(name string, key string) (int64, error) {
	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return 0, err
	}
	defer conn.Close()

	sql := "delete from const where name=? and key=?"
	if stmt, err = conn.Prepare(sql); err != nil {
		return 0, err
	}
	defer stmt.Close()

	// execute
	res, err := stmt.Exec(name, key)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func DeleteByName(name string) (int64, error) {
	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return 0, err
	}
	defer conn.Close()

	sql := "delete from const where name=?"
	if stmt, err = conn.Prepare(sql); err != nil {
		return 0, err
	}
	defer stmt.Close()

	// execute
	res, err := stmt.Exec(name)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func GetList(name string) ([]*model.Const, error) {
	parser := em.NewQueryParser().Select().Where("name", name)
	// parser.DefaultOrderBy("create_time", db.DESC)
	// parser.DefaultLimit(0, config.LIST_PAGE_SIZE)
	return _list(parser)
}

// the last part, read the list from rows
func _list(query *db.QueryParser) ([]*model.Const, error) {
	models := make([]*model.Const, 0)
	if err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			m := &model.Const{}
			err := rows.Scan(&m.Id, &m.Name, &m.Key, &m.Value, &m.FloatValue, &m.Time)
			models = append(models, m)
			return true, err
		},
	); err != nil {
		return nil, err
	}
	return models, nil
}
