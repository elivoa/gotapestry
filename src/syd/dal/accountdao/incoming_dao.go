package accountdao

import (
	"database/sql"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"syd/model"
	"time"
)

// create a new entity.
var em = &db.Entity{
	Table:        "incoming",
	PK:           "id",
	Fields:       []string{"id", "customer_id", "incoming", "time"},
	CreateFields: []string{"customer_id", "incoming", "time"},
	UpdateFields: []string{"customer_id", "incoming", "time"},
}

func init() {
	db.RegisterEntity("incoming", em)
}

// --------------------------------------------------------------------------------
// TODO:
//   Order by time
//   restrict in time range.
//

func MyIncoming() ([]*model.AccountIncoming, error) {
	return list_incoming(em.Select())
}

func IncomingHistory(customerId int) ([]*model.AccountIncoming, error) {
	return list_incoming(em.Select().Where("customer_id", customerId))
}

// list_incoming is an common function that accept a query and query a list of result, and error.
func list_incoming(query *db.QueryParser) ([]*model.AccountIncoming, error) {
	incomings := make([]*model.AccountIncoming, 0)
	err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			p := new(model.AccountIncoming)
			err := rows.Scan(&p.Id, &p.CustomeId, &p.Incoming, &p.Time)
			incomings = append(incomings, p)
			return true, err
		},
	)
	if err != nil {
		return nil, err
	}
	return incomings, nil
}

// CreateIncoming creates
func CreateIncoming(incoming *model.AccountIncoming) (*model.AccountIncoming, error) {
	res, err := em.Insert().Exec(
		incoming.CustomeId, incoming.Incoming, time.Now(),
	)
	if err != nil {
		return nil, err
	}
	liid, err := res.LastInsertId()
	incoming.Id = int(liid)
	return incoming, nil
}

func DeleteIncoming(id int) (int64, error) {
	res, err := em.Delete().Exec(id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
