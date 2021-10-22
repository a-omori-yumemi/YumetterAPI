package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/querier"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"github.com/a-omori-yumemi/YumetterAPI/usecase"
	"github.com/labstack/echo/v4"
)

type GetUser struct {
	UserQuerier querier.IUserQuerier
}

func (h GetUser) Handle(c echo.Context) error {
	usrID, err := strconv.Atoi(c.Param("usr_id"))
	if err != nil {
		return echo.NewHTTPError(400, err)
	}

	user, err := h.UserQuerier.FindUser(model.UsrIDType(usrID))
	if err != nil {
		return err
	}

	return c.JSON(200, user)
}

type RegisterUser struct {
	UserRepo repository.IUserRepository
}

func (h RegisterUser) GetParams(c echo.Context) (model.UserName, model.Password, error) {
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

func (h RegisterUser) Handle(c echo.Context) error {
	name, password, err := h.GetParams(c)
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
	user, err = h.UserRepo.AddUser(user)
	if err == model.ErrDuplicateKey {
		return echo.NewHTTPError(409, "the username has already taken")
	} else if err != nil {
		return err
	}
	return c.JSON(200, user)
}

type LoginUser struct {
	Auth usecase.IAuthenticator
}

func (h LoginUser) GetParams(c echo.Context) (name model.UserName, pass model.Password, err error) {
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
func (h LoginUser) Handle(c echo.Context) error {
	name, pass, err := h.GetParams(c)
	if err != nil {
		return err
	}
	token, err := h.Auth.Login(name, pass)
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

type GetMe struct {
	UserQuerier repository.IUserRepository
}

func (h GetMe) Handle(c echo.Context) error {
	usrID := GetSessionUserID(c)
	if usrID == nil {
		return echo.NewHTTPError(401)
	}

	user, err := h.UserQuerier.FindUser(model.UsrIDType(*usrID))
	if err == model.ErrNotFound {
		return echo.NewHTTPError(500, "user not found")
	} else if err != nil {
		return err
	}

	return c.JSON(200, user)
}

type DeleteMe struct {
	UserRepo repository.IUserRepository
}

func (h DeleteMe) Handle(c echo.Context) error {
	usrID := GetSessionUserID(c)
	if usrID == nil {
		return echo.NewHTTPError(401)
	}

	err := h.UserRepo.DeleteUser(model.UsrIDType(*usrID))
	if err == model.ErrNotFound {
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

type PatchMe struct {
	UserRepo repository.IUserRepository
}

func (h PatchMe) GetParams(c echo.Context) (model.UsrIDType, *model.UserName, *model.Password, error) {
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

func (h PatchMe) Handle(c echo.Context) error {
	usrID, name, pass, err := h.GetParams(c)
	if err != nil {
		return err
	}

	if name != nil {
		err = h.UserRepo.UpdateName(usrID, *name)
		if err == model.ErrDuplicateKey {
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
		err = h.UserRepo.UpdatePassword(usrID, hashed)
		if err != nil {
			return err
		}
	}
	return c.NoContent(204)
}
