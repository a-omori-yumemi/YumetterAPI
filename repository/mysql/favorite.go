package mysql

import (
	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
)

type MySQLFavoriteRepository struct {
	repository.IFavoriteRepositoty
}

func NewMySQLFavoriteRepository(DB db.MySQLDB) MySQLFavoriteRepository {
	return MySQLFavoriteRepository{}
}
