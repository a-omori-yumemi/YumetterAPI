package querier_wire

import (
	"github.com/a-omori-yumemi/YumetterAPI/querier"
	querier_tweet_detail "github.com/a-omori-yumemi/YumetterAPI/querier/findTweetDetails"
	querier_tweet_detail_wire "github.com/a-omori-yumemi/YumetterAPI/querier/findTweetDetails/wire"

	querier_mysql "github.com/a-omori-yumemi/YumetterAPI/querier/common/mysql"
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	querier_mysql.NewMySQLFavoriteQuerier,
	querier_mysql.NewMySQLTweetQuerier,
	querier_tweet_detail.NewTweetDetailQuerier,
	querier_mysql.NewMySQLUserQuerier,

	wire.Bind(new(querier.IFavoriteQuerier), new(*querier_mysql.MySQLFavoriteQuerier)),
	wire.Bind(new(querier.IUserQuerier), new(*querier_mysql.MySQLUserQuerier)),
	wire.Bind(new(querier.ITweetQuerier), new(*querier_mysql.MySQLTweetQuerier)),
	wire.Bind(new(querier.IFindTweetDetailsQuerier), new(*querier_tweet_detail.TweetDetailsQuerier)),
	querier_tweet_detail_wire.SuperSet,
)
