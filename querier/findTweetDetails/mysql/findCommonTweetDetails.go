package querier_tweet_detail_mysql

import (
	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/model"

	querier_tweet_detail "github.com/a-omori-yumemi/YumetterAPI/querier/findTweetDetails"
)

type CommonTweetDetailsQuerier struct {
	db db.MySQLReadOnlyDB
}

func NewCommonTweetDetailQuerier(db db.MySQLReadOnlyDB) *CommonTweetDetailsQuerier {
	return &CommonTweetDetailsQuerier{db: db}
}

func (u *CommonTweetDetailsQuerier) FindCommonTweetDetails(
	limit int,
	replied_to *model.TwIDType) ([]querier_tweet_detail.CommonTweetDetail, error) {

	type FlattenCommonTweetDetail struct {
		UserName model.UserName `db:"user_name"`
		model.Tweet
		ReplyCount int `db:"reply_count"`
	}
	fTweetDetails := make([]FlattenCommonTweetDetail, 0, limit)

	whereClause := ""
	args := []interface{}{}
	if replied_to != nil {
		whereClause = `WHERE replied_to=?`
		args = append(args, replied_to)
	}
	args = append(args, limit)

	err := u.db.DB.Select(&fTweetDetails,
		`SELECT 
		 U.name user_name,
		 (SELECT count(1) FROM Tweet WHERE replied_to=T.tw_id) reply_count,
		 T.*
		 FROM Tweet T JOIN User U USING(usr_id) `+
			whereClause+
			` ORDER BY tw_id DESC limit ?`,
		args...,
	)
	if err != nil || len(fTweetDetails) == 0 {
		return []querier_tweet_detail.CommonTweetDetail{}, err
	}

	favorites := []model.Favorite{}

	firstTwID := fTweetDetails[len(fTweetDetails)-1].TwID
	lastTwID := fTweetDetails[0].TwID

	err = u.db.DB.Select(&favorites,
		`SELECT tw_id, usr_id FROM Favorite WHERE tw_id BETWEEN ? AND ? ORDER BY tw_id DESC`,
		firstTwID, lastTwID,
	)
	if err != nil {
		return []querier_tweet_detail.CommonTweetDetail{}, err
	}

	tweetDetails := make([]querier_tweet_detail.CommonTweetDetail, 0, len(fTweetDetails))
	favoritesIdx := 0
	for _, d := range fTweetDetails {
		favoritesToThisTweet := map[model.UsrIDType]bool{}
		for len(favorites) > favoritesIdx && favorites[favoritesIdx].TwID == d.TwID {
			favoritesToThisTweet[favorites[favoritesIdx].UsrID] = true
			favoritesIdx++
		}

		tweetDetails = append(tweetDetails,
			querier_tweet_detail.CommonTweetDetail{
				UserName:   d.UserName,
				Favorites:  favoritesToThisTweet,
				Tweet:      d.Tweet,
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
