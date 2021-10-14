package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
)

type MySQLDB struct {
	DB *sqlx.DB
}

type DBConfig struct {
	Port     string
	Host     string
	Password string
	Database string
	User     string
}

func NewMySQLDB(conf DBConfig) (MySQLDB, error) {
	dsn := conf.User + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Database + "?parseTime=true&multiStatements=true"
	log.Info(dsn)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return MySQLDB{}, err
	}

	// _, err = sqlx.LoadFile(db, "db/init.sql")
	return MySQLDB{db}, err
}
