package usecase

import (
	"fmt"
	"time"

	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/querier"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type JWTAuthenticator struct {
	IAuthenticator
	userQuerier querier.IUserQuerier
	secretKey   string
}

type Claim struct {
	jwt.StandardClaims
	UsrID model.UsrIDType
}

func NewJWTAuthenticator(userQuerier querier.IUserQuerier, secretKey string) *JWTAuthenticator {
	return &JWTAuthenticator{userQuerier: userQuerier, secretKey: secretKey}
}

func (a *JWTAuthenticator) Login(name model.UserName, password model.Password) (string, error) {
	user, err := a.userQuerier.FindUserByName(name)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err != nil {
		return "", fmt.Errorf("failed to login")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	expire := time.Now().Add(time.Hour * 24 * 3)
	token.Claims = &Claim{
		jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.UsrID,
	}
	return token.SignedString([]byte(a.secretKey))
}

func (a *JWTAuthenticator) UseSession(tokenString string) (model.UsrIDType, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("invalid session")
	}
	if claim, ok := token.Claims.(*Claim); ok {
		return claim.UsrID, nil
	} else {
		return 0, fmt.Errorf("invalid session")
	}
}
