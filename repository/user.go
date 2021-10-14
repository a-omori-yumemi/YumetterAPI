package repository

import "github.com/a-omori-yumemi/YumetterAPI/model"

type IUserRepository interface {
	FindUser(usrID model.UsrIDType) (model.User, error)
	AddUser(user model.User) (model.User, error)
	DeleteUser(usrID model.UsrIDType) error
	UpdateName(usrID model.UsrIDType, name model.UserName) error
	UpdatePassword(usrID model.UsrIDType, password model.HashedPassword) error
}
