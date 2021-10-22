package querier

type Queriers struct {
	TweetDetailQuerier IFindTweetDetailsQuerier
	TweetQuerier       ITweetQuerier
	UserQuerier        IUserQuerier
	FavQuerier         IFavoriteQuerier
}
