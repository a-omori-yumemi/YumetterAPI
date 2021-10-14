package mysql

import (
	"github.com/a-omori-yumemi/YumetterAPI/db"
	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
)

type MySQLUserRepository struct {
	repository.IUserRepository
	db db.MySQLDB
}

func NewMySQLUserRepository(DB db.MySQLDB) *MySQLUserRepository {
	return &MySQLUserRepository{db: DB}
}

func (r *MySQLFavoriteRepository) FindUser(usrID model.UsrIDType) (user model.User, err error) {
	err = r.db.DB.Get(&user, "SELECT * FROM User WHERE usr_id=?", usrID)
	return user, err
}
func (r *MySQLFavoriteRepository) AddUser(user model.User) (ret model.User, err error) {
	res, err := r.db.DB.Exec("INSERT INTO User (name, password) VALUES (?,?)", user.Name, user.HashedPassword)
	if err != nil {
		return ret, err
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
func (r *MySQLFavoriteRepository) DeleteUser(usrID model.UsrIDType) error {
	_, err := r.db.DB.Exec("DELETE FROM User WHERE usr_id=?", usrID)
	return err
}
func (r *MySQLFavoriteRepository) UpdateName(usrID model.UsrIDType, name model.UserName) error {
	_, err := r.db.DB.Exec("UPDATE User SET name=? WHERE usr_id=?", name, usrID)
	return err
}
func (r *MySQLFavoriteRepository) UpdatePassword(usrID model.UsrIDType, password model.HashedPassword) error {
	_, err := r.db.DB.Exec("UPDATE User SET password=? WHERE usr_id=?", password, usrID)
	return err
}
