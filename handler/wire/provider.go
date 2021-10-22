package handler_wire

import (
	"github.com/a-omori-yumemi/YumetterAPI/handler"
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	wire.Struct(new(handler.AuthUserMiddleware), "*"),
	wire.Struct(new(handler.GetTweet), "*"),
	wire.Struct(new(handler.DeleteTweet), "*"),
	wire.Struct(new(handler.GetTweets), "*"),
	wire.Struct(new(handler.PostTweet), "*"),
	wire.Struct(new(handler.GetUser), "*"),
	wire.Struct(new(handler.RegisterUser), "*"),
	wire.Struct(new(handler.LoginUser), "*"),
	wire.Struct(new(handler.GetMe), "*"),
	wire.Struct(new(handler.DeleteMe), "*"),
	wire.Struct(new(handler.PatchMe), "*"),
	wire.Struct(new(handler.GetFavorites), "*"),
	wire.Struct(new(handler.PutFavorite), "*"),
	wire.Struct(new(handler.DeleteFavorite), "*"),
)
