package usecase_wire

import (
	"github.com/a-omori-yumemi/YumetterAPI/usecase"
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	usecase.NewTweetDeleteUsecase,
	usecase.NewJWTAuthenticator,

	wire.Bind(new(usecase.IAuthenticator), new(*usecase.JWTAuthenticator)),
	wire.Bind(new(usecase.ITweetDeleteUsecase), new(*usecase.TweetDeleteUsecase)),
)
