package service

import (
	"fmt"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/core/exception"
	"github.com/elivoa/got/coreservice/sessions"
	"log"
	"net/http"
	"syd/model"
)

var TimeZoneKey string = config.TIMEZONE_SESSION_KEY // "USER_TOKEN_SESSION_KEY"
var EmptyTimeZone = &model.TimeZoneInfo{}

// TODO change session into longtime session.

var TimeZone = &TimeZoneService{}

type TimeZoneService struct {
	// TODO: Inject request...
}

func (s *TimeZoneService) UserTimeZoneSafe(request *http.Request) *model.TimeZoneInfo {
	// Check TimeZone settings.
	sessionId := sessions.GetSessionId(request)
	if tzi := sessions.Get(sessionId, TimeZoneKey); tzi != nil {
		if timezoneinfo, ok := (tzi).(*model.TimeZoneInfo); ok {
			return timezoneinfo
		} else {
			log.Printf("Warrning! TimeZone in session is not type model.TimeZoneInfo!!")
		}
	} else {
		// if not, check long live session.i.e. cookie.
		cookies := sessions.LongCookieSession(request)
		if offset, ok := cookies.Values[TimeZoneKey]; ok && nil != offset {
			if offsetInt, ok := offset.(int); ok {
				// if cookie has timezone, pass. Set to session
				sessions.Set(sessionId, TimeZoneKey, model.NewTimeZoneInfo(offsetInt))
			} else {
				log.Printf("Warrning! Offset in cookie are not int, offset: %v", offset)
			}
		} else {
			// if not found in cookie; jump to exception page.
			fmt.Println("\n\n\njjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj")
			panic(exception.NewTimeZoneNotFoundErrorf("TimeZoneInfo not found in session."))
		}

	}
	return EmptyTimeZone
}

// temp function
func (s *TimeZoneService) UserTimeZoneDontCheckCookie(request *http.Request) *model.TimeZoneInfo {
	// Check TimeZone settings.
	sessionId := sessions.GetSessionId(request)
	if tzi := sessions.Get(sessionId, TimeZoneKey); tzi != nil {
		if timezoneinfo, ok := (tzi).(*model.TimeZoneInfo); ok {
			return timezoneinfo
		} else {
			log.Printf("Warrning! TimeZone in session is not type model.TimeZoneInfo!!")
		}
	} else {

		// 有这个else， 每次重启都需要登陆。否则不需要登陆，会记录cookie。

		// if not, check long live session.i.e. cookie.
		cookies := sessions.LongCookieSession(request)
		if offset, ok := cookies.Values[TimeZoneKey]; ok && nil != offset {
			if offsetInt, ok := offset.(int); ok {
				// if cookie has timezone, pass. Set to session
				sessions.Set(sessionId, TimeZoneKey, model.NewTimeZoneInfo(offsetInt))
			} else {
				log.Printf("Warrning! Offset in cookie are not int, offset: %v", offset)
			}
		} else {
			// 这个版本不redirect到登陆页面。因为调用者已经在登陆页面了。
			// if not found in cookie; jump to exception page.
			// fmt.Println("\n\n\njjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj")
			// panic(exception.NewTimeZoneNotFoundErrorf("TimeZoneInfo not found in session."))
		}

	}
	return EmptyTimeZone
}

func (s *TimeZoneService) SaveTimeZone(response http.ResponseWriter, request *http.Request,
	timezone *model.TimeZoneInfo) {
	// set both cookie and session
	sessionId := sessions.GetSessionId(request)
	sessions.Set(sessionId, TimeZoneKey, timezone)
	cookies := sessions.LongCookieSession(request)
	cookies.Values[TimeZoneKey] = timezone.Offset
	cookies.Save(request, response)
	// fmt.Println("Set timezone to ", offset, " to session : ", sessionId)
}
