package querier_tweet_detail_wire

import (
	querier_tweet_detail "github.com/a-omori-yumemi/YumetterAPI/querier/findTweetDetails"
	querier_tweet_detail_mysql "github.com/a-omori-yumemi/YumetterAPI/querier/findTweetDetails/mysql"
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	querier_tweet_detail_mysql.NewCommonTweetDetailQuerier,
	querier_tweet_detail_mysql.NewFindFavoritesByRangeQuerier,
	querier_tweet_detail.NewCacheStore,

	wire.Bind(new(querier_tweet_detail.ICommonTweetDetailsQuerier), new(*querier_tweet_detail_mysql.CommonTweetDetailsQuerier)),
	wire.Bind(new(querier_tweet_detail.IFindFavoritesByRangeQuerier), new(*querier_tweet_detail_mysql.FindFavoritesByRangeQuerier)),
	wire.Bind(new(querier_tweet_detail.CacheStore), new(*querier_tweet_detail.MemCacheStore)),
)
