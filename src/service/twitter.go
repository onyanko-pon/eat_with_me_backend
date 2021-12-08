package service

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"github.com/onyanko-pon/eat_with_me_backend/src/entity"
)

type TwitterAuthService struct{}

/*
oauthトークンからアクセストークンを取得する
*/
func (s TwitterAuthService) GenAccessToken(oauthToken string, oauthSecret string, oauthVerifier string) (string, string, error) {
	config := oauth1.Config{
		ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		CallbackURL:    os.Getenv("TWITTER_OAUTH_CALLBACK_URL"),
		Endpoint:       twitter.AuthorizeEndpoint,
	}

	accessToken, accessSecret, err := config.AccessToken(oauthToken, oauthSecret, oauthVerifier)

	return accessToken, accessSecret, err
}

/*
アクセストークンからTwitterユーザーを取得する
*/
func (s TwitterAuthService) VerifyUser(accessToken string, accessSecret string) (*entity.TwitterUser, error) {
	// TODO 然るべき configにする
	config := oauth1.Config{
		ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		CallbackURL:    os.Getenv("TWITTER_OAUTH_CALLBACK_URL"),
		Endpoint:       twitter.AuthorizeEndpoint,
	}
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	path := "https://api.twitter.com/1.1/account/verify_credentials.json"
	resp, _ := httpClient.Get(path)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var twitterUser entity.TwitterUser
	err := json.Unmarshal(body, &twitterUser)
	if err != nil {
		return nil, err
	}

	return &twitterUser, nil
}
