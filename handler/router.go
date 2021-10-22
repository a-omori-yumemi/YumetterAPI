package handler

import (
	"github.com/a-omori-yumemi/YumetterAPI/querier"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"github.com/a-omori-yumemi/YumetterAPI/usecase"
	"github.com/labstack/echo/v4"
)

func SetRoute(e *echo.Echo, repos repository.Repositories, usecases usecase.Usecases, queriers querier.Queriers) {
	g := e.Group("/v1")
	g.Use(AuthUserMiddleware(usecases.Authenticator))

	tweetsg := g.Group("/tweets")
	tweetsg.GET("/:tw_id", GetTweet(queriers.TweetQuerier))
	tweetsg.DELETE("/:tw_id", DeleteTweet(usecases.TweetDeleteUsecase))
	tweetsg.GET("", GetTweets(queriers.TweetDetailQuerier))
	tweetsg.POST("", PostTweet(repos.TweetRepo))

	usersg := g.Group("/users")
	usersg.GET("/:usr_id", GetUser(queriers.UserQuerier))
	usersg.POST("", RegisterUser(repos.UserRepo))
	usersg.POST("/login", LoginUser(usecases.Authenticator))
	usersg.GET("/me", GetMe(repos.UserRepo))
	usersg.DELETE("/me", DeleteMe(repos.UserRepo))
	usersg.PATCH("/me", PatchMe(repos.UserRepo))

	favoritesg := tweetsg.Group("/:tw_id/favorites")
	favoritesg.GET("", GETFavorites(queriers.FavQuerier))
	favoritesg.PUT("/:usr_id", PUTFavorite(repos.FavRepo))
	favoritesg.DELETE("/:usr_id", DELETEFavorite(repos.FavRepo))
}
