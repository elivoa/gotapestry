package inventorytrackdao

import (
	"database/sql"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"syd/base/inventory"
	"syd/model"
)

var core_fields = []string{
	"product_id", "color", "size", "stock_change_to", "old_stock", "delta",
	"user_id", "reason", "context",
}
var em = &db.Entity{
	Table:        "inventory_track",
	PK:           "id",
	Fields:       append(append([]string{"id"}, core_fields...), "time"),
	CreateFields: core_fields,
	UpdateFields: core_fields, // update not needed
}

func init() {
	db.RegisterEntity("inventory_track", em)
}

func EntityManager() *db.Entity {
	return em
}

func CreateInventoryTrack(m *model.InventoryTrackItem) (*model.InventoryTrackItem, error) {
	res, err := em.Insert().Exec(
		m.ProductId, m.Color, m.Size, m.StockChagneTo, m.OldStock, m.Delta,
		m.UserId, m.Reason, m.Context,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	m.Id = id
	return m, nil
}

func DeleteInventoryTrack(id int64) (int64, error) {
	return em.DeleteByPK(id)
}

// func AllInventoryTracks() ([]*model.InventoryTrackItem, error) {
// 	return nil, nil
// }

// list all accounts by id.
func ListInventoryTracks(productId int64) ([]*model.InventoryTrackItem, error) {
	return _list(em.Select().
		Where(inventory.F_Track_ProductId, productId).
		OrderBy("id", db.DESC).Limit(300),
	)
}

func List(parser *db.QueryParser) ([]*model.InventoryTrackItem, error) {
	return _list(parser)
}

// the last part, read the list from rows
func _list(query *db.QueryParser) ([]*model.InventoryTrackItem, error) {
	models := make([]*model.InventoryTrackItem, 0)
	if err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			m := &model.InventoryTrackItem{}
			err := rows.Scan(
				&m.Id, &m.ProductId, &m.Color, &m.Size, &m.StockChagneTo, &m.OldStock, &m.Delta,
				&m.UserId, &m.UserId, &m.Reason, &m.Context, &m.Time,
			)
			models = append(models, m)
			return true, err
		},
	); err != nil {
		return nil, err
	}
	return models, nil
}
