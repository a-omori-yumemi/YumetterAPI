package handler

import (
	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/labstack/echo/v4"
)

const ContextUserKey = "USER_ID_KEY"

func GetSessionUserID(c echo.Context) *model.UsrIDType {
	var usrID *model.UsrIDType = nil
	usrIDTmp, ok := c.Get(ContextUserKey).(model.UsrIDType)
	if ok {
		usrID = new(model.UsrIDType)
		*usrID = model.UsrIDType(usrIDTmp)
	}
	return usrID
}

func AuthUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set(ContextUserKey, 12345)
		return next(c)
	}
}
