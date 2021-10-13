package service

import (
	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
)

type TweetDetail struct {
	UserName   string
	Tweet      model.Tweet
	FavCount   int32
	ReplyCount int32
	Favorited  bool
}

type ITweetService interface {
	FindTweet(twID model.TwIDType) (model.Tweet, error)
	FindTweets(limit int, replied_to *model.TwIDType) ([]model.Tweet, error)
	AddTweet(tweet model.Tweet) (model.Tweet, error)
	DeleteTweet(twID model.TwIDType) error
	FindTweetDetails(watcherID model.TwIDType, limit int32, replied_to *model.TwIDType) ([]TweetDetail, error)
}

type TweetService struct {
	ITweetService
}

func NewTweetService(
	favRepo repository.IFavoriteRepository,
	tweetRepo repository.ITweetRepository,
	userRepo repository.IUserRepository) TweetService {
	return TweetService{}
}
