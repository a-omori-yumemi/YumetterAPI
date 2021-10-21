package usecase

type Usecases struct {
	Authenticator IAuthenticator
	TweetService  ITweetService
}
