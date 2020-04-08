package service

import (
	"fmt"
	"net/http"
	"strings"
	"syd/base"
	"syd/dal/useractiondao"
	"syd/dal/userdao"
	"syd/model"
	"time"

	"github.com/elivoa/got/config"
	"github.com/elivoa/got/core/exception"
	"github.com/elivoa/got/coreservice/sessions"
	"github.com/elivoa/got/db"
	"github.com/elivoa/got/logs"
)

var (
	USER_TOKEN_SESSION_KEY string = config.USER_TOKEN_SESSION_KEY // "USER_TOKEN_SESSION_KEY"
	USER_TOKEN_COOKIE_KEY  string = config.USER_TOKEN_COOKIE_KEY  // "USER_TOKEN_COOKIE_KEY"
)
var DEEP_TRACE = false

// TODO change session into longtime session.

// TODO Need interface & implements Design pattern.
type UserService struct {
	logs *logs.Logger // TODO: Inject request...
}

func (s *UserService) EntityManager() *db.Entity {
	return userdao.EntityManager()
}

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
	if s.logs.Trace() {
		s.logs.Printf("Enter function: RequireLogin")
	}

	if userToken := s.GetLogin(w, r); userToken == nil {
		panic(&base.LoginError{Message: "User not login.", Reason: "some reason"})
	} else {
		return userToken
	}
}

// RequireRole including RequireLogin
func (s *UserService) RequireRole(w http.ResponseWriter, r *http.Request, roles ...string) *model.UserToken {
	userToken := s.RequireLogin(w, r)
	for _, requiredRole := range roles {
		requiredRole = strings.ToLower(requiredRole)
		for _, r := range userToken.Roles {
			if r == requiredRole {
				return userToken
			}
		}
	}
	panic(exception.NewAccessDeniedErrorf(
		"Access Denied. You need to be one of the following roles: %v", roles))
}

// will be very fast.
// return true if user is login and login is available.
// return false if
func (s *UserService) GetLogin(w http.ResponseWriter, r *http.Request) *model.UserToken {
	if s.logs.Trace() {
		s.logs.Printf("Enter function: GetLogin [session]")
	}

	session := sessions.LongCookieSession(r)

	// session.Values["55667788"] = &model.UserToken{Name: "84983"}
	// session.Values["55667788--"] = "&model.UserToken{Name: 84983}"
	// session.Save(r, w)
	// for k, v := range session.Values {
	// 	fmt.Println("\t)))->:", k, " -> ", v)
	// }

	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")

	// deep trace.
	if DEEP_TRACE && s.logs.Trace() {
		s.logs.Printf("  DEEP TRACE: everything in long-session: %s {", session.ID)
		for k, v := range session.Values {
			s.logs.Printf("  DEEP TRACE: %v : %v", k, v)
		}
		s.logs.Printf("  }")
	}

	if userTokenRaw, ok := session.Values[config.USER_TOKEN_SESSION_KEY]; ok && userTokenRaw != nil {
		if s.logs.Trace() {
			s.logs.Printf("Got userToken : %v.", userTokenRaw)
		}

		if userToken := userTokenRaw.(*model.UserToken); userToken != nil {
			if s.logs.Trace() {
				s.logs.Printf("Got userToken.outdated = %v.", "TODO: check if userToken is outdated.")
			}

			// TODO: check if userToken is outdated.
			if outdated := false; !outdated {
				// TODO: update userToken.Tiemout
				if s.logs.Trace() {
					s.logs.Printf("GetLogin success, return cached user.")
				}
				return userToken
			}
		}
	}

	// if not in session, try cookie
	if userToken, err := s.LoginFromCookie(r); err != nil {
		if s.logs.Error() {
			s.logs.Printf("Login from cookie failed, reason: %s\n", err.Error())
		}
		return nil
	} else {
		if s.logs.Info() {
			s.logs.Printf("[Cookie] Login from cookie succeed, username: %s, password(hash): %s\n",
				userToken.Username, userToken.Password)
		}

		// if success, update it to session.
		s.setToSession(w, r, userToken)
		if s.logs.Trace() {
			s.logs.Printf("GetLogin success, return user.")
		}
		return userToken
	}
}

func (s *UserService) LoginFromCookie(r *http.Request) (*model.UserToken, error) {
	if s.logs.Trace() {
		s.logs.Printf("enter function: LoginFromCookie")
	}

	if credential := s.loadUserTokenFromCookie(r); credential != nil {
		if s.logs.Trace() {
			s.logs.Printf("Credential in cookie is : %v", credential)
		}
		user, err := userdao.GetUserWithCredential(credential[0], credential[1])
		if s.logs.Trace() {
			s.logs.Printf("[DB] Get user with Credenial, user is : %v", user)
		}

		if nil != user && err == nil {
			// if pass := userdao.VerifyLogin(userToken.Username, userToken.Password); pass {
			// fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			// fmt.Println("login success.")
			// fmt.Println(userToken)
			return user.ToUserToken(), nil
		} else {
			if err == nil {
				err = &base.LoginError{Message: "Username and password not matched."}
			}
			return nil, err
		}
	}
	return nil, &base.LoginError{Message: "User not login.^Y^"}
}

