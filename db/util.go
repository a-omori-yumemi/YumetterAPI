package db

import (
	"database/sql"

	"github.com/VividCortex/mysqlerr"
	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/go-sql-driver/mysql"
)

func InterpretMySQLError(err error) error {
	if err != nil {
		if err == sql.ErrNoRows {
			return model.ErrNotFound
		}
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			switch driverErr.Number {
			case mysqlerr.ER_DUP_ENTRY:
				return model.ErrDuplicateKey
			}
		}
	}
	return err
}
