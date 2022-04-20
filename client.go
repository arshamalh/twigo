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
	userID            string
}

// ** Requests ** //
func (c *Client) request(method, route string, params map[string]interface{}) (*http.Response, error) {
	// OAuth_1a is always true for post and put routes
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

func (c *Client) get_request(route string, oauth_1a bool, params map[string]interface{}, endpoint_parameters []string) (*http.Response, error) {
	// oauth_1a ==> Whether or not to use OAuth 1.0a User context
	parsedRoute, err := url.Parse(route)
	if err != nil {
		return nil, err
	}

	parameters := url.Values{}
	for param_name, param_value := range params {
		if !contains(endpoint_parameters, param_name) {
			fmt.Printf(" it seems endpoint parameter '%s' is not supported", param_name)
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

		default:
			return nil, fmt.Errorf("%s with value of %s is not supported, please contact us", param_name, param_value)
		}

	}
	parsedRoute.RawQuery = parameters.Encode()
	fullRoute := base_route + parsedRoute.String()
	fmt.Println("Route:>> ", fullRoute)
	if c.read_only_access { oauth_1a = false }
	if c.bearerToken == "" { oauth_1a = true }
	if oauth_1a {
		//%% TODO: Should we define authorizedClient here? or tweepy is doing it wrong?
		return c.authorizedClient.Get(fullRoute)
	} else {
		request, err := http.NewRequest("GET", fullRoute, nil)
		if err != nil {
			return nil, err
		}

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.bearerToken))
		client := http.Client{}
		return client.Do(request)
	}
}

func (c *Client) delete_request(route string) (*http.Response, error) {
	// OAuth_1a is always true for delete routes
	request, err := http.NewRequest("DELETE", base_route+route, nil)
	if err != nil {
		return nil, err
	}
	return c.authorizedClient.Do(request)
}

// %%TODO: maybe using map[string]interface{} for params is not a good approach, we can use a predefined struct instead.
// the second approch will help other developers to understand what's going on and which params to pass.
// but first approach is more convenient for those who know which params they should pass.

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

func (c *Client) GetLikingUsers(tweet_id string, oauth_1a bool, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"expansions", "media.fields", "place.fields",
		"poll.fields", "tweet.fields", "user.fields",
		"max_results",
	}

	route := fmt.Sprintf("tweets/%s/liking_users", tweet_id)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UsersResponse{}).Parse(response)
}

func (c *Client) GetLikedTweets(user_id string, oauth_1a bool, params map[string]interface{}) (*TweetsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "media.fields",
		"pagination_token", "place.fields", "poll.fields",
		"tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("users/%s/liked_tweets", user_id)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&TweetsResponse{}).Parse(response)
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
func (c *Client) Retweet(tweet_id string) (*http.Response, error) {
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

func (c *Client) UnRetweet(tweet_id string) (*http.Response, error) {
	route := fmt.Sprintf("users/%s/retweets/%s", c.userID, tweet_id)
	return c.delete_request(route)
}

func (c *Client) GetRetweeters(tweet_id string, oauth_1a bool, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"expansions", "media.fields", "place.fields",
		"poll.fields", "tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("tweets/%s/retweeted_by", tweet_id)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UsersResponse{}).Parse(response)
}

// ** Search tweets ** //
// func (c *Client) SearchRecentTweets() (*http.Response, error)
// func (c *Client) SearchAllTweets(query string, params map[string]interface{}) (*http.Response, error)
// func QueryMaker()

// ** Timelines ** //
func (c *Client) GetUserTweets(user_id string, oauth_1a bool, params map[string]interface{}) (*TweetsResponse, error) {
	endpoint_parameters := []string{
		"end_time", "exclude", "expansions", "max_results",
		"media.fields", "pagination_token", "place.fields",
		"poll.fields", "since_id", "start_time", "tweet.fields",
		"until_id", "user.fields",
	}
	route := fmt.Sprintf("users/%s/tweets", user_id)

	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&TweetsResponse{}).Parse(response)
}

func (c *Client) GetUserMentions(user_id string, oauth_1a bool, params map[string]interface{}) (*TweetsResponse, error) {
	endpoint_parameters := []string{
		"end_time", "expansions", "max_results", "media.fields",
		"pagination_token", "place.fields", "poll.fields", "since_id",
		"start_time", "tweet.fields", "until_id", "user.fields",
	}

	route := fmt.Sprintf("users/%s/mentions", user_id)

	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&TweetsResponse{}).Parse(response)
}

// ** Tweet counts ** //
func (c *Client) GetAllTweetsCount(query string, params map[string]interface{}) (*http.Response, error) {
	endpoint_parameters := []string{
		"end_time", "granularity", "next_token", "query",
		"since_id", "start_time", "until_id",
	}
	if params == nil {
		params = make(map[string]interface{})
	}
	params["query"] = query
	return c.get_request("tweets/counts/all", false, params, endpoint_parameters)
}

