package handler

import (
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"github.com/a-omori-yumemi/YumetterAPI/service"
	"github.com/labstack/echo/v4"
)

func SetRoute(e *echo.Echo, repos repository.Repositories, services service.Services) {
	g := e.Group("/v1")

	tweetsg := g.Group("/tweets")
	tweetsg.POST("/", PostTweet(services.TweetService))
	tweetsg.GET("/:tw_id", GetTweet(services.TweetService))

}
