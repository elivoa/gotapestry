package accountdao

import (
	"database/sql"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"syd/model"
)

// create a new entity.
var aclEm = &db.Entity{
	Table:        "account_changelog",
	PK:           "id",
	Fields:       []string{"id", "customer_id", "delta", "account", "type", "related_order_tn", "reason", "time"},
	CreateFields: []string{"customer_id", "delta", "account", "type", "related_order_tn", "reason"},
	UpdateFields: []string{}, // update not needed
}

func init() {
	db.RegisterEntity("account_changelog", aclEm)
}

// CreateAccountChangeLog creates Account chagnelog
func CreateAccountChangeLog(acl *model.AccountChangeLog) (*model.AccountChangeLog, error) {
	res, err := aclEm.Insert().Exec(
		acl.CustomerId, acl.Delta, acl.Account, acl.Type, acl.RelatedOrderTN, acl.Reason,
	)
	if err != nil {
		return nil, err
	}
	aclId, err := res.LastInsertId()
	acl.Id = aclId
	return acl, nil
}

// DeleteAccountChagneLog deletes changelog by id.
func DeleteAccountChangeLog(id int) (int64, error) {
	res, err := aclEm.Delete().Exec(id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func AllAccountChangeLogs() ([]*model.AccountChangeLog, error) {
	return nil, nil
}

// list all accounts by id.
func ListAccountChangeLogsByCustomerId(customerId int) ([]*model.AccountChangeLog, error) {
	return list_account_changelog(aclEm.Select().Where("customer_id", customerId).OrderBy("id desc").Limit(20))
}

// list_incoming is an common function that accept a query and query a list of result, and error.
func list_account_changelog(query *db.QueryParser) ([]*model.AccountChangeLog, error) {
	changeLogs := make([]*model.AccountChangeLog, 0)
	err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			p := new(model.AccountChangeLog)
			err := rows.Scan(
				&p.Id, &p.CustomerId, &p.Delta, &p.Account, &p.Type,
				&p.RelatedOrderTN, &p.Reason, &p.Time,
			)
			changeLogs = append(changeLogs, p)
			return true, err
		},
	)
	if err != nil {
		return nil, err
	}
	return changeLogs, nil
}
