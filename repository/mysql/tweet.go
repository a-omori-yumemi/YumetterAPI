package mysql

import (
	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
)

type MySQLTweetRepository struct {
	repository.ITweetRepository
	db db.MySQLDB
}

func NewMySQLTweetRepository(DB db.MySQLDB) *MySQLTweetRepository {
	return &MySQLTweetRepository{db: DB}
}

func (r *MySQLTweetRepository) FindTweet(twID model.TwIDType) (tweet model.Tweet, err error) {
	err = r.db.DB.Get(&tweet, "SELECT * FROM Tweet WHERE tw_id=? ORDER BY created_at", twID)
	return tweet, err
}
func (r *MySQLTweetRepository) FindTweets(limit int, replied_to *model.TwIDType) (tweets []model.Tweet, err error) {
	err = r.db.DB.Select(&tweets, "SELECT * FROM Tweet ORDER BY created_at")
	return tweets, err
}
func (r *MySQLTweetRepository) AddTweet(tweet model.Tweet) (ret model.Tweet, err error) {
	res, err := r.db.DB.Exec("INSERT INTO Tweet (body, tw_id, usr_id, replied_to) VALUES (?,?,?,?)", tweet.Body, tweet.TwID, tweet.UsrID, tweet.RepliedTo)
	if err != nil {
		return ret, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		//INSERTには成功しているため、このエラーは握りつぶす（エラーを返したくない）
		tweet.TwID = model.TwIDType(id)
		return tweet, nil
	}
	ret, err = r.FindTweet(model.TwIDType(id))
	if err != nil {
		//INSERTには成功しているため、このエラーは握りつぶす（エラーを返したくない）
		tweet.TwID = model.TwIDType(id)
		return tweet, nil
	}
	return ret, nil
}
func (r MySQLTweetRepository) DeleteTweet(twID model.TwIDType) error {
	_, err := r.db.DB.Exec("DELETE FROM Tweet WHERE tw_id=?", twID)
	return err
}
