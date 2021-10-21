package querier

import (
	"github.com/a-omori-yumemi/YumetterAPI/model"
)

type IFavoriteQuerier interface {
	FindFavorite(TwID model.TwIDType, UsrID model.UsrIDType) (model.Favorite, error)
	FindFavorites(TwID model.TwIDType) ([]model.Favorite, error)
}
