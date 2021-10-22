package handler

type Handlers struct {
	AuthUserMiddleware AuthUserMiddleware
	GetTweet           GetTweet
	DeleteTweet        DeleteTweet
	GetTweets          GetTweets
	PostTweet          PostTweet
	GetUser            GetUser
	RegisterUser       RegisterUser
	LoginUser          LoginUser
	GetMe              GetMe
	DeleteMe           DeleteMe
	PatchMe            PatchMe
	GetFavorites       GetFavorites
	PutFavorite        PutFavorite
	DeleteFavorite     DeleteFavorite
}
