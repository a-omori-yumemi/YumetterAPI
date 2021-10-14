package usecase

import (
	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
)

type ITweetService interface {
	repository.ITweetRepository
	DeleteTweetWithAuth(requestUserID model.UsrIDType, twID model.TwIDType) error
	FindTweetDetails(requestUserID *model.UsrIDType, limit int, replied_to *model.TwIDType) ([]TweetDetail, error)
}

type TweetDetail struct {
	UserName   model.UserName
	Tweet      model.Tweet
	FavCount   int
	ReplyCount int
	Favorited  bool
}

func (t TweetDetail) Validate() error {
	return t.UserName.Validate()
}
