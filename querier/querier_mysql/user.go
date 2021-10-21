package querier_mysql

import (
	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/model"
)

type MySQLUserQuerier struct {
	db db.MySQLReadOnlyDB
}

func NewMySQLUserQuerier(DB db.MySQLReadOnlyDB) *MySQLUserQuerier {
	return &MySQLUserQuerier{db: DB}
}

func (r *MySQLUserQuerier) FindUser(usrID model.UsrIDType) (user model.User, err error) {
	err = r.db.DB.Get(&user, "SELECT * FROM User WHERE usr_id=?", usrID)
	return user, db.InterpretMySQLError(err)
}
func (r *MySQLUserQuerier) FindUserByName(name model.UserName) (user model.User, err error) {
	err = r.db.DB.Get(&user, "SELECT * FROM User WHERE name=?", name)
	return user, db.InterpretMySQLError(err)
}
