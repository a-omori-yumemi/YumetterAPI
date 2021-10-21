package handler

import (
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"github.com/a-omori-yumemi/YumetterAPI/usecase"
	"github.com/labstack/echo/v4"
)

func SetRoute(e *echo.Echo, repos repository.Repositories, usecases usecase.Usecases) {
	g := e.Group("/v1")
	g.Use(AuthUserMiddleware(usecases.Authenticator))

	tweetsg := g.Group("/tweets")
	tweetsg.GET("/:tw_id", GetTweet(repos.TweetRepo))
	tweetsg.DELETE("/:tw_id", DeleteTweet(usecases.TweetDeleteUsecase))
	tweetsg.GET("", GetTweets(usecases.TweetDetailUsecase))
	tweetsg.POST("", PostTweet(repos.TweetRepo))

	usersg := g.Group("/users")
	usersg.GET("/:usr_id", GetUser(repos.UserRepo))
	usersg.POST("", RegisterUser(repos.UserRepo))
	usersg.POST("/login", LoginUser(usecases.Authenticator))
	usersg.GET("/me", GetMe(repos.UserRepo))
	usersg.DELETE("/me", DeleteMe(repos.UserRepo))
	usersg.PATCH("/me", PatchMe(repos.UserRepo))

	favoritesg := tweetsg.Group("/:tw_id/favorites")
	favoritesg.GET("", GETFavorites(repos.FavRepo))
	favoritesg.PUT("/:usr_id", PUTFavorite(repos.FavRepo))
	favoritesg.DELETE("/:usr_id", DELETEFavorite(repos.FavRepo))
}
