package usecase

import (
	"github.com/a-omori-yumemi/YumetterAPI/model"
)

type IAuthenticator interface {
	//returns token
	Login(model.UserName, model.Password) (string, error)
	UseSession(token string) (model.UsrIDType, error)
}
