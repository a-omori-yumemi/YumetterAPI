package mysql

import (
	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
)

type MySQLUserRepository struct {
	repository.IUserRepository
}

func NewMySQLUserRepository(DB db.MySQLDB) *MySQLUserRepository {
	return &MySQLUserRepository{}
}
