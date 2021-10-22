package repository_wire

import (
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	repo_mysql "github.com/a-omori-yumemi/YumetterAPI/repository/mysql"
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	repo_mysql.NewMySQLFavoriteRepository,
	repo_mysql.NewMySQLUserRepository,
	repo_mysql.NewMySQLTweetRepository,

	wire.Bind(new(repository.IFavoriteRepository), new(*repo_mysql.MySQLFavoriteRepository)),
	wire.Bind(new(repository.IUserRepository), new(*repo_mysql.MySQLUserRepository)),
	wire.Bind(new(repository.ITweetRepository), new(*repo_mysql.MySQLTweetRepository)),
)
