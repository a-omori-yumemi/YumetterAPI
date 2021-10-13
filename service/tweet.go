package service

import (
	"github.com/a-omori-yumemi/YumetterAPI/model"
)

type ITweetService interface {
	FindTweet(twID model.TwIDType) (model.Tweet, error)
	FindTweets(limit int, replied_to *model.TwIDType) ([]model.Tweet, error)
	AddTweet(tweet model.Tweet) (model.Tweet, error)
	DeleteTweet(twID model.TwIDType) error
	FindTweetDetail()
}
