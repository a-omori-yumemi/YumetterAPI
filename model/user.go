package model

import (
	"fmt"
	"time"
)

type UsrIDType int32
type UserName string

type User struct {
	UsrID     UsrIDType `db:"usr_id" json:"usr_id"`
	Name      UserName  `db:"name" json:"name"`
	Password  []byte    `db:"password"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpadtedAt time.Time `db:"updated_at" json:"updated_at"`
}

const UserNameMaxLength = 64
const UserNameMinLength = 1

func (u User) Validate() error {
	return u.Name.Validate()
}

func (n UserName) Validate() error {
	if len(n) < UserNameMinLength || len(n) > UserNameMaxLength {
		return fmt.Errorf("user name is too long")
	}
	return nil
}
