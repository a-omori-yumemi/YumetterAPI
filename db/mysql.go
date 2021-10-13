package db

import (
	"github.com/jmoiron/sqlx"
)

type DB struct {
	DB *sqlx.DB
}

type DBConfig struct {
	Port     string
	Host     string
	Password string
	Database string
	User     string
}

func NewDB(conf DBConfig) (DB, error) {
	return DB{}, nil
}
