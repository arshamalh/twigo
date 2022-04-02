package twigo

import "fmt"

type Client struct {
	APIKey          string
	APISecret       string
	BearerToken     string
	AccessToken     string
	AccessSecret    string
	waitOnRateLimit bool
}

type Response struct {
	Data     interface{}
	Includes interface{}
	Errors   interface{}
	Meta     interface{}
}

func NewClient(api_key, api_secret, bearer_token, access_token, access_secret string, wait_on_rate_limit bool) *Client {
	return &Client{
		APIKey:          api_key,
		APISecret:       api_secret,
		BearerToken:     bearer_token,
		AccessToken:     access_token,
		AccessSecret:    access_secret,
		waitOnRateLimit: wait_on_rate_limit,
	}
}

func (client *Client) make_request(method, url string, params map[string]string, body interface{}, user_auth bool) (*Response, error) {
	// base_url := "https://api.twitter.com"
	headers := map[string]string{"User-Agent": "twigo"}
	if user_auth {
		headers["Authorization"] = fmt.Sprintf("Basic %s", client.AccessToken)
	} else {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", client.BearerToken)
	}
	return nil, nil
}

func (client *Client) GetMe() *Response
func (client *Client) CreateTweet(text string, options ...interface{}) *Response

// tweet_id and target_user_id can be string or int
func (client *Client) Like(tweet_id string, options ...interface{}) *Response
func (client *Client) Retweet(tweet_id string, options ...interface{}) *Response

func (client *Client) SearchAllTweets(query string, options ...interface{}) *Response
func (client *Client) FollowUser(target_user_id string, options ...interface{}) *Response
func (client *Client) UnfollowUser(target_user_id string, options ...interface{}) *Response
func (client *Client) Unlike(tweet_id string, options ...interface{}) *Response
func (client *Client) Unretweet(tweet_id string, options ...interface{}) *Response
