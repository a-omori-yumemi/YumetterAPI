package model

import (
	"time"
)

type TwIDType int32

type Tweet struct {
	Body      string    `db:"body"`
	TwID      TwIDType  `db:"tw_id"`
	UsrID     UsrIDType `db:"usr_id"`
	CreatedAt time.Time `db:"craeted_at"`
	RepliedTo *TwIDType `db:"replied_to"`
}
