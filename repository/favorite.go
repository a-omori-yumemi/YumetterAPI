package repository

import (
	"github.com/a-omori-yumemi/YumetterAPI/model"
)

type IFavoriteRepositoty interface {
	FindFavorite(TwID model.TwIDType, UsrID model.UsrIDType) (model.Favorite, error)
	FindFavorites(TwID model.TwIDType) ([]model.Favorite, error)
	AddFavorite(model.Favorite) (model.Favorite, error)
	DeleteFavorite(TwID model.TwIDType, UsrID model.UsrIDType) error
}
