package twigo

import (
	"fmt"
	"strings"

	"github.com/mrjones/oauth"
)

func NewClient(consumerKey, consumerSecret, accessToken, accessTokenSecret, bearerToken string) (*Client, error) {
	keys_exists := consumerKey != "" && consumerSecret != "" && accessToken != "" && accessTokenSecret != ""
	read_only_access := bearerToken != "" && !keys_exists

	if !read_only_access && !keys_exists {
		return nil, fmt.Errorf("consumer key, consumer secret, access token and access token secret must be provided")
	}

	if !read_only_access {
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
		consumer.Debug(false)

		t := oauth.AccessToken{
			Token:  accessToken,
			Secret: accessTokenSecret,
		}

		authorizedClient, err := consumer.MakeHttpClient(&t)
		return &Client{authorizedClient, consumerKey, consumerSecret, accessToken, accessTokenSecret, bearerToken, read_only_access, userID}, err
	}

	return &Client{nil, consumerKey, consumerSecret, accessToken, accessTokenSecret, bearerToken, read_only_access, ""}, nil
}
