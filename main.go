package main

import (
	"os"
	"strconv"

	"net/http"
	_ "net/http/pprof"

	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/handler"
	querier_tweet_detail "github.com/a-omori-yumemi/YumetterAPI/querier/findTweetDetails"
	"github.com/a-omori-yumemi/YumetterAPI/usecase"
	"github.com/felixge/fgprof"
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
	http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
	go func() {
		log.Print(http.ListenAndServe(":6060", nil))
	}()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Pre(middleware.RemoveTrailingSlash())
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	handlers, err := buildHandlers()
	if err != nil {
		log.Fatal(err)
	}
	SetRoute(e, handlers)
	e.Logger.Fatal("failed to start server", e.Start(":"+port))
}

func ProvideRODBConfig() db.RODBConfig {
	return db.RODBConfig{
		Port:            os.Getenv("MYSQL_PORT"),
		Host:            os.Getenv("MYSQL_READ_HOST"),
		User:            os.Getenv("MYSQL_USER"),
		Password:        os.Getenv("MYSQL_PASSWORD"),
		Database:        os.Getenv("MYSQL_DATABASE"),
		MaxOpenConns:    os.Getenv("MYSQL_READ_MAX_OPEN_CONNS"),
		MaxIdleConns:    os.Getenv("MYSQL_READ_MAX_IDLE_CONNS"),
		ConnMaxIdletime: os.Getenv("MYSQL_READ_CONN_MAX_IDLE_TIME"),
	}
}
func ProvideDBConfig() db.DBConfig {
	return db.DBConfig{
		Port:            os.Getenv("MYSQL_PORT"),
		Host:            os.Getenv("MYSQL_WRITE_HOST"),
		User:            os.Getenv("MYSQL_USER"),
		Password:        os.Getenv("MYSQL_PASSWORD"),
		Database:        os.Getenv("MYSQL_DATABASE"),
		MaxOpenConns:    os.Getenv("MYSQL_WRITE_MAX_OPEN_CONNS"),
		MaxIdleConns:    os.Getenv("MYSQL_WRITE_MAX_IDLE_CONNS"),
		ConnMaxIdletime: os.Getenv("MYSQL_WRITE_CONN_MAX_IDLE_TIME"),
	}
}

func ProvideTweetDetailsCacheLifeTime() (querier_tweet_detail.TweetDetailsCacheLifeTime, error) {
	t, err := strconv.Atoi(os.Getenv("TWEET_DETAILS_CACHE_LIFE_TIME"))
	return querier_tweet_detail.TweetDetailsCacheLifeTime(t), err
}

func ProvideSecretKey() usecase.SecretKey {
	return usecase.SecretKey(os.Getenv("SECRET_KEY"))
}
func SetRoute(e *echo.Echo, h handler.Handlers) {
	g := e.Group("/v1")
	g.Use(h.AuthUserMiddleware.Handle)

	tweetsg := g.Group("/tweets")
	tweetsg.GET("/:tw_id", h.GetTweet.Handle)
	tweetsg.DELETE("/:tw_id", h.DeleteTweet.Handle)
	tweetsg.GET("", h.GetTweets.Handle)
	tweetsg.POST("", h.PostTweet.Handle)

	usersg := g.Group("/users")
	usersg.GET("/:usr_id", h.GetUser.Handle)
	usersg.POST("", h.RegisterUser.Handle)
	usersg.POST("/login", h.LoginUser.Handle)
	usersg.GET("/me", h.GetMe.Handle)
	usersg.DELETE("/me", h.DeleteMe.Handle)
	usersg.PATCH("/me", h.PatchMe.Handle)

	favoritesg := tweetsg.Group("/:tw_id/favorites")
	favoritesg.GET("", h.GetFavorites.Handle)
	favoritesg.PUT("/:usr_id", h.PutFavorite.Handle)
	favoritesg.DELETE("/:usr_id", h.DeleteFavorite.Handle)
}
