package mysql

import (
	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
)

type MySQLFavoriteRepository struct {
	repository.IFavoriteRepository
}

func NewMySQLFavoriteRepository(DB db.MySQLDB) MySQLFavoriteRepository {
	return MySQLFavoriteRepository{}
}
