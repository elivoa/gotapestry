package account

import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"syd/dal/userdao"
	"syd/model"
	"syd/service"
)

type AccountChangePassword struct {
	core.Page
	Title string

	User        *model.User
	NewPassword string
	FormMessage string `scope:"flash"` // Move this message to form component.
	FormError   string `scope:"flash"` // Move this message to form component.
}

func (p *AccountChangePassword) Setup() {
	userToken := service.User.RequireLogin(p.W, p.R)
	p.User = &model.User{Username: userToken.Username}
}

func (p *AccountChangePassword) OnSuccessFromChangePasswordForm() *exit.Exit {
	fmt.Printf("-------------- login form success -----------------\n")
	fmt.Println("Username ", p.User)

	// verify login
	userToken := service.User.RequireLogin(p.W, p.R)

	// verify login old password
	user, err := userdao.GetUserWithCredential(userToken.Username, p.User.Password) // p.User.Password is old password;
	if err != nil {
		//		panic(err)
		p.FormError = "Error: Login failed!" + err.Error()
		return nil
	} else if user == nil {
		p.FormError = "Error: Login failed!"
		return nil
	}

	// update new password
	user.Password = p.NewPassword // set new password
	if _, err := service.User.UpdateUser(user); err != nil {
		p.FormError = "Error: Login failed!" + err.Error()
		return nil
		// panic(err)
	}
	return exit.Redirect("/")

	// _, err := service.User.Login(p.User.Username, p.User.Password, p.W, p.R)
	// if err != nil {
	// 	// error can't login, How to redirect to the current page and show errors.
	// 	p.FormError = "Error: Login failed!"

	// 	// TODO return to this
	// 	return nil

	// } else {

	// 	// service already set userToken to session and cookie. redirect if needed.

	// 	// TODO:  why this not works.
	// 	p.FormMessage = "Login Success!"

	// 	return exit.Redirect("/")
	// }
}
