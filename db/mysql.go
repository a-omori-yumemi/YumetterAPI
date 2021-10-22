package db

import (
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
)

type MySQLDB struct {
	DB *sqlx.DB
}

type ReadOnlyDB interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}

type MySQLReadOnlyDB struct {
	DB ReadOnlyDB
}

type DBConfig struct {
	Port            string
	Host            string
	Password        string
	Database        string
	User            string
	MaxOpenConns    string
	MaxIdleConns    string
	ConnMaxLifetime string
}

func NewMySQLDB(conf DBConfig) (MySQLDB, error) {
	dsn := conf.User + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Database + "?parseTime=true&multiStatements=true"
	log.Info(dsn)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return MySQLDB{}, err
	}

	maxOpenConns, err := strconv.Atoi(conf.MaxOpenConns)
	if err != nil {
		return MySQLDB{}, err
	}
	maxIdleConns, err := strconv.Atoi(conf.MaxIdleConns)
	if err != nil {
		return MySQLDB{}, err
	}
	connMaxLifetime, err := strconv.Atoi(conf.ConnMaxLifetime)
	if err != nil {
		return MySQLDB{}, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Second * time.Duration(connMaxLifetime))

	// _, err = sqlx.LoadFile(db, "db/init.sql")
	return MySQLDB{db}, err
}

func NewMySQLReadOnlyDB(conf DBConfig) (MySQLReadOnlyDB, error) {
	db, err := NewMySQLDB(conf)
	return MySQLReadOnlyDB{db.DB}, err
}
