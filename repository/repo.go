package repository

import (
	"fmt"
)

type Repositories struct {
	TweetRepo ITweetRepository
	UserRepo  IUserRepository
	FavRepo   IFavoriteRepository
}

var ErrDuplicateKey = fmt.Errorf("duplicate key")
var ErrNotFound = fmt.Errorf("not found")
