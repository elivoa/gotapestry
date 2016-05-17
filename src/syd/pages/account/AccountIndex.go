package account

import (
	"github.com/elivoa/got/core"
	"strconv"
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

func (p *AccountIndex) StoreName(store int) string {
	if value, err := service.Const.GetStringValue("store", strconv.Itoa(store)); err != nil {
		panic(err)
	} else {
		return value
	}
}
