package querier_mysql

import (
	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/model"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLFavoriteQuerier struct {
	db db.MySQLReadOnlyDB
}

func NewMySQLFavoriteQuerier(DB db.MySQLReadOnlyDB) *MySQLFavoriteQuerier {
	return &MySQLFavoriteQuerier{db: DB}
}

func (r *MySQLFavoriteQuerier) FindFavorite(TwID model.TwIDType, UsrID model.UsrIDType) (ret model.Favorite, err error) {
	err = r.db.DB.Get(&ret, "SELECT * FROM Favorite WHERE tw_id=? and usr_id=?", TwID, UsrID)
	return ret, db.InterpretMySQLError(err)
}
func (r *MySQLFavoriteQuerier) FindFavorites(TwID model.TwIDType) (favs []model.Favorite, err error) {
	favs = []model.Favorite{}
	err = r.db.DB.Select(&favs, "SELECT * FROM Favorite WHERE tw_id=?", TwID)
	return favs, db.InterpretMySQLError(err)
}
