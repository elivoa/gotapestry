package account

import (
	"fmt"
	"github.com/elivoa/got/route/exit"
	"got/core"
	"syd/dal/userdao"
	"syd/model"
)

type AccountRegister struct {
	core.Page
	Title string

	User        *model.User
	FormMessage string `scope:"flash"` // Move this message to form component.
	FormError   string `scope:"flash"` // Move this message to form component.
}

func (p *AccountRegister) OnSuccessFromRegisterForm() *exit.Exit {
	fmt.Printf("-------------- register form success -----------------\n")
	fmt.Println("Username ", p.User)

	// TODO: validate user.

	if user, err := userdao.CreateUser(p.User); err != nil {
		panic(err)
	} else {
		p.User = user
	}

	// TODO: log create action.

	return exit.Redirect("/")
}
