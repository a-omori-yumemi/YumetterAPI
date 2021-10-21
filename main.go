package main

import (
	"os"

	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/handler"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"github.com/a-omori-yumemi/YumetterAPI/repository/mysql"
	"github.com/a-omori-yumemi/YumetterAPI/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e := echo.New()
	e.Use(middleware.Logger())
	e.Pre(middleware.RemoveTrailingSlash())
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	wconf := db.DBConfig{
		Port:     os.Getenv("MYSQL_PORT"),
		Host:     os.Getenv("MYSQL_WRITE_HOST"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
	}
	rconf := db.DBConfig{
		Port:     os.Getenv("MYSQL_PORT"),
		Host:     os.Getenv("MYSQL_READ_HOST"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
	}
	repos, usecases := construct(wconf, rconf)
	handler.SetRoute(e, repos, usecases)
	e.Logger.Fatal("failed to start server", e.Start(":"+port))
}

func construct(wconf db.DBConfig, rconf db.DBConfig) (repository.Repositories, usecase.Usecases) {
	DB, err := db.NewMySQLDB(wconf)
	ReplicaDB, err := db.NewMySQLReadOnlyDB(rconf)
	if err != nil {
		log.Fatal("failed to connect DB ", err)
	}
	repos := repository.Repositories{
		FavRepo:   mysql.NewMySQLFavoriteRepository(DB),
		TweetRepo: mysql.NewMySQLTweetRepository(DB),
		UserRepo:  mysql.NewMySQLUserRepository(DB),
	}
	services := usecase.Usecases{
		TweetDetailUsecase: usecase.NewTweetDetailQuerier(DB),
		TweetDeleteUsecase: usecase.NewTweetDeleteUsecase(repos.TweetRepo),
		Authenticator: usecase.NewJWTAuthenticator(
			repos.UserRepo,
			os.Getenv("SECRET_KEY"),
		),
	}

	return repos, services
}
