package main

import (
	"fmt"
	"os"

	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/handler"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"github.com/a-omori-yumemi/YumetterAPI/repository/mysql"
	"github.com/a-omori-yumemi/YumetterAPI/usecase"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
)

type DBConfig struct {
	Port     string
	Host     string
	Password string
	Database string
	User     string
}

func main() {
	conf := DBConfig{
		Port:     os.Getenv("MYSQL_PORT"),
		Host:     os.Getenv("MYSQL_HOST"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
	}
	dsn := conf.User + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Database + "?parseTime=true&multiStatements=true"
	fmt.Print(dsn)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Error(err)
		return
	}
	repos, usecases := construct(conf)
	handler.SetRoute(e, repos, usecases)
	e.Logger.Fatal("failed to start server", e.Start(":"+port))
}

func construct(conf db.DBConfig) (repository.Repositories, usecase.Usecases) {
	DB, err := db.NewMySQLDB(conf)
	if err != nil {
		log.Fatal("failed to connect DB ", err)
	}
	repos := repository.Repositories{
		FavRepo:   mysql.NewMySQLFavoriteRepository(DB),
		TweetRepo: mysql.NewMySQLTweetRepository(DB),
		UserRepo:  mysql.NewMySQLUserRepository(DB),
	}
	services := usecase.Usecases{
		TweetService: usecase.NewTweetService(
			repos.FavRepo,
			repos.TweetRepo,
			repos.UserRepo,
		),
		Authenticator: usecase.NewJWTAuthenticator(
			repos.UserRepo,
			os.Getenv("SECRET_KEY"),
		),
	}

	return repos, services
}
