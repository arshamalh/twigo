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
	dataPayload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, base_url+url, bytes.NewBuffer(dataPayload))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := c.authorizedClient.Do(request)

	return response, err
}

// ** Manage Tweets ** //
func (c *Client) CreateTweet(text string, options ...interface{}) (*http.Response, error) {
	data := map[string]interface{}{
		"text": text,
	}
	return c.request(
		"POST",
		"tweets",
		data,
	)
}

func (c *Client) DeleteTweet(tweet_id string, options ...interface{}) (*http.Response, error) {
	url := fmt.Sprintf("tweets/%s", tweet_id)

	return c.request(
		"DELETE",
		url,
		nil,
	)
}

// ** Likes ** //
func (c *Client) Like(tweet_id string, options ...interface{}) (*http.Response, error) {
	data := map[string]interface{}{
		"tweet_id": tweet_id,
	}
	url := fmt.Sprintf("users/%s/likes", c.userID)
	return c.request(
		"POST",
		url,
		data,
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

// ** Hide replies ** //
func (c *Client) HideReply(reply_id string) (*http.Response, error) {
	data := map[string]interface{}{
		"hidden": true,
	}
	url := fmt.Sprintf("tweets/%s/hidden", reply_id)
	
	return c.request(
		"PUT",
		url,
		data,
	)
}

func (c *Client) UnHideReply(reply_id string) (*http.Response, error) {
	data := map[string]interface{}{
		"hidden": false,
	}
	url := fmt.Sprintf("tweets/%s/hidden", reply_id)
	
	return c.request(
		"PUT",
		url,
		data,
	)
}

// ** Retweets ** //
func (c *Client) Retweet(tweet_id string, options ...interface{}) (*http.Response, error) {
	data := map[string]interface{}{
		"tweet_id": tweet_id,
	}
	url := fmt.Sprintf("users/%s/retweets", c.userID)
	return c.request(
		"POST",
		url,
		data,
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
// func (c *Client) GetRetweeters() (*http.Response, error)

// ** Search tweets ** //
// func (c *Client) SearchRecentTweets() (*http.Response, error)
// func (c *Client) SearchAllTweets(query string, options ...interface{}) (*http.Response, error)
// func QueryMaker()

// ** Timelines ** //
func (c *Client) GetUserTweets(user_id string) (*http.Response, error) {
	url := fmt.Sprintf("users/%s/tweets", user_id)

	return c.request(
		"GET",
		url,
		nil,
	)
}

func (c *Client) GetUserMentions(user_id string) (*http.Response, error) {
	url := fmt.Sprintf("users/%s/mentions", user_id)

	return c.request(
		"GET",
		url,
		nil,
	)
}

// ** Tweet counts ** //
// func (c *Client) GetAllTweetsCount() (*http.Response, error)
// func (c *Client) GetRecentTweetsCount() (*http.Response, error)

// ** Tweet lookup ** //
// func (c *Client) GetTweet() (*http.Response, error)
// func (c *Client) GetTweets() (*http.Response, error)

// ** Blocks ** //
func (c *Client) Block(target_user_id string) (*http.Response, error) {
	data := map[string]interface{}{
		"target_user_id": target_user_id,
	}
	url := fmt.Sprintf("users/%s/blocking", c.userID)

	return c.request(
		"POST",
		url,
		data,
	)
}

func (c *Client) UnBlock(target_user_id string) (*http.Response, error) {
	url := fmt.Sprintf("users/%s/blocking/%s", c.userID, target_user_id)
	return c.request(
		"DELETE",
		url,
		nil,
	)
}
// func (c *Client) GetBlocked() (*http.Response, error)

// ** Follows ** //
func (c *Client) FollowUser(target_user_id string, options ...interface{}) (*http.Response, error) {
	data := map[string]interface{}{
		"target_user_id": target_user_id,
	}

	url := fmt.Sprintf("users/%s/following", c.userID)

	return c.request(
		"POST",
		url,
		data,
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

// func (c *Client) GetUserFollowers(user_id string) (*http.Response, error)
// func (c *Client) GetUserFollowing(user_id string) (*http.Response, error)

// ** Mutes ** //
func (c *Client) Mute(target_user_id string) (*http.Response, error) {
	data := map[string]interface{}{
		"target_user_id": target_user_id,
	}
	
	url := fmt.Sprintf("users/%s/muting", c.userID)

	return c.request(
		"POST", 
		url, 
		data,
	)
}

func (c *Client) UnMute(target_user_id string) (*http.Response, error) {
	url := fmt.Sprintf("users/%s/muting/%s", c.userID, target_user_id)

	return c.request(
		"DELETE", 
		url, 
		nil,
	)
}

func (c *Client) GetMuted() (*http.Response, error) {
	url := fmt.Sprintf("users/%s/muting", c.userID)
	return c.request(
		"GET",
		url,
		nil,
	)
}

// ** User lookup ** //
// func (c *Client) GetUser(user_id, username string) (*http.Response, error)
// func (c *Client) GetUsers(user_ids, usernames []string) (*http.Response, error)

// ** Spaces ** //
// func (c *Client) SearchSpaces(query string) (*http.Response, error)
// func (c *Client) GetSpaces(space_ids, user_ids []string) (*http.Response, error)
// func (c *Client) GetSpace(space_id string) (*http.Response, error)
// func (c *Client) GetSpaceBuyers(space_id string) (*http.Response, error)

// ** Batch Compliance ** //
// func (c *Client) GetComplianceJobs(_type string) (*http.Response, error)
// func (c *Client) GetComplianceJob(id string) (*http.Response, error)
// func (c *Client) CreateComplianceJobs(_type, name, resumable string) (*http.Response, error)

// func tweetResponseParser() {}
// func userResponseParser(){}

// func (c *Client) GetMe() *Response

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
