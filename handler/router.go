package handler

import (
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"github.com/a-omori-yumemi/YumetterAPI/usecase"
	"github.com/labstack/echo/v4"
)

func SetRoute(e *echo.Echo, repos repository.Repositories, usecases usecase.Usecases) {
	g := e.Group("/v1")

	tweetsg := g.Group("/tweets")
	tweetsg.GET("/:tw_id", GetTweet(usecases.TweetService))
	tweetsg.DELETE("/:tw_id", DeleteTweet(usecases.TweetService), AuthUserMiddleware(usecases.Authenticator))
	tweetsg.GET("/", GetTweets(usecases.TweetService))
	tweetsg.POST("/", PostTweet(usecases.TweetService), AuthUserMiddleware(usecases.Authenticator))

	usersg := g.Group("/users")
	usersg.GET("/:usr_id", GetUser(repos.UserRepo))
	usersg.POST("/", RegisterUser(repos.UserRepo))
	usersg.POST("/login", LoginUser(usecases.Authenticator))
	usersg.GET("/me", GetMe(repos.UserRepo), AuthUserMiddleware(usecases.Authenticator))
	usersg.DELETE("/me", DeleteMe(repos.UserRepo), AuthUserMiddleware(usecases.Authenticator))
	usersg.PATCH("/me", PatchMe(repos.UserRepo), AuthUserMiddleware(usecases.Authenticator))

	favoritesg := tweetsg.Group("/:tw_id/favorites")
	favoritesg.GET("/", GETFavorites(repos.FavRepo))
	favoritesg.PUT("/:usr_id", PUTFavorite(repos.FavRepo), AuthUserMiddleware(usecases.Authenticator))
	favoritesg.DELETE("/:usr_id", DELETEFavorite(repos.FavRepo), AuthUserMiddleware(usecases.Authenticator))
}
