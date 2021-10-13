package model

import (
	"fmt"
	"time"
)

type TwIDType uint64

type Tweet struct {
	Body      string    `db:"body" json:"body"`
	TwID      TwIDType  `db:"tw_id" json:"tw_id"`
	UsrID     UsrIDType `db:"usr_id" json:"usr_id"`
	CreatedAt time.Time `db:"craeted_at" json:"created_at"`
	RepliedTo *TwIDType `db:"replied_to" json:"replied_to,omitempty"`
}

const BodyMaxLength = 280
const BodyMinLength = 1

func (t Tweet) Validate() error {
	if len(t.Body) < BodyMinLength {
		return fmt.Errorf("body is too short")
	}
	if len(t.Body) > BodyMaxLength {
		return fmt.Errorf("body is too long")
	}
	return nil
}