func (c *Client) GetRecentTweetsCount(query string, params map[string]interface{}) (*http.Response, error) {
	endpoint_parameters := []string{
		"end_time", "granularity", "query",
		"since_id", "start_time", "until_id",
	}
	if params == nil {
		params = make(map[string]interface{})
	}
	params["query"] = query
	return c.get_request("tweets/counts/recent", false, params, endpoint_parameters)
}

// ** Tweet lookup ** //
func (c *Client) GetTweet(tweet_id string, oauth_1a bool, params map[string]interface{}) (*TweetResponse, error) {
	endpoint_parameters := []string{
		"expansions", "media.fields", "place.fields",
		"poll.fields", "tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("tweets/%s", tweet_id)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&TweetResponse{}).Parse(response)
}

func (c *Client) GetTweets(tweet_ids []string, oauth_1a bool, params map[string]interface{}) (*TweetsResponse, error) {
	endpoint_parameters := []string{
		"ids", "expansions", "media.fields", "place.fields",
		"poll.fields", "tweet.fields", "user.fields",
	}
	if params == nil {
		params = make(map[string]interface{})
	}
	params["ids"] = tweet_ids
	response, err := c.get_request("tweets", oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&TweetsResponse{}).Parse(response)
}

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
	endpoint_parameters := []string{
		"expansions", "max_results", "tweet.fields",
		"user.fields", "pagination_token",
	}
	route := fmt.Sprintf("users/%s/blocking", c.userID)
	return c.get_request(route, true, params, endpoint_parameters)
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

func (c *Client) UnfollowUser(target_user_id string) (*http.Response, error) {
	route := fmt.Sprintf("users/%s/following/%s", c.userID, target_user_id)
	return c.delete_request(route)
}

func (c *Client) GetUserFollowers(user_id string, oauth_1a bool, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "tweet.fields",
		"user.fields", "pagination_token",
	}
	route := fmt.Sprintf("users/%s/followers", user_id)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UsersResponse{}).Parse(response)
}

func (c *Client) GetUserFollowing(user_id string, oauth_1a bool, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "tweet.fields",
		"user.fields", "pagination_token",
	}
	route := fmt.Sprintf("users/%s/following", user_id)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UsersResponse{}).Parse(response)
}

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

func (c *Client) GetMuted(oauth_1a bool, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "tweet.fields",
		"user.fields", "pagination_token",
	}
	route := fmt.Sprintf("users/%s/muting", c.userID)
	response, err := c.get_request(route, true, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UsersResponse{}).Parse(response)
}

// ** User lookup ** //
func (c *Client) GetUserByID(user_id string, oauth_1a bool, params map[string]interface{}) (*UserResponse, error) {
	endpoint_parameters := []string{
		"expansions", "tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("users/%s", user_id)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UserResponse{}).Parse(response)
}

func (c *Client) GetUserByUsername(username string, oauth_1a bool, params map[string]interface{}) (*UserResponse, error) {
	endpoint_parameters := []string{
		"expansions", "tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("users/by/username/%s", username)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UserResponse{}).Parse(response)
}

func (c *Client) GetUsersByIDs(user_ids []string, oauth_1a bool, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"usernames", "ids", "expansions",
		"tweet.fields", "user.fields",
	}

	if user_ids == nil {
		return nil, fmt.Errorf("user_ids are required")
	}
	if params == nil {
		params = make(map[string]interface{})
	}
	params["ids"] = user_ids

	response, err := c.get_request("users", oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UsersResponse{}).Parse(response)
}

func (c *Client) GetUsersByUsernames(usernames []string, oauth_1a bool, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"usernames", "ids", "expansions",
		"tweet.fields", "user.fields",
	}

	if usernames == nil {
		return nil, fmt.Errorf("usernames are required")
	}
	if params == nil {
		params = make(map[string]interface{})
	}
	params["usernames"] = usernames

	response, err := c.get_request("users/by", oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UsersResponse{}).Parse(response)
}

// ** Spaces ** //
// func (c *Client) SearchSpaces(query string) (*http.Response, error)
// func (c *Client) GetSpaces(space_ids, user_ids []string) (*http.Response, error)
// func (c *Client) GetSpace(space_id string) (*http.Response, error)
// func (c *Client) GetSpaceBuyers(space_id string) (*http.Response, error)

// ** List Tweets lookup ** //
func (c *Client) GetListTweets(list_id string, oauth_1a bool, params map[string]interface{}) (*TweetsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "pagination_token",
		"tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("lists/%s/tweets", list_id)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&TweetsResponse{}).Parse(response)
}

// ** List follows ** //
func (c *Client) FollowList(list_id string) (*http.Response, error) {
	data := map[string]interface{}{
		"list_id": list_id,
	}

	route := fmt.Sprintf("users/%s/followed_lists", c.userID)

	return c.request(
		"POST",
		route,
		data,
	)
}

