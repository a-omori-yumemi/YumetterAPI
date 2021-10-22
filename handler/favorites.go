package handler

import (
	"strconv"

	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/querier"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"github.com/labstack/echo/v4"
)

func GetTwUsrParams(c echo.Context) (twID model.TwIDType, usrID model.UsrIDType, err error) {
	if twIDTmp, err := strconv.Atoi(c.Param("tw_id")); err != nil {
		return twID, usrID, echo.NewHTTPError(400, err)
	} else {
		twID = model.TwIDType(twIDTmp)
	}
	if usrIDTmp, err := strconv.Atoi(c.Param("usr_id")); err != nil {
		return twID, usrID, echo.NewHTTPError(400, err)
	} else {
		usrID = model.UsrIDType(usrIDTmp)
	}
	return twID, usrID, nil
}

type PutFavorite struct {
	FavRepo repository.IFavoriteRepository
}

func (h PutFavorite) Handle(c echo.Context) error {
	twID, usrID, err := GetTwUsrParams(c)
	if err != nil {
		return err
	}

	fav := model.Favorite{
		TwID:  twID,
		UsrID: usrID,
	}
	fav, err = h.FavRepo.AddFavorite(fav)
	if err != nil {
		return err
	}
	return c.NoContent(200)
}

type DeleteFavorite struct {
	FavRepo repository.IFavoriteRepository
}

func (h DeleteFavorite) Handle(c echo.Context) error {
	twID, usrID, err := GetTwUsrParams(c)
	if err != nil {
		return err
	}

	err = h.FavRepo.DeleteFavorite(twID, usrID)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

type GetFavorites struct {
	FavQuerier querier.IFavoriteQuerier
}

func (h GetFavorites) GetParams(c echo.Context) (twID model.TwIDType, err error) {
	if twIDTmp, err := strconv.Atoi(c.Param("tw_id")); err != nil {
		return twID, echo.NewHTTPError(400, err)
	} else {
		twID = model.TwIDType(twIDTmp)
	}
	return twID, nil
}

func (h GetFavorites) Handle(c echo.Context) error {
	twID, err := h.GetParams(c)
	if err != nil {
		return err
	}

	favs, err := h.FavQuerier.FindFavorites(twID)
	if err != nil {
		return err
	}
	return c.JSON(200, favs)
}
