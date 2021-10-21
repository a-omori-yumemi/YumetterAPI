package repo_mysql

import (
	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/model"
)

type MySQLUserRepository struct {
	db db.MySQLDB
}

func NewMySQLUserRepository(DB db.MySQLDB) *MySQLUserRepository {
	return &MySQLUserRepository{db: DB}
}

func (r *MySQLUserRepository) FindUser(usrID model.UsrIDType) (user model.User, err error) {
	err = r.db.DB.Get(&user, "SELECT * FROM User WHERE usr_id=?", usrID)
	return user, db.InterpretMySQLError(err)
}
func (r *MySQLUserRepository) FindUserByName(name model.UserName) (user model.User, err error) {
	err = r.db.DB.Get(&user, "SELECT * FROM User WHERE name=?", name)
	return user, db.InterpretMySQLError(err)
}
func (r *MySQLUserRepository) AddUser(user model.User) (ret model.User, err error) {
	res, err := r.db.DB.Exec("INSERT INTO User (name, password) VALUES (?,?)", user.Name, user.HashedPassword)
	if err != nil {
		return ret, db.InterpretMySQLError(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		//INSERTには成功しているため、このエラーは握りつぶす（エラーを返したくない）
		user.UsrID = model.UsrIDType(id)
		return user, nil
	}
	ret, err = r.FindUser(model.UsrIDType(id))
	if err != nil {
		//INSERTには成功しているため、このエラーは握りつぶす（エラーを返したくない）
		user.UsrID = model.UsrIDType(id)
		return user, nil
	}
	return ret, nil
}
func (r *MySQLUserRepository) DeleteUser(usrID model.UsrIDType) error {
	res, err := r.db.DB.Exec("DELETE FROM User WHERE usr_id=?", usrID)
	if err != nil {
		return db.InterpretMySQLError(err)
	}
	if cou, err := res.RowsAffected(); err != nil && cou == 0 {
		return model.ErrNotFound
	}
	return nil
}
func (r *MySQLUserRepository) UpdateName(usrID model.UsrIDType, name model.UserName) error {
	res, err := r.db.DB.Exec("UPDATE User SET name=? WHERE usr_id=?", name, usrID)
	if err != nil {
		return db.InterpretMySQLError(err)
	}
	if cou, err := res.RowsAffected(); err != nil && cou == 0 {
		return model.ErrNotFound
	}
	return nil
}
func (r *MySQLUserRepository) UpdatePassword(usrID model.UsrIDType, password model.HashedPassword) error {
	res, err := r.db.DB.Exec("UPDATE User SET password=? WHERE usr_id=?", password, usrID)
	if err != nil {
		return db.InterpretMySQLError(err)
	}
	if cou, err := res.RowsAffected(); err != nil && cou == 0 {
		return model.ErrNotFound
	}
	return nil
}
