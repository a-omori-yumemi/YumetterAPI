package handler

import (
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"github.com/labstack/echo/v4"
)

func SetRoute(e *echo.Echo, repos repository.Repositories) {
	g := e.Group("/v1")

	tweetsg := g.Group("/tweets")
	tweetsg.POST("/", PostTweet(repos.TweetRepo))
	tweetsg.POST("/:tw_id", GetTweet(repos.TweetRepo))

}
