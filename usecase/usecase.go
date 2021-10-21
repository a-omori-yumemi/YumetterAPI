package usecase

type Usecases struct {
	Authenticator      IAuthenticator
	TweetDeleteUsecase ITweetDeleteUsecase
	TweetDetailUsecase ITweetDetailQuerier
}
