package handler

import (
	"strconv"

	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"github.com/a-omori-yumemi/YumetterAPI/service"
	"github.com/labstack/echo/v4"
)

func PostTweet(tweetRepo repository.ITweetRepository) echo.HandlerFunc {

	GetParams := func(c echo.Context) (string, *model.TwIDType, error) {
		body := c.FormValue("body")

		var repliedTo *model.TwIDType = nil
		if repliedToStr := c.FormValue("replied_to"); repliedToStr != "" {
			if repliedToTmp, err := strconv.Atoi(repliedToStr); err == nil {
				repliedTo = new(model.TwIDType)
				*repliedTo = model.TwIDType(repliedToTmp)
			} else {
				return body, repliedTo, echo.NewHTTPError(400, "invalid parameter: {replied_to}")
			}
		}

		return body, repliedTo, nil
	}

	return func(c echo.Context) error {
		body, repliedTo, err := GetParams(c)
		if err != nil {
			return err
		}
		tweet := model.Tweet{
			Body:      body,
			RepliedTo: repliedTo,
		}
		if tweet.Validate() != nil {
			return echo.NewHTTPError(400, "body is too short or the tweet you replied to is missing")
		}

		_, err = tweetRepo.AddTweet(tweet)
		return err
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

func GetTweets(tweetService service.ITweetService) echo.HandlerFunc {
	const DefaultLimitValue = 30

	GetParams := func(c echo.Context) (*model.UsrIDType, int, *model.TwIDType, error) {
		var usrID *model.UsrIDType = nil
		usrIDTmp, ok := c.Get(ContextUserKey).(string)
		if ok {
			if usrIDTmp, err := strconv.Atoi(usrIDTmp); err == nil {
				usrID = new(model.UsrIDType)
				*usrID = model.UsrIDType(usrIDTmp)
			}
		}

		limit, err := strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			limit = DefaultLimitValue
		}

		var repliedTo *model.TwIDType = nil
		if repliedToTmp, err := strconv.Atoi(c.QueryParam("replied_to")); err == nil {
			repliedTo = new(model.TwIDType)
			*repliedTo = model.TwIDType(repliedToTmp)
		}

		return usrID, limit, repliedTo, nil
	}

	return func(c echo.Context) error {
		usrID, limit, repliedTo, err := GetParams(c)
		if err != nil {
			return err
		}

		tweet, err := tweetService.FindTweetDetails(usrID, limit, repliedTo)
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
