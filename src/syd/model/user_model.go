package model

import (
	"strings"
	"time"
)

type User struct {
	Id         int64
	Username   string
	Password   string // no use?
	HashPass   string
	Gender     string
	QQ         string
	Mobile     string
	City       string
	Role       string
	CreateTime time.Time
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
