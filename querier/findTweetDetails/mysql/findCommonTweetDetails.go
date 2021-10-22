package querier_tweet_detail_mysql

import (
	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/querier"
)

type CommonTweetDetailQuerier struct {
	db db.MySQLReadOnlyDB
}

func NewCommonTweetDetailQuerier(db db.MySQLReadOnlyDB) *CommonTweetDetailQuerier {
	return &CommonTweetDetailQuerier{db: db}
}

func (u *CommonTweetDetailQuerier) FindCommonTweetDetails(
	requestUserID *model.UsrIDType,
	limit int,
	replied_to *model.TwIDType) ([]querier.TweetDetail, error) {

	type FlattenCommonTweetDetail struct {
		UserName model.UserName `db:"user_name"`
		model.Tweet
		FavCount   int `db:"fav_count"`
		ReplyCount int `db:"reply_count"`
	}
	fTweetDetails := make([]FlattenCommonTweetDetail, 0, limit)

	whereClause := ""
	args := []interface{}{}
	args = append(args, requestUserID)
	if replied_to != nil {
		whereClause = `WHERE replied_to=?`
		args = append(args, replied_to)
	}
	args = append(args, limit)

	err := u.db.DB.Select(&fTweetDetails,
		`SELECT 
		 U.name user_name,
		 (SELECT count(1) FROM Favorite WHERE tw_id=T.tw_id) fav_count,
		 (SELECT count(1) FROM Tweet WHERE replied_to=T.tw_id) reply_count,
		 T.*
		 FROM Tweet T JOIN User U USING(usr_id) `+
			whereClause+
			` ORDER BY tw_id DESC limit ?`,
		args...,
	)

	tweetDetails := make([]querier.TweetDetail, 0, len(fTweetDetails))
	for _, d := range fTweetDetails {
		tweetDetails = append(tweetDetails,
			querier.TweetDetail{
				UserName:   d.UserName,
				Tweet:      d.Tweet,
				FavCount:   d.FavCount,
				ReplyCount: d.ReplyCount,
			})
	}
	return tweetDetails, err
}

// SELECT (SELECT count(1) FROM Favorite WHERE tw_id=T.tw_id) fav_count,
// (SELECT count(1) FROM Tweet WHERE replied_to=T.tw_id) reply_count,
// (F.usr_id is not NULL) favorited , F.usr_id, T.*
// FROM Tweet T LEFT OUTER JOIN
// (SELECT * FROM Favorite WHERE usr_id=25 ORDER BY tw_id DESC limit 30) F USING(tw_id)
// ORDER BY tw_id DESC limit 30
