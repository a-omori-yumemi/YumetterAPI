package querier_tweet_detail

import (
	"fmt"

	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/querier"
)

type ICommonTweetDetailsQuerier interface {
	FindCommonTweetDetails(limit int, replied_to *model.TwIDType) ([]CommonTweetDetail, error)
}

type CommonTweetDetail struct {
	UserName model.UserName `db:"user_name"`
	model.Tweet
	FavCount   int `db:"fav_count"`
	ReplyCount int `db:"reply_count"`
}

type IFindFavoritesByRangeQuerier interface {
	// required to be sorted by tw_id (DESC)
	FindFavoritesByRange(firstTwID model.TwIDType, lastTwID model.TwIDType, usrID model.UsrIDType) ([]model.Favorite, error)
}

type TweetDetailsQuerier struct {
	cache                       CacheStore
	commonTweetDetailQuerier    ICommonTweetDetailsQuerier
	findFavoritesByRangeQuerier IFindFavoritesByRangeQuerier
}

type CacheStore interface {
	Set(interface{})
	Get() interface{}
}

type MemCacheStore struct {
	CacheStore
	cache string
}

func NewCacheStore() *MemCacheStore {
	return &MemCacheStore{cache: ""}
}

func NewTweetDetailQuerier(
	commonTweetDetailQuerier ICommonTweetDetailsQuerier,
	findFavoritesByRangeQuerier IFindFavoritesByRangeQuerier,
	cache CacheStore) *TweetDetailsQuerier {

	return &TweetDetailsQuerier{
		cache:                       cache,
		commonTweetDetailQuerier:    commonTweetDetailQuerier,
		findFavoritesByRangeQuerier: findFavoritesByRangeQuerier,
	}
}

var ErrCommonTweetDetailArrayIsEmpty = fmt.Errorf("common tweet details is empty")

func LastTwID(ds []CommonTweetDetail) (model.TwIDType, error) {
	if len(ds) == 0 {
		return 0, ErrCommonTweetDetailArrayIsEmpty
	}
	//dsはソート済み
	return ds[0].TwID, nil
}

func FirstTwID(ds []CommonTweetDetail) (model.TwIDType, error) {
	if len(ds) == 0 {
		return 0, ErrCommonTweetDetailArrayIsEmpty
	}
	//dsはソート済み
	return ds[len(ds)-1].TwID, nil
}

func (q *TweetDetailsQuerier) FindTweetDetails(requestUserID *model.UsrIDType, limit int, replied_to *model.TwIDType) ([]querier.TweetDetail, error) {
	commonTweetDetails, err := q.commonTweetDetailQuerier.FindCommonTweetDetails(limit, replied_to)
	if err != nil {
		return []querier.TweetDetail{}, err
	}

	lastID, err := LastTwID(commonTweetDetails)
	if err != nil {
		return []querier.TweetDetail{}, err
	}
	firstID, err := FirstTwID(commonTweetDetails)
	if err != nil {
		return []querier.TweetDetail{}, err
	}

	favorites := []model.Favorite{}
	if requestUserID != nil {
		favorites, err = q.findFavoritesByRangeQuerier.FindFavoritesByRange(firstID, lastID, *requestUserID)
		if err != nil {
			return []querier.TweetDetail{}, err
		}
	}

	favIdx := 0
	ret := make([]querier.TweetDetail, 0, len(commonTweetDetails))
	for _, ctd := range commonTweetDetails {
		cur := querier.TweetDetail{
			Tweet:      ctd.Tweet,
			FavCount:   ctd.FavCount,
			ReplyCount: ctd.ReplyCount,
			UserName:   ctd.UserName,
		}

		if favIdx < len(favorites) && favorites[favIdx].TwID == ctd.TwID {
			cur.Favorited = true
			favIdx += 1
		}
		ret = append(ret, cur)
	}

	return ret, nil
}
