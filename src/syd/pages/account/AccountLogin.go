package account

import (
	"fmt"
	"github.com/elivoa/got/route/exit"
	"got/core"
	"syd/model"
	"syd/service"
)

type AccountLogin struct {
	core.Page
	Title string

	User        *model.User
	FormMessage string `scope:"flash"` // Move this message to form component.
	FormError   string `scope:"flash"` // Move this message to form component.
}

func (p *AccountLogin) OnSuccessFromLoginForm() *exit.Exit {
	fmt.Printf("-------------- login form success -----------------\n")
	fmt.Println("Username ", p.User)

	_, err := service.User.Login(p.User.Username, p.User.Password, p.W, p.R)
	if err != nil {
		// error can't login, How to redirect to the current page and show errors.
		p.FormError = "Error: Login failed!"

		// TODO return to this
		return nil

	} else {

		// service already set userToken to session and cookie. redirect if needed.

		// TODO:  why this not works.
		p.FormMessage = "Login Success!"

		return exit.Redirect("/")
	}
}
