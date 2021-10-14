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
	err = r.db.DB.Get(&tweet, "SELECT * FROM Tweet WHERE tw_id=? ORDER BY created_at DESC", twID)
	return tweet, interpretMySQLError(err)
}
func (r *MySQLTweetRepository) FindTweets(limit int, replied_to *model.TwIDType) (tweets []model.Tweet, err error) {
	tweets = []model.Tweet{}
	if replied_to == nil {
		err = r.db.DB.Select(&tweets, "SELECT * FROM Tweet ORDER BY created_at DESC LIMIT ?", limit)
	} else {
		err = r.db.DB.Select(&tweets, "SELECT * FROM Tweet WHERE replied_to=? ORDER BY created_at DESC LIMIT ?", replied_to, limit)
	}
	return tweets, interpretMySQLError(err)
}
func (r *MySQLTweetRepository) AddTweet(tweet model.Tweet) (ret model.Tweet, err error) {
	res, err := r.db.DB.Exec("INSERT INTO Tweet (body, usr_id, replied_to) VALUES (?,?,?)", tweet.Body, tweet.UsrID, tweet.RepliedTo)
	if err != nil {
		return ret, interpretMySQLError(err)
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
func (r *MySQLTweetRepository) DeleteTweet(twID model.TwIDType) error {
	res, err := r.db.DB.Exec("DELETE FROM Tweet WHERE tw_id=?", twID)
	if err != nil {
		return interpretMySQLError(err)
	}
	if cou, err := res.RowsAffected(); err != nil && cou == 0 {
		return repository.ErrNotFound
	}
	return nil
}
