package querier

import "github.com/a-omori-yumemi/YumetterAPI/model"

type IUserQuerier interface {
	FindUser(usrID model.UsrIDType) (model.User, error)
	FindUserByName(name model.UserName) (model.User, error)
}
