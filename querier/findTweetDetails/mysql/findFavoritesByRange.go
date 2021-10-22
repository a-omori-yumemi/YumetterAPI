package querier_tweet_detail_mysql

import (
	"math"

	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/model"
)

type FindFavoritesByRangeQuerier struct {
	db db.MySQLReadOnlyDB
}

func NewFindFavoritesByRangeQuerier(db db.MySQLReadOnlyDB) *FindFavoritesByRangeQuerier {
	return &FindFavoritesByRangeQuerier{
		db: db,
	}
}

func (q *FindFavoritesByRangeQuerier) FindFavoritesByRange(
	firstTwID model.TwIDType,
	lastTwID model.TwIDType,
	usrID model.UsrIDType) ([]model.Favorite, error) {

	favorites := make([]model.Favorite, 0, int(math.Max(float64(lastTwID-firstTwID), 0)))
	err := q.db.DB.Select(&favorites,
		"SELECT tw_id, usr_id FROM Favorite WHERE tw_id BETWEEN ? AND ? AND usr_id=? ORDER BY tw_id DESC",
		firstTwID,
		lastTwID,
		usrID,
	)
	return favorites, db.InterpretMySQLError(err)
}
