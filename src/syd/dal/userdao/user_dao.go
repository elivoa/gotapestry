// latest-tag: [user_dao.go] Time-stamp: <[user_dao.go] Elivoa @ Friday, 2015-03-27 13:19:26>
package userdao

import (
	"database/sql"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"syd/model"
	"time"
)

var core_fields = []string{
	"username", "password", "gender", "qq", "mobile", "city", "role", "create_time",
}

var em = &db.Entity{
	Table:        "user",
	PK:           "id",
	Fields:       append(append([]string{"id"}, core_fields...), "update_time"),
	CreateFields: core_fields,
	UpdateFields: core_fields,
}

func init() {
	db.RegisterEntity("syd/user", em)
}

func EntityManager() *db.Entity {
	return em
}

func _one(query *db.QueryParser) (*model.User, error) {
	p := new(model.User)
	err := query.Query(
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
	return nil, nil
}

func _list(query *db.QueryParser) ([]*model.User, error) {
	models := make([]*model.User, 0)
	if err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			p := &model.User{}
			err := rows.Scan(
				&p.Id, &p.Username, &p.Password, &p.Gender, &p.QQ, &p.Mobile, &p.City, &p.Role,
				&p.CreateTime, &p.UpdateTime,
			)
			models = append(models, p)
			return true, err
		},
	); err != nil {
		return nil, err
	}
	return models, nil
}

func CreateUser(user *model.User) (*model.User, error) {
	// password???
	res, err := em.Insert().Exec(
		user.Username, user.Password, user.Gender, user.QQ, user.Mobile, user.City, user.Role, time.Now(),
		// TODO: change time.Now() to user.CreateTime, need to assign all create time outside.
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	user.Id = id
	return user, nil
}

func UpdateUser(user *model.User) (int64, error) {
	// update order
	res, err := em.Update().Exec(
		user.Username, user.Password, user.Gender, user.QQ, user.Mobile, user.City, user.Role,
		user.UpdateTime,
		user.Id, // condition
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func DeleteUser(id int64) (int64, error) {
	res, err := em.Delete().Exec(id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func GetUser(field string, value interface{}) (*model.User, error) {
	return _one(em.Select().Where(field, value))
}

// TODO: password should be encripted.
func GetUserWithCredential(username string, password string) (*model.User, error) {
	return _one(em.Select().Where("username", username).And("password", password))
}

func GetUserById(id int64) (*model.User, error) {
	return GetUser(em.PK, id)
}

// TODO: password should be encripted.
func VerifyLogin(username string, password string) bool {
	user, err := GetUserWithCredential(username, password)
	return user.Id > 0 && err == nil // return p.Id if true.
}

func ListUser() ([]*model.User, error) {
	var query *db.QueryParser
	parser := em.Select().Where()
	query = parser.OrderBy("id", db.DESC)
	return _list(query)
}

func ListUserByIdSet(ids ...int64) (map[int64]*model.User, error) {
	if nil == ids || len(ids) == 0 {
		return nil, nil
	}
	var query *db.QueryParser
	parser := em.Select().Where()
	query = parser.InInt64("id", ids...).OrderBy("id", db.DESC)

	users, err := _list(query)
	if err != nil {
		return nil, err
	}

	// users := make([]*model.User, 0)
	// if err := query.Query(
	// 	func(rows *sql.Rows) (bool, error) {
	// 		p := new(model.User)
	// 		err := rows.Scan(
	// 			&p.Id, &p.Name, &p.Position, &p.Username, &p.Password, &p.Gender, &p.QQ,
	// 			&p.Mobile, &p.Mobile2, &p.Phone, &p.Country, &p.City, &p.Address, &p.Store, &p.Role, &p.Note,
	// 			&p.CreateTime, &p.UpdateTime,
	// 		)
	// 		users = append(users, p)
	// 		return true, err
	// 	},
	// ); err != nil {
	// 	return nil, err
	// }

	// return the map
	var usermap = map[int64]*model.User{}
	for _, u := range users {
		usermap[u.Id] = u
	}
	return usermap, nil
}

// old things

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
