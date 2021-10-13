package model

import (
	"time"
)

type TwIDType uint64

type Tweet struct {
	Body      string    `db:"body"`
	TwID      TwIDType  `db:"tw_id"`
	UsrID     UsrIDType `db:"usr_id"`
	CreatedAt time.Time `db:"craeted_at"`
	RepliedTo *TwIDType `db:"replied_to"`
}
