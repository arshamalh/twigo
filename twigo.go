package twigo

import (
	"fmt"
	"strings"

	"github.com/arshamalh/twigo/utils"
	"github.com/mrjones/oauth"
)

func NewClient(consumerKey, consumerSecret, accessToken, accessTokenSecret, bearerToken string) (*Client, error) {
	keys_exists := consumerKey != "" && consumerSecret != "" && accessToken != "" && accessTokenSecret != ""

	if !keys_exists {
		if bearerToken == "" {
			return nil, fmt.Errorf("consumer key, consumer secret, access token and access token secret must be provided")
		}

		return &Client{
			nil,
			consumerKey,
			consumerSecret,
			accessToken,
			accessTokenSecret,
			bearerToken,
			true,
			"",
			OAuth_Default,
		}, nil
	}

	if bearer_token, err := utils.BearerFinder(consumerKey, consumerSecret); bearerToken == "" && err == nil {
		bearerToken = bearer_token
	}

	userID := strings.Split(accessToken, "-")[0]

	// TODO: I'm authenticating here, but Do I need to authenticate every once in a while?
	consumer := oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})

	t := oauth.AccessToken{
		Token:  accessToken,
		Secret: accessTokenSecret,
	}

	authorizedClient, err := consumer.MakeHttpClient(&t)
	return &Client{
		authorizedClient,
		consumerKey,
		consumerSecret,
		accessToken,
		accessTokenSecret,
		bearerToken,
		false,
		userID,
		OAuth_Default,
	}, err
}
