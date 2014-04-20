package userdao

import (
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
)

// user action history

var historyEm = &db.Entity{
	Table:        "user_action_history",
	PK:           "id",
	Fields:       []string{"id"},
	CreateFields: []string{"username", "password", "gender", "qq", "mobile", "city", "role"},
	UpdateFields: []string{"username", "password", "gender", "qq", "mobile", "city", "role"},
}

func init() {
	db.RegisterEntity("user_action_history", historyEm)
}

//  log user action
func LogUserAction(userId int64, action int) error {
	// password???
	_, err := historyEm.Insert().Exec()
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserAction(id int64) (int64, error) {
	res, err := historyEm.Delete().Exec(id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// func GetUserById(id int64) (*model.User, error) {
// 	p := new(model.User)
// 	err := historyEm.Select().Where("id", id).Query(
// 		func(rows *sql.Rows) (bool, error) {
// 			return false, rows.Scan(
// 				&p.Id, &p.Username, &p.Password, &p.Gender, &p.QQ, &p.Mobile, &p.City, &p.Role,
// 				&p.CreateTime, &p.UpdateTime,
// 			)
// 		},
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if p.Id > 0 {
// 		return p, nil
// 	}
// 	return nil, errors.New("Person not found!")
// }

// list all accounts by id.
// func ListAccountChangeLogsByCustomerId(customerId int) ([]*model.AccountChangeLog, error) {
// 	return list_account_changelog(aclEm.Select().Where("customer_id", customerId).OrderBy("id desc").Limit(20))
// }

// list_incoming is an common function that accept a query and query a list of result, and error.
// func list_account_changelog(query *db.QueryParser) ([]*model.AccountChangeLog, error) {
// 	changeLogs := make([]*model.AccountChangeLog, 0)
// 	err := query.Query(
// 		func(rows *sql.Rows) (bool, error) {
// 			p := new(model.AccountChangeLog)
// 			err := rows.Scan(
// 				&p.Id, &p.CustomerId, &p.Delta, &p.Account, &p.Type,
// 				&p.RelatedOrderTN, &p.Reason, &p.Time,
// 			)
// 			changeLogs = append(changeLogs, p)
// 			return true, err
// 		},
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return changeLogs, nil
// }
