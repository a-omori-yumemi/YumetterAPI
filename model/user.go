package model

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UsrIDType int32
type UserName string
type Password string
type HashedPassword []byte

type User struct {
	UsrID          UsrIDType      `db:"usr_id" json:"usr_id"`
	Name           UserName       `db:"name" json:"name"`
	HashedPassword HashedPassword `db:"password"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpadtedAt      time.Time      `db:"updated_at" json:"updated_at"`
}

func (u User) Validate() error {
	return u.Name.Validate()
}

const UserNameMaxLength = 64
const UserNameMinLength = 1

func (n UserName) Validate() error {
	if len(n) < UserNameMinLength {
		return fmt.Errorf("user name is too short")
	}
	if len(n) > UserNameMaxLength {
		return fmt.Errorf("user name is too long")
	}
	return nil
}

const PasswordMaxLength = 30
const PasswordMinLength = 8

func (p Password) Validate() error {
	if len(p) < PasswordMinLength {
		return fmt.Errorf("password is too short")
	}
	if len(p) > PasswordMaxLength {
		return fmt.Errorf("password is too long")
	}
	return nil
}

func (p Password) Hash() (HashedPassword, error) {
	return bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
}
