package mysql

import (
	"database/sql"

	"github.com/VividCortex/mysqlerr"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"github.com/go-sql-driver/mysql"
)

func interpretMySQLError(err error) error {
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.ErrNotFound
		}
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			switch driverErr.Number {
			case mysqlerr.ER_DUP_ENTRY:
				return repository.ErrDuplicateKey
			}
		}
	}
	return err
}
