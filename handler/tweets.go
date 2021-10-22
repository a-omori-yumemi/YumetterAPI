package handler

import (
	"strconv"

	"github.com/a-omori-yumemi/YumetterAPI/model"
	"github.com/a-omori-yumemi/YumetterAPI/querier"
	"github.com/a-omori-yumemi/YumetterAPI/repository"
	"github.com/a-omori-yumemi/YumetterAPI/usecase"
	"github.com/labstack/echo/v4"
)

func PostTweet(tweetRepo repository.ITweetRepository) echo.HandlerFunc {

	GetParams := func(c echo.Context) (usrID model.UsrIDType, body string, repliedTo *model.TwIDType, err error) {
		if usrIDPtr := GetSessionUserID(c); usrIDPtr != nil {
			usrID = *usrIDPtr
		} else {
			return usrID, body, repliedTo, echo.NewHTTPError(401)
		}

		body = c.FormValue("body")

		if repliedToStr := c.FormValue("replied_to"); repliedToStr != "" {
			if repliedToTmp, err := strconv.Atoi(repliedToStr); err == nil {
				repliedTo = new(model.TwIDType)
				*repliedTo = model.TwIDType(repliedToTmp)
			} else {
				return usrID, body, repliedTo, echo.NewHTTPError(400, "invalid parameter: {replied_to}")
			}
		}

		return usrID, body, repliedTo, nil
	}

	return func(c echo.Context) error {
		usrID, body, repliedTo, err := GetParams(c)
		if err != nil {
			return err
		}
		tweet := model.Tweet{
			Body:      body,
			UsrID:     usrID,
			RepliedTo: repliedTo,
		}
		if tweet.Validate() != nil {
			return echo.NewHTTPError(400, "body is too short or the tweet you replied to is missing")
		}

		tweet, err = tweetRepo.AddTweet(tweet)
		if err != nil {
			return err
		}
		return c.JSON(200, tweet)
	}
}

func GetTweet(tweetQuerier querier.ITweetQuerier) echo.HandlerFunc {

	return func(c echo.Context) error {
		twID, err := strconv.Atoi(c.Param("tw_id"))
		if err != nil {
			return echo.NewHTTPError(400, err)
		}

		tweet, err := tweetQuerier.FindTweet(model.TwIDType(twID))
		if err == model.ErrNotFound {
			return echo.NewHTTPError(404, "tweet not found")
		} else if err != nil {
			return err
		}

		return c.JSON(200, tweet)
	}
}

func GetTweets(tweetDetailUsecase querier.ITweetDetailQuerier) echo.HandlerFunc {
	const DefaultLimitValue = 30

	GetParams := func(c echo.Context) (*model.UsrIDType, int, *model.TwIDType, error) {
		usrID := GetSessionUserID(c)

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

		tweet, err := tweetDetailUsecase.FindTweetDetails(usrID, limit, repliedTo)
		if err != nil {
			return err
		}

		return c.JSON(200, tweet)
	}
}

func DeleteTweet(tweetDeleteUsecase usecase.ITweetDeleteUsecase) echo.HandlerFunc {

	GetParams := func(c echo.Context) (model.UsrIDType, model.TwIDType, error) {
		var usrID model.UsrIDType
		if usrIDPtr := GetSessionUserID(c); usrIDPtr != nil {
			usrID = *usrIDPtr
		} else {
			return 0, 0, echo.NewHTTPError(401)
		}

		twID, err := strconv.Atoi(c.Param("tw_id"))
		if err != nil {
			return 0, 0, echo.NewHTTPError(400, err)
		}
		return usrID, model.TwIDType(twID), nil
	}

	return func(c echo.Context) error {
		usrID, twID, err := GetParams(c)
		if err != nil {
			return err
		}

		err = tweetDeleteUsecase.DeleteTweetWithAuth(usrID, twID)
		if err == model.ErrNotFound {
			return echo.NewHTTPError(404, "tweet not found")
		} else if err == model.ErrForbidden {
			return echo.NewHTTPError(403, "Only author can delete this tweet")
		} else if err != nil {
			return err
		}

		return c.NoContent(204)
	}
}
