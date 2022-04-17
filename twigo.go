package twigo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mrjones/oauth"
)

const (
	base_url = "https://api.twitter.com/2/"
)

type Client struct {
	authorizedClient  *http.Client
	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string
	bearerToken       string
	read_only_access  bool
	bearer_token_auth bool
}

type Response struct {
	Data     interface{}
	Includes interface{}
	Errors   interface{}
	Meta     interface{}
}

func NewClient(consumerKey, consumerSecret, accessToken, accessTokenSecret, bearerToken string, wait_on_rate_limit bool) (*Client, error) {
	keys_exists := consumerKey != "" && consumerSecret != "" && accessToken != "" && accessTokenSecret != ""
	bearer_token_auth := bearerToken != ""
	read_only_access := bearer_token_auth && !keys_exists

	if !read_only_access && !keys_exists {
		return nil, fmt.Errorf("consumer key, consumer secret, access token and access token secret must be provided")
	}

	if !read_only_access {
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
		return &Client{authorizedClient, consumerKey, consumerSecret, accessToken, accessTokenSecret, bearerToken, read_only_access, bearer_token_auth}, err
	}

	return &Client{nil, consumerKey, consumerSecret, accessToken, accessTokenSecret, bearerToken, read_only_access, bearer_token_auth}, nil
}

func (c *Client) request(method, url string, params map[string]interface{}) (*http.Response, error) {
	if method == "GET" {
		return c.authorizedClient.Get(base_url+url)
	}
	DataPayload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, base_url+url, bytes.NewBuffer(DataPayload))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := c.authorizedClient.Do(request)

	return response, err
}

func (c *Client) GetUserTweets(user_id string) (*http.Response, error) {
	return c.request(
		"GET", 
		fmt.Sprintf("users/%s/tweets", user_id), 
		nil,
	)
}

// tweet_id and target_user_id can be string or int
func (c *Client) CreateTweet(text string, options ...interface{}) (*http.Response, error) {
	Data := map[string]interface{}{
		"text": text,
	}
	return c.request(
		"POST",
		"tweets",
		Data,
	)
}

func (c *Client) Like(tweet_id string, options ...interface{}) (*http.Response, error) {
	Data := map[string]interface{}{
		"tweet_id": tweet_id,
	}
	user_id := strings.Split(c.accessToken, "-")[0]	
	url := fmt.Sprintf("users/%s/likes", user_id)
	return c.request(
		"POST",
		url,
		Data,
	)
}

// func (c *Client) GetMe() *Response
// func (c *Client) Retweet(tweet_id string, options ...interface{}) *Response
// func (c *Client) SearchAllTweets(query string, options ...interface{}) *Response
// func (c *Client) FollowUser(target_user_id string, options ...interface{}) *Response
// func (c *Client) UnfollowUser(target_user_id string, options ...interface{}) *Response
// func (c *Client) Unlike(tweet_id string, options ...interface{}) *Response
// func (c *Client) Unretweet(tweet_id string, options ...interface{}) *Response

// Some wants only bearer
// Some wants only consumer
// Some (writes) only can use consumer
// Some (gets) is better to use bearer
