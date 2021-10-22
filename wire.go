//go:build wireinject
// +build wireinject

package main

import (
	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/handler"
	handler_wire "github.com/a-omori-yumemi/YumetterAPI/handler/wire"
	querier_wire "github.com/a-omori-yumemi/YumetterAPI/querier/wire"
	repository_wire "github.com/a-omori-yumemi/YumetterAPI/repository/wire"
	usecase_wire "github.com/a-omori-yumemi/YumetterAPI/usecase/wire"

	"github.com/google/wire"
)

func buildHandlers() (handler.Handlers, error) {
	wire.Build(
		ProvideDBConfig,
		ProvideRODBConfig,
		ProvideSecretKey,
		db.NewMySQLDB,
		db.NewMySQLReadOnlyDB,
		handler_wire.SuperSet,
		repository_wire.SuperSet,
		usecase_wire.SuperSet,
		querier_wire.SuperSet,
		wire.Struct(new(handler.Handlers), "*"),
	)
	return handler.Handlers{}, nil
}
