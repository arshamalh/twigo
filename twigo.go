package twigo

import (
	"fmt"
	"strings"

	"github.com/arshamalh/twigo/utils"
	"github.com/mrjones/oauth"
)

type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	BearerToken    string
}

func NewClient(config *Config) (*Client, error) {
	keys_exists := config.ConsumerKey != "" && config.ConsumerSecret != "" && config.AccessToken != "" && config.AccessSecret == ""

	if !keys_exists {
		if config.BearerToken == "" {
			return nil, fmt.Errorf("consumer key, consumer secret, access token and access token secret must be provided")
		}

		userID := ""
		if config.AccessToken != "" {
			userID = strings.Split(config.AccessToken, "-")[0]
		}

		return &Client{
			nil,
			config.ConsumerKey,
			config.ConsumerSecret,
			config.AccessToken,
			config.AccessSecret,
			config.BearerToken,
			true,
			userID,
			OAuth_2,
		}, nil
	}

	if config.BearerToken == "" {
		if bearer_token, err := utils.BearerFinder(config.ConsumerKey, config.ConsumerSecret); err == nil {
			config.BearerToken = bearer_token
		} else {
			return nil, err
		}
	}

	userID := strings.Split(config.AccessToken, "-")[0]

	// TODO: I'm authenticating here, but Do I need to authenticate every once in a while?
	consumer := oauth.NewConsumer(
		config.ConsumerKey,
		config.ConsumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})

	t := oauth.AccessToken{
		Token:  config.AccessToken,
		Secret: config.AccessSecret,
	}

	authorizedClient, err := consumer.MakeHttpClient(&t)
	return &Client{
		authorizedClient,
		config.ConsumerKey,
		config.ConsumerSecret,
		config.AccessToken,
		config.AccessSecret,
		config.BearerToken,
		false,
		userID,
		OAuth_Default,
	}, err
}

func NewBearerOnlyClient(bearerToken string) (*Client, error) {
	return &Client{
		nil,
		"",
		"",
		"",
		"",
		bearerToken,
		true,
		"",
		OAuth_2,
	}, nil
}
