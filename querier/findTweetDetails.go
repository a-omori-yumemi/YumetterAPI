package querier

import (
	"github.com/a-omori-yumemi/YumetterAPI/model"
)

const MaxLimitValue = 100

type IFindTweetDetailsQuerier interface {
	FindTweetDetails(requestUserID *model.UsrIDType, limit int, replied_to *model.TwIDType) ([]TweetDetail, error)
}

type TweetDetail struct {
	UserName   model.UserName `db:"user_name" json:"user_name"`
	Tweet      model.Tweet    `db:"tweetw" json:"tweet"`
	FavCount   int            `json:"fav_count"`
	ReplyCount int            `json:"reply_count"`
	Favorited  bool           `json:"favorited"`
}

func (t TweetDetail) Validate() error {
	return t.UserName.Validate()
}
