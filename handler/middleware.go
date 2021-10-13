package handler

import (
	"github.com/labstack/echo/v4"
)

const ContextUserKey = "USER_ID_KEY"

func AuthUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set(ContextUserKey, 12345)
		return next(c)
	}
}
