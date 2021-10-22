package querier

import "github.com/a-omori-yumemi/YumetterAPI/model"

type ITweetQuerier interface {
	FindTweet(twID model.TwIDType) (model.Tweet, error)
	FindTweets(limit int, replied_to *model.TwIDType) ([]model.Tweet, error)
}
