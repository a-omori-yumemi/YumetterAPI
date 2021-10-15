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

	conf := db.DBConfig{
		Port:     os.Getenv("MYSQL_PORT"),
		Host:     os.Getenv("MYSQL_HOST"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
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
