package mysql

import (
	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLFavoriteRepository struct {
	repository.IFavoriteRepository
	db db.MySQLDB
}

func NewMySQLFavoriteRepository(DB db.MySQLDB) *MySQLFavoriteRepository {
	return &MySQLFavoriteRepository{db: DB}
}

func (r *MySQLFavoriteRepository) FindFavorite(TwID model.TwIDType, UsrID model.UsrIDType) (ret model.Favorite, err error) {
	err = r.db.DB.Get(&ret, "SELECT FROM Favorite WHERE tw_id=? and usr_id=?", TwID, UsrID)
	return ret, interpretMySQLError(err)
}
func (r *MySQLFavoriteRepository) FindFavorites(TwID model.TwIDType) (favs []model.Favorite, err error) {
	err = r.db.DB.Select(&favs, "SELECT FROM Favorite WHERE tw_id=?", TwID)
	return favs, interpretMySQLError(err)
}

// ignores duplicate key
func (r *MySQLFavoriteRepository) AddFavorite(fav model.Favorite) (model.Favorite, error) {
	_, err := r.db.DB.Exec("INSERT IGNORE INTO Favorite (tw_id, usr_id) VALUES (?,?)", fav.TwID, fav.UsrID)
	return fav, err
}

// ignores not found key
func (r *MySQLFavoriteRepository) DeleteFavorite(TwID model.TwIDType, UsrID model.UsrIDType) error {
	_, err := r.db.DB.Exec("DELETE IGNORE FROM Favorite where tw_id=? and usr_id=?", TwID, UsrID)
	return err
}
