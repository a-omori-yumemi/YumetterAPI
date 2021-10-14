package service

import (
	"database/sql"
	"fmt"

	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"go.uber.org/multierr"
)

type TweetService struct {
	ITweetService
	favRepo   repository.IFavoriteRepository
	tweetRepo repository.ITweetRepository
	userRepo  repository.IUserRepository
}

func NewTweetService(
	favRepo repository.IFavoriteRepository,
	tweetRepo repository.ITweetRepository,
	userRepo repository.IUserRepository) *TweetService {
	return &TweetService{favRepo: favRepo, tweetRepo: tweetRepo, userRepo: userRepo}
}

var ErrForbidden error = fmt.Errorf("forbidden operation")

func (s *TweetService) DeleteTweetWithAuth(requestUserID model.UsrIDType, twID model.TwIDType) error {
	tw, err := s.FindTweet(twID)
	if err != nil {
		return err
	}

	if tw.UsrID != requestUserID {
		return ErrForbidden
	}
	return nil
}

func (s *TweetService) FindTweetDetails(
	requestUserID *model.UsrIDType,
	limit int,
	replied_to *model.TwIDType) ([]TweetDetail, error) {

	tweets, err := s.FindTweets(limit, replied_to)
	if err != nil {
		return []TweetDetail{}, err
	}

	tweetDetails := make([]TweetDetail, len(tweets))
	// キャッシュで解決しよう
	for _, tweet := range tweets {
		user, err1 := s.userRepo.FindUser(tweet.UsrID)
		favs, err2 := s.favRepo.FindFavorites(tweet.TwID)
		replies, err3 := s.FindTweets(10000000, tweet.RepliedTo)
		favorited := false
		var err4 error
		if requestUserID != nil {
			_, err4 = s.favRepo.FindFavorite(tweet.TwID, *requestUserID)
			favorited = err == sql.ErrNoRows
		}
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return tweetDetails, multierr.Combine(err1, err2, err3, err4)
		}

		tweetDetails = append(tweetDetails, TweetDetail{
			UserName:   user.Name,
			Tweet:      tweet,
			FavCount:   len(favs),
			ReplyCount: len(replies),
			Favorited:  favorited,
		})
	}

	return tweetDetails, nil
}

func (s *TweetService) FindTweet(twID model.TwIDType) (model.Tweet, error) {
	return s.tweetRepo.FindTweet(twID)
}
func (s *TweetService) FindTweets(limit int, replied_to *model.TwIDType) ([]model.Tweet, error) {
	return s.tweetRepo.FindTweets(limit, replied_to)
}
func (s *TweetService) AddTweet(tweet model.Tweet) (model.Tweet, error) {
	return s.tweetRepo.AddTweet(tweet)
}
func (s *TweetService) DeleteTweet(twID model.TwIDType) error {
	return s.tweetRepo.DeleteTweet(twID)
}
