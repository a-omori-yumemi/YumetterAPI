package model

import (
	"time"
)

type UsrIDType int32

type User struct {
	UsrID     UsrIDType `db:"usr_id"`
	Name      string    `db:"name"`
	Password  []byte    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpadtedAt time.Time `db:"updated_at"`
}
