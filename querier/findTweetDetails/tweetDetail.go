package querier_tweet_detail

import (
	"fmt"
	"time"

	data_source_wrapper "github.com/a-omori-yumemi/YumetterAPI/dataSourceWrapper"
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
	dataSource                  data_source_wrapper.DataSourceWrapper
	commonTweetDetailQuerier    ICommonTweetDetailsQuerier
	findFavoritesByRangeQuerier IFindFavoritesByRangeQuerier
}

type TweetDetailsCacheLifeTime int

type CommonTweetDetailDataSourceMaker interface {
	data_source_wrapper.DataSourceMaker
}

func NewCommonTweetDetailDataSourceMaker(lifeTime TweetDetailsCacheLifeTime) CommonTweetDetailDataSourceMaker {
	return data_source_wrapper.NewCacheMaker(time.Duration(lifeTime) * time.Second)
}

func NewTweetDetailQuerier(
	commonTweetDetailQuerier ICommonTweetDetailsQuerier,
	findFavoritesByRangeQuerier IFindFavoritesByRangeQuerier,
	dataSourceMaker CommonTweetDetailDataSourceMaker) *TweetDetailsQuerier {

	dataSource := dataSourceMaker.NewDataSourceWrapper(func() (interface{}, error) {
		return commonTweetDetailQuerier.FindCommonTweetDetails(querier.MaxLimitValue, nil)
	})

	return &TweetDetailsQuerier{
		dataSource:                  dataSource,
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

func FirstTwID(ds []CommonTweetDetail, limit int) (model.TwIDType, error) {
	if len(ds) == 0 || limit == 0 {
		return 0, ErrCommonTweetDetailArrayIsEmpty
	}
	//dsはソート済み
	return ds[limit-1].TwID, nil
}

func (q *TweetDetailsQuerier) FindTweetDetails(requestUserID *model.UsrIDType, limit int, replied_to *model.TwIDType) ([]querier.TweetDetail, error) {
	commonTweetDetails, ok := q.dataSource.Get().([]CommonTweetDetail)
	if !ok {
		return []querier.TweetDetail{}, fmt.Errorf("failed to get TimeLine")
	}

	if limit > len(commonTweetDetails) {
		limit = len(commonTweetDetails)
	}

	lastID, err := LastTwID(commonTweetDetails)
	if err != nil {
		return []querier.TweetDetail{}, err
	}
	firstID, err := FirstTwID(commonTweetDetails, limit)
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
	ret := make([]querier.TweetDetail, 0, limit)
	for i := 0; limit > i; i++ {
		cur := querier.TweetDetail{
			Tweet:      commonTweetDetails[i].Tweet,
			FavCount:   commonTweetDetails[i].FavCount,
			ReplyCount: commonTweetDetails[i].ReplyCount,
			UserName:   commonTweetDetails[i].UserName,
		}

		if favIdx < len(favorites) && favorites[favIdx].TwID == commonTweetDetails[i].TwID {
			cur.Favorited = true
			favIdx += 1
		}
		ret = append(ret, cur)
	}

	return ret, nil
}
