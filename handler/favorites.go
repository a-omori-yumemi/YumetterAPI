package handler

import (
	"strconv"

	"github.com/a-omori-yumemi/YumetterAPI/model"
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

func PUTFavorite(favRepo repository.IFavoriteRepository) echo.HandlerFunc {

	return func(c echo.Context) error {
		twID, usrID, err := GetTwUsrParams(c)
		if err != nil {
			return err
		}

		fav := model.Favorite{
			TwID:  twID,
			UsrID: usrID,
		}
		fav, err = favRepo.AddFavorite(fav)
		if err != nil {
			return err
		}
		return c.NoContent(200)
	}
}
func DELETEFavorite(favRepo repository.IFavoriteRepository) echo.HandlerFunc {

	return func(c echo.Context) error {
		twID, usrID, err := GetTwUsrParams(c)
		if err != nil {
			return err
		}

		err = favRepo.DeleteFavorite(twID, usrID)
		if err != nil {
			return err
		}
		return c.NoContent(204)
	}
}

func GETFavorites(favRepo repository.IFavoriteRepository) echo.HandlerFunc {
	GetParams := func(c echo.Context) (twID model.TwIDType, err error) {
		if twIDTmp, err := strconv.Atoi(c.Param("tw_id")); err != nil {
			return twID, echo.NewHTTPError(400, err)
		} else {
			twID = model.TwIDType(twIDTmp)
		}
		return twID, nil
	}

	return func(c echo.Context) error {
		twID, err := GetParams(c)
		if err != nil {
			return err
		}

		favs, err := favRepo.FindFavorites(twID)
		if err != nil {
			return err
		}
		return c.JSON(200, favs)
	}

}