func (c *Client) UnfollowList(list_id string) (*http.Response, error) {
	route := fmt.Sprintf("users/%s/followed_lists/%s", c.userID, list_id)
	return c.delete_request(route)
}

func (c *Client) GetListFollowers(list_id string, oauth_1a bool, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "pagination_token",
		"tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("lists/%s/followers", list_id)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UsersResponse{}).Parse(response)
}

func (c *Client) GetFollowedLists(user_id string, oauth_1a bool, params map[string]interface{}) (*ListsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "pagination_token",
		"list.fields", "user.fields",
	}
	route := fmt.Sprintf("users/%s/followed_lists", user_id)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&ListsResponse{}).Parse(response)
}

// ** List lookup ** //
func (c *Client) GetList(list_id string, oauth_1a bool, params map[string]interface{}) (*ListResponse, error) {
	endpoint_parameters := []string{
		"expansions", "list.fields", "user.fields",
	}
	route := fmt.Sprintf("lists/%s", list_id)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&ListResponse{}).Parse(response)
}

func (c *Client) GetOwnedLists(user_id string, oauth_1a bool, params map[string]interface{}) (*ListsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "pagination_token",
		"list.fields", "user.fields",
	}
	route := fmt.Sprintf("users/%s/owned_lists", user_id)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&ListsResponse{}).Parse(response)
}

// ** List members ** //
func (c *Client) AddListMemeber(list_id, user_id string) (*http.Response, error) {
	data := map[string]interface{}{
		"user_id": user_id,
	}

	route := fmt.Sprintf("lists/%s/members", list_id)

	return c.request(
		"POST",
		route,
		data,
	)
}

func (c *Client) RemoveListMember(list_id, user_id string) (*http.Response, error) {
	route := fmt.Sprintf("lists/%s/members/%s", list_id, user_id)
	return c.delete_request(route)
}

func (c *Client) GetListMembers(list_id string, oauth_1a bool, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "pagination_token",
		"tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("lists/%s/members", list_id)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UsersResponse{}).Parse(response)
}

func (c *Client) GetListMemberships(user_id string, oauth_1a bool, params map[string]interface{}) (*ListsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "pagination_token",
		"list.fields", "user.fields",
	}
	route := fmt.Sprintf("users/%s/list_memberships", user_id)
	response, err := c.get_request(route, oauth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&ListsResponse{}).Parse(response)
}

// ** Manage Lists ** //
func (c *Client) CreateList(name string, description string, private bool, params map[string]interface{}) (*http.Response, error) {
	data := map[string]interface{}{
		"name":        name,
		"description": description,
		"private":     private,
	}

	route := "lists"

	return c.request(
		"POST",
		route,
		data,
	)
}

func (c *Client) UpdateList(list_id string, name string, description string, private bool, params map[string]interface{}) (*http.Response, error) {
	data := map[string]interface{}{
		"name":        name,
		"description": description,
		"private":     private,
	}

	route := fmt.Sprintf("lists/%s", list_id)

	return c.request(
		"PUT",
		route,
		data,
	)
}

func (c *Client) DeleteList(list_id string) (*http.Response, error) {
	route := fmt.Sprintf("lists/%s", list_id)
	return c.delete_request(route)
}

// ** Pinned Lists ** //
func (c *Client) PinList(list_id string) (*http.Response, error) {
	data := map[string]interface{}{
		"list_id": list_id,
	}

	route := fmt.Sprintf("users/%s/pinned_lists", c.userID)

	return c.request(
		"POST",
		route,
		data,
	)
}

func (c *Client) UnpinList(list_id string) (*http.Response, error) {
	route := fmt.Sprintf("users/%s/pinned_lists/%s", c.userID, list_id)
	return c.delete_request(route)
}

func (c *Client) GetPinnedLists(params map[string]interface{}) (*ListsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "list.fields", "user.fields",
	}
	route := fmt.Sprintf("users/%s/pinned_lists", c.userID)
	response, err := c.get_request(route, true, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&ListsResponse{}).Parse(response)
}

// ** Batch Compliance ** //
// func (c *Client) GetComplianceJobs(_type string) (*http.Response, error)
// func (c *Client) GetComplianceJob(id string) (*http.Response, error)
// func (c *Client) CreateComplianceJobs(_type, name, resumable string) (*http.Response, error)

// func tweetResponseParser() {}
// func userResponseParser(){}

// func (c *Client) GetMe() *Response

// I don't know when (where) to use oauth 1, and when (where) to use oauth 2
// Is oauth 2 only using bearer token?
// Is bearer token readonly?
// Why tweepy has separated some get routes and is using bearer for them?
// Why twitter is using bearer for tweeting but once it said that brearer is readonly?
// Some wants only bearer
// Some wants only consumer
// Some (writes) only can use consumer
// Some (gets) is better to use bearer
