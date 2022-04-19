package twigo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	base_route = "https://api.twitter.com/2/"
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

// ** Requests ** //
func (c *Client) request(method, route string, params map[string]interface{}) (*http.Response, error) {
	dataPayload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, base_route+route, bytes.NewBuffer(dataPayload))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := c.authorizedClient.Do(request)

	return response, err
}

func (c *Client) get_request(route string, params map[string]interface{}, endpoint_parameters []string) (*http.Response, error) {
	parsedRoute, err := url.Parse(route)
	if err != nil {
		return nil, err
	}

	parameters := url.Values{}
	for param_name, param_value := range params {
		if !contains(endpoint_parameters, param_name) {
			return nil, fmt.Errorf("endpoint parameter '%s' is not supported", param_name)
		}
		switch param_valt := param_value.(type) {
		case int:
			parameters.Add(param_name, strconv.Itoa(param_valt))
		case string:
			parameters.Add(param_name, param_valt)
		case []string:
			parameters.Add(param_name, strings.Join(param_valt, ","))
		// TODO: case for arrays of anything else not only string
		// TODO: case datetime
		// 	if param_value.tzinfo is not None:
		// 		param_value = param_value.astimezone(datetime.timezone.utc)
		// 	request_params[param_name] = param_value.strftime("%Y-%m-%dT%H:%M:%S.%fZ")
		// if param_name is not in the endpoint_parameters, we should warn user, but anyway it's not that much important!
		// Actually it doesn't cause an error, but maybe user had a typo! so I think it's better to return an error.

		default:
			fmt.Println(param_name, param_value)
		}

	}
	parsedRoute.RawQuery = parameters.Encode()
	return c.authorizedClient.Get(base_route + route)
}

func (c *Client) delete_request(route string) (*http.Response, error) {
	request, err := http.NewRequest("DELETE", base_route+route, nil)
	if err != nil {
		return nil, err
	}
	return c.authorizedClient.Do(request)
}

// ** Manage Tweets ** //
func (c *Client) CreateTweet(text string, params map[string]interface{}) (*http.Response, error) {
	data := map[string]interface{}{
		"text": text,
	}
	return c.request(
		"POST",
		"tweets",
		data,
	)
}

func (c *Client) DeleteTweet(tweet_id string) (*http.Response, error) {
	route := fmt.Sprintf("tweets/%s", tweet_id)

	return c.delete_request(route)
}

// ** Likes ** //
func (c *Client) Like(tweet_id string) (*http.Response, error) {
	data := map[string]interface{}{
		"tweet_id": tweet_id,
	}
	route := fmt.Sprintf("users/%s/likes", c.userID)
	return c.request(
		"POST",
		route,
		data,
	)
}

func (c *Client) Unlike(tweet_id string) (*http.Response, error) {
	route := fmt.Sprintf("users/%s/likes/%s", c.userID, tweet_id)
	return c.delete_request(route)
}

func (c *Client) GetLikingUsers(tweet_id string, params map[string]interface{}) (*http.Response, error) {
	route := fmt.Sprintf("tweets/%s/liking_users", tweet_id)
	return c.get_request(route, params, nil)
}

func (c *Client) GetLikedTweets(user_id string, params map[string]interface{}) (*http.Response, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "media.fields",
		"pagination_token", "place.fields", "poll.fields",
		"tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("users/%s/liked_tweets", user_id)
	return c.get_request(route, params, endpoint_parameters)
}

// ** Hide replies ** //
func (c *Client) HideReply(reply_id string) (*http.Response, error) {
	data := map[string]interface{}{
		"hidden": true,
	}
	route := fmt.Sprintf("tweets/%s/hidden", reply_id)

	return c.request(
		"PUT",
		route,
		data,
	)
}

func (c *Client) UnHideReply(reply_id string) (*http.Response, error) {
	data := map[string]interface{}{
		"hidden": false,
	}
	route := fmt.Sprintf("tweets/%s/hidden", reply_id)

	return c.request(
		"PUT",
		route,
		data,
	)
}

