package useractiondao

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"syd/model"
	"time"
)

// user action dao
var core_fields = []string{"user_id", "action", "context", "create_time"}
var historyEm = &db.Entity{
	Table:        "user_action",
	PK:           "id",
	Fields:       append([]string{"id"}, core_fields...),
	CreateFields: core_fields,
	UpdateFields: core_fields,
}

func init() {
	db.RegisterEntity("syd/user_action", historyEm)
}

func EntityManager() *db.Entity {
	return historyEm
}

//  log user action
func LogUserAction(userId int64, action model.ActionType, contexts ...interface{}) error {
	// contexts ==> string
	var buf bytes.Buffer
	if nil != contexts && len(contexts) > 0 {
		for idx, context := range contexts {
			if idx > 0 {
				buf.WriteRune(',')
			}
			str := fmt.Sprint(context)
			if strings.Contains(str, ",") {
				buf.WriteString("\"")
				buf.WriteString(str)
				buf.WriteString("\"")
			} else {
				buf.WriteString(str)
			}
		}
	}
	// write to db.
	_, err := historyEm.Insert().Exec(userId, action, buf.String(), time.Now())
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

func ListUserAction(parser *db.QueryParser) ([]*model.UserAction, error) {
	parser.SetEntity(historyEm) // set entity manager into query parser.
	parser.Reset()              // to prevent if parser is used before. TODO:Is this necessary?
	// append default behavore.
	parser.DefaultOrderBy("create_time", db.DESC)
	parser.DefaultLimit(0, config.LIST_PAGE_SIZE)
	parser.Select()
	return _listUserAction(parser)
}

// the last part, read the list from rows
func _listUserAction(query *db.QueryParser) ([]*model.UserAction, error) {
	models := make([]*model.UserAction, 0)
	if err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			m := &model.UserAction{}
			err := rows.Scan(
				&m.Id, &m.UserId, &m.Action, &m.Context, &m.CreateTime,
			)
			models = append(models, m)
			return true, err
		},
	); err != nil {
		return nil, err
	}
	return models, nil
}
