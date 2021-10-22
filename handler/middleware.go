package handler

import (
	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/usecase"
	"github.com/labstack/echo/v4"
)

const ContextUserKey = "USER_ID_KEY"
const SessionCookieName = "SESSION"

func GetSessionUserID(c echo.Context) *model.UsrIDType {
	var usrID *model.UsrIDType = nil
	usrIDTmp, ok := c.Get(ContextUserKey).(model.UsrIDType)
	if ok {
		usrID = new(model.UsrIDType)
		*usrID = model.UsrIDType(usrIDTmp)
	}
	return usrID
}

type AuthUserMiddleware struct {
	Auth usecase.IAuthenticator
}

func (m AuthUserMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(SessionCookieName)
		if err == nil {
			usrID, err := m.Auth.UseSession(cookie.Value)
			if err == nil {
				c.Set(ContextUserKey, usrID)
			}
		}
		return next(c)
	}
}
