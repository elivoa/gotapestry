package account

import (
	"github.com/elivoa/got/builtin/services"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"syd/model"
	"syd/service"
)

type AccountLogin struct {
	core.Page
	Title string

	LoginUser   *model.User
	FormMessage string // `scope:"flash"` //
	FormError   string `query:"errmsg"` // use query to immulate Flash message.
}

func (p *AccountLogin) OnSuccessFromLoginForm() *exit.Exit {
	// fmt.Printf("-------------- login form success -----------------\n")
	// fmt.Println("Username ", p.LoginUser)

	_, err := service.User.Login(p.LoginUser.Username, p.LoginUser.Password, p.W, p.R)
	if err != nil {
		// error can't login, How to redirect to the current page and show errors.
		p.FormError = "Error: Login failed!"

		// TODO: immulate flash message. automatically return empty page with parameter.
		url := services.Link.GeneratePageUrlWithContextAndQueryParameters("account/login",
			map[string]interface{}{"errmsg": "Login failed! " + err.Error()},
		)
		return exit.Redirect(url) // return nil // <-- should return nil
	} else {
		// service already set userToken to session and cookie. redirect if needed.

		p.FormMessage = "Login Success!" // nouse! No one can see this.
		return exit.Redirect("/")        // Return to homepage; TODO: return to where I comes from!
	}
}

// TODO: Should be moved to common place.
// func (p *AccountLogin) OnSetTimeZone(offset int) *exit.Exit {
// 	timezone := model.NewTimeZoneInfo(offset)
// 	service.TimeZone.SaveTimeZone(p.W, p.R, timezone)
// 	return exit.RenderText(timezone.String())
// }
