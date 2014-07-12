// latest-tag: Time-stamp: <[user_model.go] Elivoa @ Saturday, 2014-06-21 21:25:11>
package model

import (
	"strings"
	"time"
)

type User struct {
	Id         int64
	Name       string
	Position   string
	Username   string
	Password   string // no use?
	HashPass   string // non-persist field.
	Gender     string
	QQ         string // 车模项目中没用
	Mobile     string
	Mobile2    string // !!车模项目要求3个电话。
	Phone      string
	Country    string
	City       string
	Address    string // 车模项目添加。可作为公用项。
	Role       string
	Note       string
	CreateTime time.Time // 试试看，是否可作为入职时间。
	UpdateTime time.Time
}

type UserToken struct {
	Username string
	Password string
	Roles    []string // roles
	// TODO: timeout?
}

func (u *User) ToUserToken() *UserToken {
	userToken := &UserToken{
		Username: u.Username,
		Password: u.Password,
	}
	rawRoles := strings.Split(u.Role, ",")
	userToken.Roles = make([]string, 0)
	for _, r := range rawRoles {
		userToken.Roles = append(userToken.Roles, strings.ToLower(strings.TrimSpace(r)))
	}
	return userToken
}
