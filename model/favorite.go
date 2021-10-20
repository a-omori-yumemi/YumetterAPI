package model

import "time"

type Favorite struct {
	TwID      TwIDType  `db:"tw_id" json:"tw_id"`
	UsrID     UsrIDType `db:"usr_id" json:"usr_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
