package handler

import (
	"strconv"

	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"github.com/labstack/echo/v4"
)

func PostTweet(tweetRepo repository.ITweetRepository) echo.HandlerFunc {

	return func(c echo.Context) error {
		return nil
	}
}

func GetTweet(tweetRepo repository.ITweetRepository) echo.HandlerFunc {

	return func(c echo.Context) error {
		tw_id, err := strconv.Atoi(c.Param("tw_id"))
		if err != nil {
			return err
		}

		tweet, err := tweetRepo.FindTweet(model.TwIDType(tw_id))
		if err != nil {
			return err
		}

		return c.JSON(200, tweet)
	}
}

func DeleteTweet(tweetRepo repository.ITweetRepository) echo.HandlerFunc {

	return func(c echo.Context) error {
		tw_id, err := strconv.Atoi(c.Param("tw_id"))
		if err != nil {
			return err
		}

		err = tweetRepo.DeleteTweet(model.TwIDType(tw_id))
		if err != nil {
			return err
		}

		return c.NoContent(204)
	}
}
