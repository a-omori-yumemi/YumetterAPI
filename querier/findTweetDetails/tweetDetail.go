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
	UserName model.UserName
	model.Tweet
	Favorites  map[model.UsrIDType]bool
	ReplyCount int
}

type IFindTweetDetailsQuerier interface {
	FindTweetDetails(
		requestUserID *model.UsrIDType,
		limit int,
		replied_to *model.TwIDType) ([]querier.TweetDetail, error)
}

type TweetDetailsQuerier struct {
	dataSource               data_source_wrapper.DataSourceWrapper
	commonTweetDetailQuerier ICommonTweetDetailsQuerier
	findTweetDetailsQuerier  IFindTweetDetailsQuerier
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
	dataSourceMaker CommonTweetDetailDataSourceMaker,
	findTweetDetailsQuerier IFindTweetDetailsQuerier) *TweetDetailsQuerier {

	dataSource := dataSourceMaker.NewDataSourceWrapper(func() (interface{}, error) {
		return commonTweetDetailQuerier.FindCommonTweetDetails(querier.MaxLimitValue, nil)
	})

	return &TweetDetailsQuerier{
		dataSource:               dataSource,
		commonTweetDetailQuerier: commonTweetDetailQuerier,
		findTweetDetailsQuerier:  findTweetDetailsQuerier,
	}
}

func (q *TweetDetailsQuerier) FindTweetDetails(requestUserID *model.UsrIDType, limit int, replied_to *model.TwIDType) ([]querier.TweetDetail, error) {
	var commonTweetDetails []CommonTweetDetail
	if replied_to == nil {
		var ok bool
		commonTweetDetails, ok = q.dataSource.Get().([]CommonTweetDetail)
		if !ok {
			return []querier.TweetDetail{}, fmt.Errorf("failed to get TimeLine")
		}
	} else {
		return q.findTweetDetailsQuerier.FindTweetDetails(requestUserID, limit, replied_to)
	}

	if limit > len(commonTweetDetails) {
		limit = len(commonTweetDetails)
	}

	ret := make([]querier.TweetDetail, 0, limit)
	for i := 0; limit > i; i++ {
		favorited := false
		if requestUserID != nil {
			var ok bool
			if favorited, ok = commonTweetDetails[i].Favorites[*requestUserID]; !ok {
				favorited = false
			}
		}

		cur := querier.TweetDetail{
			Tweet:      commonTweetDetails[i].Tweet,
			FavCount:   len(commonTweetDetails[i].Favorites),
			Favorited:  favorited,
			ReplyCount: commonTweetDetails[i].ReplyCount,
			UserName:   commonTweetDetails[i].UserName,
		}

		ret = append(ret, cur)
	}

	return ret, nil
}