// return username & password pair
func (s *UserService) loadUserTokenFromCookie(r *http.Request) []string {
	if c, err := r.Cookie(USER_TOKEN_COOKIE_KEY); err == nil {
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

	user, err := userdao.GetUserWithCredential(username, password)
	if nil != user && err == nil {
		userToken := user.ToUserToken()
		s.setToSession(w, r, userToken)
		s.setToCookie(w, userToken)
		return userToken, nil
	} else {
		if err == nil {
			err = &base.LoginError{Message: "Username and password not matched."}
		}
		return nil, err
	}
}

// set UserToken to session.
func (s *UserService) setToSession(w http.ResponseWriter, r *http.Request, userToken *model.UserToken) {
	session := sessions.LongCookieSession(r)
	session.Values[USER_TOKEN_SESSION_KEY] = userToken
	if err := session.Save(r, w); err != nil {
		fmt.Println("======= ERROR TO HANDLE ========================")
		fmt.Println(err)
	}
}

// set UserToken to Cookie.
func (s *UserService) setToCookie(w http.ResponseWriter, userToken *model.UserToken) {
	http.SetCookie(w, &http.Cookie{
		Name:    USER_TOKEN_COOKIE_KEY,
		Value:   fmt.Sprintf("%s|%s", userToken.Username, userToken.Password),
		Expires: time.Now().AddDate(0, 0, 7),
		Path:    "/",
	})
}

func (s *UserService) removeUserCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   USER_TOKEN_COOKIE_KEY,
		Path:   "/",
		MaxAge: -1,
	})
}

func (s *UserService) removeUserTokenSession(w http.ResponseWriter, r *http.Request) {
	session := sessions.LongCookieSession(r)
	session.Values[USER_TOKEN_SESSION_KEY] = nil
	delete(session.Values, USER_TOKEN_SESSION_KEY)
	session.Save(r, w)
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

// --------------------------------------------------------------------------------
// The following is helper function to fill user to models.
func (s *UserService) _batchFetchUsers(ids []int64) (map[int64]*model.User, error) {
	return userdao.ListUserByIdSet(ids...)
}

func (s *UserService) BatchFetchUsers(ids ...int64) (map[int64]*model.User, error) {
	return s._batchFetchUsers(ids)
}

func (s *UserService) BatchFetchUsersByIdMap(idset map[int64]bool) (map[int64]*model.User, error) {
	var idarray = []int64{}
	if idset != nil {
		for id, _ := range idset {
			idarray = append(idarray, id)
		}
	}
	return s._batchFetchUsers(idarray)
}

// func (s *UserService) GetUserById(id int64) (*model.User, error) {
// 	return userdao.GetUserById(id)
// }

func (s *UserService) CreateUser(user *model.User) (*model.User, error) {
	dbuser, err := userdao.GetUser("username", user.Username)
	if err != nil {
		// DONE: 如何使用error才不会导致调用栈的丢失？
		panic(exception.NewCoreError(err, ""))
	}
	if dbuser != nil {
		return nil, exception.NewCoreError(nil, "User already exists for name: %s", user.Username)
	}
	return userdao.CreateUser(user)
}

func (s *UserService) UpdateUser(user *model.User) (int64, error) {
	return userdao.UpdateUser(user)
}

func (s *UserService) Total() int {
	count, err := s.EntityManager().CountAll()
	if err != nil {
		panic(err)
	}
	return count
}

// --------------------------------------------------------------------------------
// UserAction related
func (s *UserService) UserActionEntityManager() *db.Entity {
	return useractiondao.EntityManager()
}

func (s *UserService) LogUserAction(userId int64, action model.ActionType, contexts ...interface{}) error {
	return useractiondao.LogUserAction(userId, action, contexts...)
}

func (s *UserService) ListUserActionWithUsers(parser *db.QueryParser) ([]*model.UserAction, error) {
	if userActions, err := useractiondao.ListUserAction(parser); err != nil {
		return nil, err
	} else {
		if err := s.FillUserActionListWithUser(userActions); err != nil {
			return nil, err
		}
		return userActions, nil
	}
}

// orderlist is passed by pointer.
func (s *UserService) FillUserActionListWithUser(userActions []*model.UserAction) error {
	var idset = map[int64]bool{}
	for _, userAction := range userActions {
		idset[userAction.UserId] = true
	}
	usermap, err := s.BatchFetchUsersByIdMap(idset)
	if err != nil {
		return err
	}
	if nil != usermap {
		for _, userAction := range userActions {
			if user, ok := usermap[userAction.UserId]; ok {
				userAction.User = user
			}
		}
	}
	return nil
}
