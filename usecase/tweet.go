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
	UserName   model.UserName `json:"user_name"`
	Tweet      model.Tweet    `json:"tweet"`
	FavCount   int            `json:"fav_count"`
	ReplyCount int            `json:"reply_count"`
	Favorited  bool           `json:"favorited"`
}

func (t TweetDetail) Validate() error {
	return t.UserName.Validate()
}
