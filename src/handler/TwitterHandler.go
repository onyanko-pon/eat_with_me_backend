package handler

import (
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"github.com/labstack/echo/v4"
)

type TwitterHandler struct{}

func NewTwitterHandler() (*TwitterHandler, error) {
	return &TwitterHandler{}, nil
}

func (h TwitterHandler) FetchRequestToken(c echo.Context) error {

	config := oauth1.Config{
		ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		CallbackURL:    os.Getenv("TWITTER_OAUTH_CALLBACK_URL"),
		Endpoint:       twitter.AuthorizeEndpoint,
	}

	requestToken, requestSecret, err := config.RequestToken()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"requestToken":  requestToken,
		"requestSecret": requestSecret,
	})
}
