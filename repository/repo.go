package repository

type Repositories struct {
	TweetRepo ITweetRepository
	UserRepo  IUserRepository
	FavRepo   IFavoriteRepositoty
}
