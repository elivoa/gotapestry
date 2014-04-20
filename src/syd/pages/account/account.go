package account

import (
	"got/core"
	"syd/model"
	"syd/service"
)

type AccountIndex struct {
	core.Page
	Title     string
	UserToken *model.UserToken
}

func (p *AccountIndex) Setup() {
	userToken := service.User.GetLogin(p.W, p.R)
	if nil != userToken {
		p.UserToken = userToken
	}
}
