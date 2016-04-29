package inventorydao

import (
	"database/sql"
	"fmt"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"syd/model"
)

var logdebug = true

func init() {
	db.RegisterEntity("syd/inventory", em)
}

var core_fields = []string{"group_id", "product_id", "color", "size", "stock", "provider_id",
	"operator_id", "status", "type", "note", "send_time", "receive_time", "create_time"}

var em = &db.Entity{
	Table:        "inventory",
	PK:           "id",
	Fields:       append(append([]string{"id"}, core_fields...), "update_time"),
	CreateFields: core_fields,
	UpdateFields: core_fields,
}

func EntityManager() *db.Entity {
	return em
}

//
// Universal one and list private
//

func _one(query *db.QueryParser) (*model.Inventory, error) {
	m := new(model.Inventory)
	err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			return false, rows.Scan(
				&m.Id, &m.GroupId, &m.ProductId, &m.Color, &m.Size, &m.Stock, &m.ProviderId, &m.OperatorId,
				&m.Status, &m.Type, &m.Note, &m.SendTime, &m.ReceiveTime, &m.CreateTime, &m.UpdateTime,
			)
		},
	)
	if err != nil {
		return nil, err
	}
	if m.Id > 0 {
		return m, nil
	}
	return nil, nil
}

// the last part, read the list from rows
func _list(query *db.QueryParser) ([]*model.Inventory, error) {
	models := make([]*model.Inventory, 0)
	if err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			m := &model.Inventory{}
			err := rows.Scan(
				&m.Id, &m.GroupId, &m.ProductId, &m.Color, &m.Size, &m.Stock, &m.ProviderId, &m.OperatorId,
				&m.Status, &m.Type, &m.Note, &m.SendTime, &m.ReceiveTime, &m.CreateTime, &m.UpdateTime,
			)
			models = append(models, m)
			return true, err
		},
	); err != nil {
		return nil, err
	}

	return models, nil
}

//
// Universal Get and List Public
//

func Get(field string, value interface{}) (*model.Inventory, error) {
	return _one(em.Select().Where(field, value))
}

func List(parser *db.QueryParser) ([]*model.Inventory, error) {
	// var query *db.QueryParser
	parser.SetEntity(em) // set entity manager into query parser.
	parser.Reset()       // to prevent if parser is used before. TODO:Is this necessary?
	// append default behavore.
	parser.DefaultOrderBy("send_time", db.DESC)
	parser.DefaultLimit(0, config.LIST_PAGE_SIZE)
	parser.Select()
	return _list(parser)
}

func GetInventoryById(id int64) (*model.Inventory, error) {
	return _one(em.Select().Where(em.PK, id))
}

//
// Create
//

func Create(m *model.Inventory) (*model.Inventory, error) {
	res, err := em.Insert().Exec(
		m.GroupId, m.ProductId, m.Color, m.Size, m.Stock, m.ProviderId, m.OperatorId,
		m.Status, m.Type, m.Note, m.SendTime, m.ReceiveTime, m.CreateTime,
	)
	if err != nil {
		return nil, err
	}
	liid, err := res.LastInsertId()
	m.Id = liid
	return m, nil
}

func Update(m *model.Inventory) (int64, error) {
	res, err := em.Update().Exec(
		m.GroupId, m.ProductId, m.Color, m.Size, m.Stock, m.ProviderId, m.OperatorId,
		m.Status, m.Type, m.Note, m.SendTime, m.ReceiveTime, m.CreateTime,
		m.Id,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func Delete(id int64) (int64, error) {
	return em.DeleteByPK(id)
}

// 这是一个落后的原始的sql

func UpdateAllInventoryItems(ig *model.InventoryGroup) error {
	if nil == ig && ig.Id >= 0 {
		panic("InventoryId is nil!")
	}
	_sql := `update inventory set provider_id=?, operator_id=?, send_time=?, receive_time=? where group_id=?`
	parser := em.RawQuery(_sql)
	_, err := parser.Exec(ig.ProviderId, ig.OperatorId, ig.SendTime, ig.ReceiveTime, ig.Id)
	return err
}

// old things.
// func SearchInventoryInUseByPattern(pattern string) ([]*model.Inventory, error) {
// 	var conn *sql.DB
// 	var stmt *sql.Stmt
// 	var err error
// 	if conn, err = db.Connect(); err != nil {
// 		return nil, err
// 	}
// 	defer conn.Close()

// 	sql := "select i.id, i.product_id, i.serialno, i.store, i.status, i.note, p.name, p.type, p.property from inventory i left join product p on p.id=i.product_id where i.status=? and p.name like ? limit 100"
// 	if stmt, err = conn.Prepare(sql); err != nil {
// 		return nil, err
// 	}
// 	defer stmt.Close()

// 	rows, err := stmt.Query(model.InventoryStatus_InUse, pattern)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	models := []*model.Inventory{}
// 	for rows.Next() {
// 		m := &model.Inventory{}
// 		p := &model.Product{}
// 		err := rows.Scan(
// 			&m.Id, &m.ProductId, &m.SerialNo, &m.Store, &m.Status, &m.Note,
// 			&p.Name, &p.Type, &p.Property,
// 		)
// 		if err != nil {
// 			panic(err)
// 		}

// 		p.Id = int(m.ProductId)
// 		m.Product = p
// 		models = append(models, m)
// 	}
// 	return models, nil
// }
