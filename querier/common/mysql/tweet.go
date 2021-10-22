package querier_mysql

import (
	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/model"
)

type MySQLTweetQuerier struct {
	db db.MySQLReadOnlyDB
}

func NewMySQLTweetQuerier(DB db.MySQLReadOnlyDB) *MySQLTweetQuerier {
	return &MySQLTweetQuerier{db: DB}
}

func (r *MySQLTweetQuerier) FindTweet(twID model.TwIDType) (tweet model.Tweet, err error) {
	err = r.db.DB.Get(&tweet, "SELECT * FROM Tweet WHERE tw_id=? ORDER BY created_at DESC", twID)
	return tweet, db.InterpretMySQLError(err)
}
func (r *MySQLTweetQuerier) FindTweets(limit int, replied_to *model.TwIDType) (tweets []model.Tweet, err error) {
	tweets = []model.Tweet{}
	if replied_to == nil {
		err = r.db.DB.Select(&tweets, "SELECT * FROM Tweet ORDER BY created_at DESC LIMIT ?", limit)
	} else {
		err = r.db.DB.Select(&tweets, "SELECT * FROM Tweet WHERE replied_to=? ORDER BY created_at DESC LIMIT ?", replied_to, limit)
	}
	return tweets, db.InterpretMySQLError(err)
}
