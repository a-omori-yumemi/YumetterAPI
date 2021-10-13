package model

import "time"

type Favorite struct {
	TwID      TwIDType  `db:"tw_id"`
	UsrID     UsrIDType `db:"usr_id"`
	CreatedAt time.Time `db:"created_at"`
}
