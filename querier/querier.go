package querier

type Queriers struct {
	TweetDetailQuerier ITweetDetailQuerier
	TweetQuerier       ITweetQuerier
	UserQuerier        IUserQuerier
	FavQuerier         IFavoriteQuerier
}
