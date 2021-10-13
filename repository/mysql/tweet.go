package mysql

import (
	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
)

type MySQLTweetRepository struct {
	repository.ITweetRepository
}

func NewMySQLTweetRepository(DB db.DB) MySQLTweetRepository {
	return MySQLTweetRepository{}
}
