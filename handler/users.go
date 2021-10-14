package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"github.com/a-omori-yumemi/YumetterAPI/usecase"
	"github.com/labstack/echo/v4"
)

func GetUser(userRepo repository.IUserRepository) echo.HandlerFunc {

	return func(c echo.Context) error {
		usrID, err := strconv.Atoi(c.Param("usr_id"))
		if err != nil {
			return echo.NewHTTPError(400, err)
		}

		user, err := userRepo.FindUser(model.UsrIDType(usrID))
		if err != nil {
			return err
		}

		return c.JSON(200, user)
	}
}

func RegisterUser(userRepo repository.IUserRepository) echo.HandlerFunc {

	GetParams := func(c echo.Context) (model.UserName, model.Password, error) {
		name := model.UserName(c.FormValue("name"))
		if err := name.Validate(); err != nil {
			return name, "", echo.NewHTTPError(400, err)
		}

		password := model.Password(c.FormValue("password"))
		if err := password.Validate(); err != nil {
			return name, password, echo.NewHTTPError(400, err)
		}
		return name, password, nil
	}

	return func(c echo.Context) error {
		name, password, err := GetParams(c)
		if err != nil {
			return err
		}

		// この処理はHandlerにあるべきではないと思うが、多少は
		hashed, err := password.Hash()
		if err != nil {
			return err
		}

		user := model.User{
			Name:           name,
			HashedPassword: hashed,
		}
		if err := user.Validate(); err != nil {
			return err
		}
		user, err = userRepo.AddUser(user)
		if err == repository.ErrDuplicateKey {
			return echo.NewHTTPError(409, "the username has already taken")
		} else if err != nil {
			return err
		}
		return c.JSON(200, user)
	}
}

func LoginUser(auth usecase.IAuthenticator) echo.HandlerFunc {

	GetParams := func(c echo.Context) (name model.UserName, pass model.Password, err error) {
		name = model.UserName(c.FormValue("name"))
		if err := name.Validate(); err != nil {
			return name, pass, echo.NewHTTPError(400, err)
		}
		pass = model.Password(c.FormValue("password"))
		if err := pass.Validate(); err != nil {
			return name, pass, echo.NewHTTPError(400, err)
		}

		return name, pass, nil
	}

	// headerのセット
	return func(c echo.Context) error {
		name, pass, err := GetParams(c)
		if err != nil {
			return err
		}
		token, err := auth.Login(name, pass)
		if err != nil {
			return echo.NewHTTPError(401, err)
		}

		c.SetCookie(&http.Cookie{
			Name:     SessionCookieName,
			Value:    token,
			MaxAge:   int(time.Hour) * 24 * 3,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
			HttpOnly: true,
		})
		return nil
	}
}

func GetMe(userRepo repository.IUserRepository) echo.HandlerFunc {

	return func(c echo.Context) error {
		usrID := GetSessionUserID(c)
		if usrID == nil {
			return echo.NewHTTPError(401)
		}

		user, err := userRepo.FindUser(model.UsrIDType(*usrID))
		if err == repository.ErrNotFound {
			return echo.NewHTTPError(500, "user not found")
		} else if err != nil {
			return err
		}

		return c.JSON(200, user)
	}
}

func DeleteMe(userRepo repository.IUserRepository) echo.HandlerFunc {

	return func(c echo.Context) error {
		usrID := GetSessionUserID(c)
		if usrID == nil {
			return echo.NewHTTPError(401)
		}

		err := userRepo.DeleteUser(model.UsrIDType(*usrID))
		if err == repository.ErrNotFound {
			return echo.NewHTTPError(500, "user not found")
		} else if err != nil {
			return err
		}

		c.SetCookie(&http.Cookie{
			Name:   SessionCookieName,
			MaxAge: -1,
		})

		return c.NoContent(204)
	}
}

func PatchMe(userRepo repository.IUserRepository) echo.HandlerFunc {

	GetParams := func(c echo.Context) (model.UsrIDType, *model.UserName, *model.Password, error) {
		usrID := GetSessionUserID(c)
		if usrID == nil {
			return 0, nil, nil, echo.NewHTTPError(401)
		}

		var name *model.UserName = nil
		if nameTmp := model.UserName(c.FormValue("name")); nameTmp != "" {
			name = new(model.UserName)
			*name = nameTmp
			if err := name.Validate(); err != nil {
				return *usrID, name, nil, echo.NewHTTPError(400, err)
			}
		}

		var password *model.Password = nil
		if passwordTmp := c.FormValue("password"); passwordTmp != "" {
			password = new(model.Password)
			*password = model.Password(passwordTmp)
			if err := password.Validate(); err != nil {
				return *usrID, name, nil, echo.NewHTTPError(400, err)
			}
		}

		return *usrID, name, password, nil
	}

	return func(c echo.Context) error {
		usrID, name, pass, err := GetParams(c)
		if err != nil {
			return err
		}

		if name != nil {
			err = userRepo.UpdateName(usrID, *name)
			if err == repository.ErrDuplicateKey {
				return echo.NewHTTPError(409, "the username has already taken")
			} else if err != nil {
				return err
			}
		}

		if pass != nil {
			hashed, err := pass.Hash()
			if err != nil {
				return err
			}
			err = userRepo.UpdatePassword(usrID, hashed)
			if err != nil {
				return err
			}
		}
		return c.NoContent(204)
	}
}
