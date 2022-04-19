package twigo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	userID            string
}

type Response struct {
	Data     interface{}
	Includes interface{}
	Errors   interface{}
	Meta     interface{}
}

func (c *Client) request(method, url string, params map[string]interface{}) (*http.Response, error) {
	if method == "GET" {
		return c.authorizedClient.Get(base_url + url)
	} else if method == "DELETE" {
		request, err := http.NewRequest(method, base_url+url, nil)
		if err != nil {
			return nil, err
		}
		return c.authorizedClient.Do(request)
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
	url := fmt.Sprintf("users/%s/likes", c.userID)
	return c.request(
		"POST",
		url,
		Data,
	)
}

func (c *Client) Unlike(tweet_id string, options ...interface{}) (*http.Response, error) {
	url := fmt.Sprintf("users/%s/likes/%s", c.userID, tweet_id)
	return c.request(
		"DELETE",
		url,
		nil,
	)
}

func (c *Client) Retweet(tweet_id string, options ...interface{}) (*http.Response, error) {
	Data := map[string]interface{}{
		"tweet_id": tweet_id,
	}
	url := fmt.Sprintf("users/%s/retweets", c.userID)
	return c.request(
		"POST",
		url,
		Data,
	)
}

func (c *Client) Unretweet(tweet_id string, options ...interface{}) (*http.Response, error) {
	url := fmt.Sprintf("users/%s/retweets/%s", c.userID, tweet_id)
	return c.request(
		"DELETE",
		url,
		nil,
	)
}

func (c *Client) FollowUser(target_user_id string, options ...interface{}) (*http.Response, error) {
	Data := map[string]interface{}{
		"target_user_id": target_user_id,
	}

	url := fmt.Sprintf("users/%s/following", c.userID)

	return c.request(
		"POST",
		url,
		Data,
	)
}
func (c *Client) UnfollowUser(target_user_id string, options ...interface{}) (*http.Response, error) {
	url := fmt.Sprintf("users/%s/following/%s", c.userID, target_user_id)

	return c.request(
		"DELETE",
		url,
		nil,
	)
}

func (c *Client) GetLikingUsers(tweet_id string, options ...interface{}) (*http.Response, error) {
	url := fmt.Sprintf("tweets/%s/liking_users", tweet_id)

	return c.request(
		"GET",
		url,
		nil,
	)
}

func (c *Client) GetLikedTweets(user_id string, options ...interface{}) (*http.Response, error) {
	// endpoint_parameters=(
	// 	"expansions", "max_results", "media.fields",
	// 	"pagination_token", "place.fields", "poll.fields",
	// 	"tweet.fields", "user.fields")
	url := fmt.Sprintf("users/%s/liked_tweets", user_id)
	return c.request(
		"GET",
		url,
		nil,
	)
}

// func tweetResponseParser() {}
// func userResponseParser(){}

// func (c *Client) GetMe() *Response
// func (c *Client) SearchAllTweets(query string, options ...interface{}) *Response

// TODO: separate
// I don't know when (where) to use oauth 1, and when (where) to use oauth 2
// Is oauth 2 only using bearer token?
// Is bearer token readonly?
// Why tweepy has separated some get routes and is using bearer for them?
// Why twitter is using bearer for tweeting but once it said that brearer is readonly?
// Some wants only bearer
// Some wants only consumer
// Some (writes) only can use consumer
// Some (gets) is better to use bearer