// ** Retweets ** //
func (c *Client) Retweet(tweet_id string, params map[string]interface{}) (*http.Response, error) {
	data := map[string]interface{}{
		"tweet_id": tweet_id,
	}
	route := fmt.Sprintf("users/%s/retweets", c.userID)
	return c.request(
		"POST",
		route,
		data,
	)
}

func (c *Client) Unretweet(tweet_id string, params map[string]interface{}) (*http.Response, error) {
	route := fmt.Sprintf("users/%s/retweets/%s", c.userID, tweet_id)
	return c.delete_request(route)
}

// func (c *Client) GetRetweeters() (*http.Response, error)

// ** Search tweets ** //
// func (c *Client) SearchRecentTweets() (*http.Response, error)
// func (c *Client) SearchAllTweets(query string, params map[string]interface{}) (*http.Response, error)
// func QueryMaker()

// ** Timelines ** //
func (c *Client) GetUserTweets(user_id string) (*http.Response, error) {
	route := fmt.Sprintf("users/%s/tweets", user_id)

	return c.get_request(route, nil, nil)
}

func (c *Client) GetUserMentions(user_id string) (*http.Response, error) {
	route := fmt.Sprintf("users/%s/mentions", user_id)

	return c.get_request(route, nil, nil)
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
	route := fmt.Sprintf("users/%s/blocking", c.userID)

	return c.request(
		"POST",
		route,
		data,
	)
}

func (c *Client) UnBlock(target_user_id string) (*http.Response, error) {
	route := fmt.Sprintf("users/%s/blocking/%s", c.userID, target_user_id)
	return c.delete_request(route)
}

func (c *Client) GetBlocked(params map[string]interface{}) (*http.Response, error) {
	// Parameters
	// ----------
	// expansions : Union[List[str], str]
	// 	:ref:`expansions_parameter`
	// max_results : int
	// 	The maximum number of results to be returned per page. This can be
	// 	a number between 1 and 1000. By default, each page will return 100
	// 	results.
	// pagination_token : str
	// 	Used to request the next page of results if all results weren't
	// 	returned with the latest request, or to go back to the previous
	// 	page of results.
	// tweet_fields : Union[List[str], str]
	// 	:ref:`tweet_fields_parameter`
	// user_fields : Union[List[str], str]
	// 	:ref:`user_fields_parameter`
	route := fmt.Sprintf("users/%s/blocking", c.userID)

	return c.get_request(route, params, nil)
}

// ** Follows ** //
func (c *Client) FollowUser(target_user_id string, params map[string]interface{}) (*http.Response, error) {
	data := map[string]interface{}{
		"target_user_id": target_user_id,
	}

	route := fmt.Sprintf("users/%s/following", c.userID)

	return c.request(
		"POST",
		route,
		data,
	)
}

func (c *Client) UnfollowUser(target_user_id string, params map[string]interface{}) (*http.Response, error) {
	route := fmt.Sprintf("users/%s/following/%s", c.userID, target_user_id)
	return c.delete_request(route)
}

// func (c *Client) GetUserFollowers(user_id string) (*http.Response, error)
// func (c *Client) GetUserFollowing(user_id string) (*http.Response, error)

// ** Mutes ** //
func (c *Client) Mute(target_user_id string) (*http.Response, error) {
	data := map[string]interface{}{
		"target_user_id": target_user_id,
	}

	route := fmt.Sprintf("users/%s/muting", c.userID)

	return c.request(
		"POST",
		route,
		data,
	)
}

func (c *Client) UnMute(target_user_id string) (*http.Response, error) {
	route := fmt.Sprintf("users/%s/muting/%s", c.userID, target_user_id)
	return c.delete_request(route)
}

func (c *Client) GetMuted() (*http.Response, error) {
	route := fmt.Sprintf("users/%s/muting", c.userID)
	return c.get_request(route, nil, nil)
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
