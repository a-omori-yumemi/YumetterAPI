package usecase

import (
	"fmt"

	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
)

type ITweetDeleteUsecase interface {
	DeleteTweetWithAuth(requestUserID model.UsrIDType, twID model.TwIDType) error
}

type TweetDeleteUsecase struct {
	repo repository.ITweetRepository
}

func NewTweetDeleteUsecase(repo repository.ITweetRepository) *TweetDeleteUsecase {
	return &TweetDeleteUsecase{repo: repo}
}

var ErrForbidden error = fmt.Errorf("forbidden operation")

func (u *TweetDeleteUsecase) DeleteTweetWithAuth(requestUserID model.UsrIDType, twID model.TwIDType) error {
	tw, err := u.repo.FindTweet(twID)
	if err != nil {
		return err
	}

	if tw.UsrID != requestUserID {
		return ErrForbidden
	}

	err = u.repo.DeleteTweet(twID)
	return err
}
