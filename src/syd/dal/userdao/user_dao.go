package userdao

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"syd/model"
)

// create a new entity.
var em = &db.Entity{
	Table:        "user",
	PK:           "id",
	Fields:       []string{"id", "username", "password", "gender", "qq", "mobile", "city", "role", "create_time", "update_time"},
	CreateFields: []string{"username", "password", "gender", "qq", "mobile", "city", "role", "update_time"},
	UpdateFields: []string{"username", "password", "gender", "qq", "mobile", "city", "role"},
}

func init() {
	db.RegisterEntity("user", em)
}

func CreateUser(user *model.User) (*model.User, error) {
	// password???
	res, err := em.Insert().Exec(
		user.Username, user.Password, user.Gender, user.QQ, user.Mobile, user.City, user.Role, time.Now(),
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	user.Id = id
	return user, nil
}

func DeleteUser(id int64) (int64, error) {
	res, err := em.Delete().Exec(id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func GetUserById(id int64) (*model.User, error) {
	p := new(model.User)
	err := em.Select().Where("id", id).Query(
		func(rows *sql.Rows) (bool, error) {
			return false, rows.Scan(
				&p.Id, &p.Username, &p.Password, &p.Gender, &p.QQ, &p.Mobile, &p.City, &p.Role,
				&p.CreateTime, &p.UpdateTime,
			)
		},
	)
	if err != nil {
		return nil, err
	}
	if p.Id > 0 {
		return p, nil
	}
	return nil, errors.New("Person not found!")
}

// TODO: password should be encripted.
func GetUserWithCredential(username string, password string) *model.User {
	fmt.Println("LoginService :> Get user with username/password pair : ", username, password)

	p := new(model.User)
	err := em.Select().Where("username", username).And("password", password).Query(
		func(rows *sql.Rows) (bool, error) {
			return false, rows.Scan(
				&p.Id, &p.Username, &p.Password, &p.Gender, &p.QQ, &p.Mobile, &p.City, &p.Role,
				&p.CreateTime, &p.UpdateTime,
			)
		},
	)
	if err != nil {
		panic(err)
	}
	if p.Id > 0 { // return p.Id if true
		return p
	}
	return nil
}

// TODO: password should be encripted.
func VerifyLogin(username string, password string) bool {
	user := GetUserWithCredential(username, password)
	return user.Id > 0 // return p.Id if true.
}

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
