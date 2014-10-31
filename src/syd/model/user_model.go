// latest-tag: Time-stamp: <[user_model.go] Elivoa @ Friday, 2014-10-31 13:36:59>
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

// UserToken
type UserToken struct {
	Id       int64
	Name     string
	Username string
	Password string
	Roles    []string // roles, do not support multi roles now.
	Store    string   // no use in syd
	// TODO: timeout?
}

func (t *UserToken) Role() string {
	if nil == t.Roles || len(t.Roles) == 0 {
		return ""
	}
	return t.Roles[0]
}

func (s *UserToken) HasRole(role string) bool {
	lowerRole := strings.ToLower(role)

	if s == nil || s.Roles == nil {
		return false
	}
	for _, r := range s.Roles {
		if r == lowerRole {
			return true
		}
	}
	return false
}

func (u *User) ToUserToken() *UserToken {
	userToken := &UserToken{
		Id:       u.Id,
		Name:     u.Name,
		Username: u.Username,
		Password: u.Password,
		// Store:    u.Store,
	}
	// roles
	rawRoles := strings.Split(u.Role, ",")
	userToken.Roles = make([]string, 0)
	for _, r := range rawRoles {
		userToken.Roles = append(userToken.Roles, strings.ToLower(strings.TrimSpace(r)))
	}
	return userToken
}
