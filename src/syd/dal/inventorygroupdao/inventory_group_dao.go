package inventorygroupdao

import (
	"database/sql"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"syd/model"
)

func init() {
	db.RegisterEntity("syd/inventorygroup", em)
}

var core_fields = []string{"status", "type", "note", "provider_id", "operator_id",
	"summary", "total_quantity", "send_time", "receive_time", "create_time"}

var em = &db.Entity{
	Table:        "inventory_group",
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

func _one(query *db.QueryParser) (*model.InventoryGroup, error) {
	m := new(model.InventoryGroup)
	err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			return false, rows.Scan(
				&m.Id, &m.Status, &m.Type, &m.Note, &m.ProviderId, &m.OperatorId,
				&m.Summary, &m.TotalQuantity, &m.SendTime, &m.ReceiveTime, &m.CreateTime, &m.UpdateTime,
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
func _list(query *db.QueryParser) ([]*model.InventoryGroup, error) {
	models := make([]*model.InventoryGroup, 0)
	if err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			m := &model.InventoryGroup{}
			err := rows.Scan(
				&m.Id, &m.Status, &m.Type, &m.Note, &m.ProviderId, &m.OperatorId,
				&m.Summary, &m.TotalQuantity, &m.SendTime, &m.ReceiveTime, &m.CreateTime, &m.UpdateTime,
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

func Get(field string, value interface{}) (*model.InventoryGroup, error) {
	return _one(em.Select().Where(field, value))
}

func List(parser *db.QueryParser) ([]*model.InventoryGroup, error) {
	// var query *db.QueryParser
	parser.SetEntity(em) // set entity manager into query parser.
	parser.Reset()       // to prevent if parser is used before. TODO:Is this necessary?
	// append default behavore.
	parser.DefaultOrderBy("send_time", db.DESC)
	parser.DefaultLimit(0, config.LIST_PAGE_SIZE)
	parser.Select()
	return _list(parser)
}

func GetInventoryGroupById(id int64) (*model.InventoryGroup, error) {
	return _one(em.Select().Where(em.PK, id))
}

//
// Create
//

func Create(m *model.InventoryGroup) (*model.InventoryGroup, error) {
	res, err := em.Insert().Exec(
		m.Status, m.Type, m.Note, m.ProviderId, m.OperatorId, m.Summary, m.TotalQuantity,
		m.SendTime, m.ReceiveTime, m.CreateTime,
	)
	if err != nil {
		return nil, err
	}
	liid, err := res.LastInsertId()
	m.Id = liid
	return m, nil
}

func Update(m *model.InventoryGroup) (int64, error) {
	res, err := em.Update().Exec(
		m.Status, m.Type, m.Note, m.ProviderId, m.OperatorId, m.Summary, m.TotalQuantity,
		m.SendTime, m.ReceiveTime, m.CreateTime,
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
