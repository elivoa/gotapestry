package service

import (
	"fmt"
	"github.com/elivoa/got/utils"
	"net/http"
	"strings"
	"syd/dal/userdao"
	"syd/exceptions"
	"syd/model"
	"time"
)

// 临时解决方案。全局唯一的service。TODO 研究一下Tapestry的IOC，看一下他们的service是每个request创建一个么？
// 按照他的方法来解决。
var User *UserService = &UserService{}

// TODO Need interface & implements Design pattern.
type UserService struct {
	// TODO: Inject request...
}

var USER_TOKEN_SESSION_KEY string = "USER_TOKEN_SESSION_KEY"

/*
  1. User enter any page. AuthService check if there are UserToken in session.
  2. For UserToken in session, check if it’s outdated. (5minutes timeout)
    1. Success — return this UserToken.
    2. False — If it’s outdated, re-auth, use database.
      1. Success — update and return new AuthService.
      2. False — go to error page/login page. —>
  3. Login page.
    1. Enter username and password, auth with Database.
*/

// used in methods.
func (s *UserService) RequireLogin(w http.ResponseWriter, r *http.Request) *model.UserToken {
	if userToken := s.GetLogin(w, r); userToken == nil {
		panic(&exceptions.LoginError{Message: "User not login.", Reason: "some reason"})
	} else {
		return userToken
	}
}

// RequireRole including RequireLogin
func (s *UserService) RequireRole(w http.ResponseWriter, r *http.Request, role string) *model.UserToken {
	userToken := s.RequireLogin(w, r)
	lowerRole := strings.ToLower(role)
	found := false
	for _, r := range userToken.Roles {
		if r == lowerRole {
			found = true
			break
		}
	}
	fmt.Println(found)
	if found {
		return userToken
	} else {
		fmt.Println("access denied?")
		panic(&exceptions.AccessDeniedError{Message: "Access Denied."})
	}
}

// will be very fast.
// return true if user is login and login is available.
// return false if
func (s *UserService) GetLogin(w http.ResponseWriter, r *http.Request) *model.UserToken {
	session := utils.Session(r)
	// { // debug print.
	// 	fmt.Printf("\t >>>>>>>>>>>>>>>>>>>>>>>>>>> Session.Values: %v\n", session.Values)
	// 	for k, v := range session.Values {
	// 		fmt.Printf("key %v --> value: %v\n", k, v)
	// 	}
	// }
	if userTokenRaw, ok := session.Values[USER_TOKEN_SESSION_KEY]; ok && userTokenRaw != nil {
		if userToken := userTokenRaw.(*model.UserToken); userToken != nil {
			// TODO: check if userToken is outdated.
			if outdated := false; !outdated {
				// TODO: update userToken.Tiemout
				return userToken
			}
		}
	}
	// if not in session, try cookie
	if userToken, err := s.LoginFromCookie(r); err != nil {
		// TODO: change to log.blablabla
		fmt.Printf("Login from cookie failed, reason: %s\n", err.Error())
		return nil
	} else {
		fmt.Printf("Login from cookie succee, username: %s, password(hash): %s\n",
			userToken.Username, userToken.Password)

		// if success, update it to session.
		s.setToSession(w, r, userToken)
		return userToken
	}
}

func (s *UserService) LoginFromCookie(r *http.Request) (*model.UserToken, error) {
	if credential := s.loadUserTokenFromCookie(r); credential != nil {
		user := userdao.GetUserWithCredential(credential[0], credential[1])
		if nil != user {
			// if pass := userdao.VerifyLogin(userToken.Username, userToken.Password); pass {
			// fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			// fmt.Println("login success.")
			// fmt.Println(userToken)
			return user.ToUserToken(), nil
		} else {
			return nil, &exceptions.LoginError{Message: "Username and password not matched."}
		}
	}
	return nil, &exceptions.LoginError{Message: "User not login."}
}

// return username & password pair
func (s *UserService) loadUserTokenFromCookie(r *http.Request) []string {
	if c, err := r.Cookie(USER_TOKEN_SESSION_KEY); err == nil {
		if nil == c {
			return nil
		}
		if splits := strings.Split(c.Value, "|"); len(splits) >= 2 {
			return splits
		}
	}
	return nil
}

// Login accept username and password then verify them.
// TODO actually I want to receive a hashed password, to reduse risk of ....
// TODO: remove w, r in parameter.
func (s *UserService) Login(username string, password string,
	w http.ResponseWriter, r *http.Request) (*model.UserToken, error) {

	// 1. verify username and password with db.
	// 2. if success, set it to session and cookie.
	// 3. if not , return error.
	// TEST: always return true.

	user := userdao.GetUserWithCredential(username, password)
	if nil != user {
		userToken := user.ToUserToken()
		s.setToSession(w, r, userToken)
		s.setToCookie(w, userToken)
		return userToken, nil
	} else {
		return nil, &exceptions.LoginError{Message: "Username and password not matched."}
	}
}

// set UserToken to session.
func (s *UserService) setToSession(w http.ResponseWriter, r *http.Request, userToken *model.UserToken) {
	session := utils.Session(r)
	session.Values[USER_TOKEN_SESSION_KEY] = userToken
	fmt.Printf("\n\nSave to Session \n")
	session.Save(r, w)
}

// set UserToken to Cookie.
func (s *UserService) setToCookie(w http.ResponseWriter, userToken *model.UserToken) {
	http.SetCookie(w, &http.Cookie{
		Name:    USER_TOKEN_SESSION_KEY,
		Value:   fmt.Sprintf("%s|%s", userToken.Username, userToken.Password),
		Expires: time.Now().AddDate(0, 0, 7),
		Path:    "/",
	})
}

func (s *UserService) removeUserCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   USER_TOKEN_SESSION_KEY,
		Path:   "/",
		MaxAge: -1,
	})
}

func (s *UserService) removeUserTokenSession(w http.ResponseWriter, r *http.Request) {
	session := utils.Session(r)
	session.Values[USER_TOKEN_SESSION_KEY] = nil
	delete(session.Values, USER_TOKEN_SESSION_KEY)
	session.Save(r, w)

	// context.Set(r, USER_TOKEN_SESSION_KEY, nil)
	// context.Delete(r, USER_TOKEN_SESSION_KEY)
}

func (s *UserService) Logout(w http.ResponseWriter, r *http.Request) {
	s.removeUserCookie(w)
	s.removeUserTokenSession(w, r)
}

// HasRole
func (s *UserService) HasRole(w http.ResponseWriter, r *http.Request, role string) bool {
	userToken := s.GetLogin(w, r)
	if nil == userToken {
		return false
	}

	lowerRole := strings.ToLower(role)
	for _, r := range userToken.Roles {
		if r == lowerRole {
			return true
		}
	}
	return false
}
