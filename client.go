package twigo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arshamalh/twigo/utils"
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
	oauth_type        OAuthType
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

// Sends a get request with specified params
//
// oauth_1a ==> Whether or not to use OAuth 1.0a User context
func (c *Client) get_request(route string, oauth_type OAuthType, params map[string]interface{}, endpoint_parameters []string) (*http.Response, error) {
	parsedRoute, err := url.Parse(route)
	if err != nil {
		return nil, err
	}

	parsedRoute.RawQuery = utils.QueryMaker(params, endpoint_parameters)
	fullRoute := base_route + parsedRoute.String()

	if oauth_type == OAuth_1a {
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

func (c *Client) SetOAuth(oauth_type OAuthType) *Client {
	if c.read_only_access {
		oauth_type = OAuth_2
	}
	if c.bearerToken == "" {
		oauth_type = OAuth_1a
	}
	c.oauth_type = oauth_type
	return c
}

// TODO: Integrate this to get_request, for each request,
// see if the oauth method is default, set it!, but we need the caller function?
// or maybe we can call this function inside each method
func (c *Client) SetDefaultOAuth(caller string) *Client {
	if doa, ok := DefaultOAuthes[caller]; ok {
		c.oauth_type = doa
	} else {
		c.oauth_type = OAuth_2
	}
	return c
}

// ** Manage Tweets ** //

// Creates a Tweet on behalf of an authenticated user
//
// Parameters
//
// text: Text of the Tweet being created. this field is required if media.media_ids is not present, otherwise pass empty string.
//
// params: A map of parameters.
// you can pass some extra parameters, such as:
// 	"direct_message_deep_link", "for_super_followers_only", "media", "geo", "poll", "reply", "reply_settings", "quote_tweet_id",
// Some of these parameters are a little special and should be passed like this:
// 	media := map[string][]string{
// 		"media_ids": []string{}
// 		"tagged_user_ids": []string{}
// 	}
// 	poll := map[string]interface{}{
// 		"options": map[string]string{},
// 		"duration_minutes": int value,
// 	}
// 	reply := map[string]interface{}{
// 		"in_reply_to_tweet_id": "",
// 		"exclude_reply_user_ids": []string{},
// 	}
// 	geo := map[string]string{"place_id": value}
// 	params := map[string]interface{}{
// 		"text": text
// 		"media": media,
// 		"geo": geo,
// 		"poll": poll,
// 		"reply": reply,
// 	}
//
// Reference
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/manage-tweets/api-reference/post-tweets
func (c *Client) CreateTweet(text string, params map[string]interface{}) (*TweetResponse, error) {
	if params == nil {
		params = make(map[string]interface{})
	}

	if text != "" {
		params["text"] = text
	} else if params["media"] == nil {
		return nil, fmt.Errorf("text or media is required")
	}

	response, err := c.request(
		"POST",
		"tweets",
		params,
	)

	if err != nil {
		return nil, err
	}

	return (&TweetResponse{}).Parse(response)
}

// Allows an authenticated user ID to delete a Tweet
//
// Parameters
//
// tweet_id: The Tweet ID you are deleting.
//
// References
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/manage-tweets/api-reference/delete-tweets-id
func (c *Client) DeleteTweet(tweet_id string) (*DeleteResponse, error) {
	route := fmt.Sprintf("tweets/%s", tweet_id)

	response, err := c.delete_request(route)

	if err != nil {
		return nil, err
	}

	return (&DeleteResponse{}).Parse(response)
}

// ** Likes ** //

// Like a Tweet.
//
// Parameters
//
// tweet_id: The ID of the Tweet that you would like to Like.
//
// References
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/likes/api-reference/post-users-id-likesx
func (c *Client) Like(tweet_id string) (*LikeResponse, error) {
	data := map[string]interface{}{
		"tweet_id": tweet_id,
	}

	route := fmt.Sprintf("users/%s/likes", c.userID)

	response, err := c.request(
		"POST",
		route,
		data,
	)

	if err != nil {
		return nil, err
	}

	return (&LikeResponse{}).Parse(response)
}

// Unlike a Tweet.
//
// The request succeeds with no action when the user sends a request to a
// user they're not liking the Tweet or have already unliked the Tweet.
//
// Parameters
//
// tweet_id: The ID of the Tweet that you would like to unlike.
//
// References
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/likes/api-reference/delete-users-id-likes-tweet_id
func (c *Client) Unlike(tweet_id string) (*LikeResponse, error) {
	route := fmt.Sprintf("users/%s/likes/%s", c.userID, tweet_id)

	response, err := c.delete_request(route)

	if err != nil {
		return nil, err
	}

	return (&LikeResponse{}).Parse(response)
}

// Allows you to get information about a Tweet’s liking users.
//
// Parameters
//
// tweet_id: Tweet ID of the Tweet to request liking users of.
//
// params (keys):
// 	"expansions", "media.fields", "place.fields",
// 	"poll.fields", "tweet.fields", "user.fields",
// 	"max_results"
//
// References
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/likes/api-reference/get-tweets-id-liking_users
func (c *Client) GetLikingUsers(tweet_id string, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"expansions", "media.fields", "place.fields",
		"poll.fields", "tweet.fields", "user.fields",
		"max_results", "pagination_token",
	}

	route := fmt.Sprintf("tweets/%s/liking_users", tweet_id)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: tweet_id, Params: params}
	users := &UsersResponse{Caller: c.GetLikingUsers, CallerData: caller_data}

	return users.Parse(response)
}

// Allows you to get information about a user’s liked Tweets.
//
// The Tweets returned by this endpoint count towards the Project-level `Tweet cap`.
//
// Parameters
//
// tweet_id: User ID of the user to request liked Tweets for.
//
// params (keys):
// 	"expansions", "media.fields", "place.fields",
// 	"poll.fields", "tweet.fields", "user.fields",
// 	"max_results"
//
// References
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/likes/api-reference/get-users-id-liked_tweets
func (c *Client) GetLikedTweets(user_id string, params map[string]interface{}) (*TweetsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "media.fields",
		"pagination_token", "place.fields", "poll.fields",
		"tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("users/%s/liked_tweets", user_id)
	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: user_id, Params: params}
	tweets := &TweetsResponse{Caller: c.GetLikedTweets, CallerData: caller_data}

	return tweets.Parse(response)
}

// ** Hide replies ** //

// Hides a reply to a Tweet
//
// Parameters
//
// reply_id:
// 	Unique identifier of the Tweet to hide. The Tweet must belong to a
// 	conversation initiated by the authenticating user.
//
// References
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/hide-replies/api-reference/put-tweets-id-hidden
func (c *Client) HideReply(reply_id string) (*HideReplyResponse, error) {
	data := map[string]interface{}{
		"hidden": true,
	}

	route := fmt.Sprintf("tweets/%s/hidden", reply_id)

	response, err := c.request(
		"PUT",
		route,
		data,
	)

	if err != nil {
		return nil, err
	}

	return (&HideReplyResponse{}).Parse(response)
}

// Unhides a reply to a Tweet
//
// Parameters
//
// reply_id:
// 	Unique identifier of the Tweet to unhide. The Tweet must belong to
//  a conversation initiated by the authenticating user.
//
// References
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/hide-replies/api-reference/put-tweets-id-hidden
func (c *Client) UnHideReply(reply_id string) (*HideReplyResponse, error) {
	data := map[string]interface{}{
		"hidden": false,
	}

	route := fmt.Sprintf("tweets/%s/hidden", reply_id)

	response, err := c.request(
		"PUT",
		route,
		data,
	)

	if err != nil {
		return nil, err
	}

	return (&HideReplyResponse{}).Parse(response)
}

// ** Retweets ** //

// Causes the user ID to Retweet the target Tweet.
//
// Parameters
//
// tweet_id: The ID of the Tweet that you would like to Retweet.
//
// References
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/retweets/api-reference/post-users-id-retweets
func (c *Client) Retweet(tweet_id string) (*RetweetResponse, error) {
	data := map[string]interface{}{
		"tweet_id": tweet_id,
	}

	route := fmt.Sprintf("users/%s/retweets", c.userID)

	response, err := c.request(
		"POST",
		route,
		data,
	)

	if err != nil {
		return nil, err
	}

	return (&RetweetResponse{}).Parse(response)
}

// Allows an authenticated user ID to remove the Retweet of a Tweet.
//
// The request succeeds with no action when the user sends a request to a
// user they're not Retweeting the Tweet or have already removed the
// Retweet of.
//
// Parameters
//
// tweet_id: The ID of the Tweet that you would like to remove the Retweet of.
//
// References
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/retweets/api-reference/delete-users-id-retweets-tweet_id
func (c *Client) UnRetweet(tweet_id string) (*RetweetResponse, error) {
	route := fmt.Sprintf("users/%s/retweets/%s", c.userID, tweet_id)

	response, err := c.delete_request(route)

	if err != nil {
		return nil, err
	}

	return (&RetweetResponse{}).Parse(response)
}

// Allows you to get information about who has Retweeted a Tweet.
//
// Parameters
//
// tweet_id: Tweet ID of the Tweet to request Retweeting users of.
//
// params (keys):
//	"expansions", "media.fields", "place.fields",
// 	"poll.fields", "tweet.fields", "user.fields",
// 	"max_results", "pagination_token"
//
// References
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/retweets/api-reference/get-tweets-id-retweeted_by
func (c *Client) GetRetweeters(tweet_id string, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"expansions", "media.fields", "place.fields",
		"poll.fields", "tweet.fields", "user.fields",
		"max_results", "pagination_token",
	}

	route := fmt.Sprintf("tweets/%s/retweeted_by", tweet_id)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: tweet_id, Params: params}
	users := &UsersResponse{Caller: c.GetRetweeters, CallerData: caller_data}

	return users.Parse(response)
}

// Returns Quote Tweets for a Tweet specified by the requested Tweet ID.
//
// The Tweets returned by this endpoint count towards the Project-level
// `Tweet cap`.
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/quote-tweets/api-reference/get-tweets-id-quote_tweets
func (c *Client) GetQuoteTweets(tweet_id string, params map[string]interface{}) (*TweetsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "media.fields", "place.fields",
		"poll.fields", "tweet.fields", "user.fields",
		"max_results", "pagination_token",
	}

	route := fmt.Sprintf("tweets/%s/quoted_tweets", tweet_id)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: tweet_id, Params: params}
	tweets := &TweetsResponse{Caller: c.GetQuoteTweets, CallerData: caller_data}

	return tweets.Parse(response)
}

// ** Search tweets ** //

// The full-archive search endpoint returns the complete history of public
// Tweets matching a search query; since the first Tweet was created March
// 26, 2006.
//
// This endpoint is only available to those users who have been approved for the `Academic Research product track`
//
// The Tweets returned by this endpoint count towards the Project-level `Tweet cap`.
//
// Parameters
//
// query : str
// 	One query for matching Tweets. Up to 1024 characters.
// end_time : Union[datetime.datetime, str]
// 	YYYY-MM-DDTHH:mm:ssZ (ISO 8601/RFC 3339). Used with ``start_time``.
// 	The newest, most recent UTC timestamp to which the Tweets will be
// 	provided. Timestamp is in second granularity and is exclusive (for
// 	example, 12:00:01 excludes the first second of the minute). If used
// 	without ``start_time``, Tweets from 30 days before ``end_time``
// 	will be returned by default. If not specified, ``end_time`` will
// 	default to [now - 30 seconds].
// max_results : int
// 	The maximum number of search results to be returned by a request. A
// 	number between 10 and the system limit (currently 500). By default,
// 	a request response will return 10 results.
// next_token : str
// 	This parameter is used to get the next 'page' of results. The value
// 	used with the parameter is pulled directly from the response
// 	provided by the API, and should not be modified. You can learn more
// 	by visiting our page on `pagination`_.
// since_id : Union[int, str]
// 	Returns results with a Tweet ID greater than (for example, more
// 	recent than) the specified ID. The ID specified is exclusive and
// 	responses will not include it. If included with the same request as
// 	a ``start_time`` parameter, only ``since_id`` will be used.
// start_time : Union[datetime.datetime, str]
// 	YYYY-MM-DDTHH:mm:ssZ (ISO 8601/RFC 3339). The oldest UTC timestamp
// 	from which the Tweets will be provided. Timestamp is in second
// 	granularity and is inclusive (for example, 12:00:01 includes the
// 	first second of the minute). By default, a request will return
// 	Tweets from up to 30 days ago if you do not include this parameter.
// until_id: string
// 	Returns results with a Tweet ID less than (that is, older than) the
// 	specified ID. Used with ``since_id``. The ID specified is exclusive
// 	and responses will not include it.
//
// References
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/search/api-reference/get-tweets-search-all
//
// Academic Research product track: https://developer.twitter.com/en/docs/projects/overview#product-track
//
// Tweet cap: https://developer.twitter.com/en/docs/projects/overview#tweet-cap
//
// pagination: https://developer.twitter.com/en/docs/twitter-api/tweets/search/integrate/paginate
func (c *Client) SearchAllTweets(query string, params map[string]interface{}) (*TweetsResponse, error) {
	endpoint_parameters := []string{
		"end_time", "expansions", "max_results", "media.fields",
		"next_token", "place.fields", "poll.fields", "query",
		"since_id", "start_time", "tweet.fields", "until_id",
		"user.fields",
	}
	route := "tweets/search/all"
	if params == nil {
		params = make(map[string]interface{})
	} else if val, ok := params["pagination_token"]; ok {
		params["next_token"] = val
		delete(params, "pagination_token")
	}

	params["query"] = query

	response, err := c.get_request(route, OAuth_2, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: query, Params: params}
	tweets := &TweetsResponse{Caller: c.GetUserTweets, CallerData: caller_data}

	return tweets.Parse(response)
}

// The recent search endpoint returns Tweets from the last seven days that match a search query.
//
// The Tweets returned by this endpoint count towards the Project-level
// `Tweet cap`.
//
// Parameters
//
// query : str
// 	One rule for matching Tweets. If you are using a
// 	`Standard Project`_ at the Basic `access level`_, you can use the
// 	basic set of `operators`_ and can make queries up to 512 characters
// 	long. If you are using an `Academic Research Project`_ at the Basic
// 	access level, you can use all available operators and can make
// 	queries up to 1,024 characters long.
// end_time : Union[datetime.datetime, str]
// 	YYYY-MM-DDTHH:mm:ssZ (ISO 8601/RFC 3339). The newest, most recent
// 	UTC timestamp to which the Tweets will be provided. Timestamp is in
// 	second granularity and is exclusive (for example, 12:00:01 excludes
// 	the first second of the minute). By default, a request will return
// 	Tweets from as recent as 30 seconds ago if you do not include this
// 	parameter.
// max_results : int
// 	The maximum number of search results to be returned by a request. A
// 	number between 10 and 100. By default, a request response will
// 	return 10 results.
// next_token : str
// 	This parameter is used to get the next 'page' of results. The value
// 	used with the parameter is pulled directly from the response
// 	provided by the API, and should not be modified.
// since_id : Union[int, str]
// 	Returns results with a Tweet ID greater than (that is, more recent
// 	than) the specified ID. The ID specified is exclusive and responses
// 	will not include it. If included with the same request as a
// 	``start_time`` parameter, only ``since_id`` will be used.
// start_time : Union[datetime.datetime, str]
// 	YYYY-MM-DDTHH:mm:ssZ (ISO 8601/RFC 3339). The oldest UTC timestamp
// 	(from most recent seven days) from which the Tweets will be
// 	provided. Timestamp is in second granularity and is inclusive (for
// 	example, 12:00:01 includes the first second of the minute). If
// 	included with the same request as a ``since_id`` parameter, only
// 	``since_id`` will be used. By default, a request will return Tweets
// 	from up to seven days ago if you do not include this parameter.
// until_id : Union[int, str]
// 	Returns results with a Tweet ID less than (that is, older than) the
// 	specified ID. The ID specified is exclusive and responses will not
// 	include it.
//
// References
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/search/api-reference/get-tweets-search-recent
//
// Tweet cap: https://developer.twitter.com/en/docs/projects/overview#tweet-cap
//
// Standard Project: https://developer.twitter.com/en/docs/projects
//
// access level: https://developer.twitter.com/en/products/twitter-api/early-access/guide.html#na_1
//
// operators: https://developer.twitter.com/en/docs/twitter-api/tweets/search/integrate/build-a-query
//
// Academic Research Project: https://developer.twitter.com/en/docs/projects
func (c *Client) SearchRecentTweets(query string, params map[string]interface{}) (*TweetsResponse, error) {
	endpoint_parameters := []string{
		"end_time", "expansions", "max_results", "media.fields",
		"next_token", "place.fields", "poll.fields", "query",
		"since_id", "start_time", "tweet.fields", "until_id",
		"user.fields",
	}
	route := "tweets/search/recent"
	if params == nil {
		params = make(map[string]interface{})
	} else if val, ok := params["pagination_token"]; ok {
		params["next_token"] = val
		delete(params, "pagination_token")
	}

	params["query"] = query

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: query, Params: params}
	tweets := &TweetsResponse{Caller: c.GetUserTweets, CallerData: caller_data}

	return tweets.Parse(response)
}

// ** Timelines ** //

// Returns Tweets composed by a single user, specified by the requested
// user ID. By default, the most recent ten Tweets are returned per
// request. Using pagination, the most recent 3,200 Tweets can be
// retrieved.
//
// The Tweets returned by this endpoint count towards the Project-level `Tweet cap`.
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/timelines/api-reference/get-users-id-tweets
func (c *Client) GetUserTweets(user_id string, params map[string]interface{}) (*TweetsResponse, error) {
	endpoint_parameters := []string{
		"end_time", "exclude", "expansions", "max_results",
		"media.fields", "pagination_token", "place.fields",
		"poll.fields", "since_id", "start_time", "tweet.fields",
		"until_id", "user.fields",
	}
	route := fmt.Sprintf("users/%s/tweets", user_id)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: user_id, Params: params}
	tweets := &TweetsResponse{Caller: c.GetUserTweets, CallerData: caller_data}

	return tweets.Parse(response)
}

// Returns Tweets mentioning a single user specified by the requested user
// ID. By default, the most recent ten Tweets are returned per request.
// Using pagination, up to the most recent 800 Tweets can be retrieved.
//
// The Tweets returned by this endpoint count towards the Project-level `Tweet cap`.
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/timelines/api-reference/get-users-id-mentions
func (c *Client) GetUserMentions(user_id string, params map[string]interface{}) (*TweetsResponse, error) {
	endpoint_parameters := []string{
		"end_time", "expansions", "max_results", "media.fields",
		"pagination_token", "place.fields", "poll.fields", "since_id",
		"start_time", "tweet.fields", "until_id", "user.fields",
	}

	route := fmt.Sprintf("users/%s/mentions", user_id)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: user_id, Params: params}
	tweets := &TweetsResponse{Caller: c.GetUserMentions, CallerData: caller_data}

	return tweets.Parse(response)
}

// ** Tweet counts ** //

// This endpoint is only available to those users who have been approved
// for the `Academic Research product track`_.
//
// The full-archive search endpoint returns the complete history of public
// Tweets matching a search query; since the first Tweet was created March
// 26, 2006.
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/counts/api-reference/get-tweets-counts-all
func (c *Client) GetAllTweetsCount(query string, params map[string]interface{}) (*TweetsCountResponse, error) {
	endpoint_parameters := []string{
		"end_time", "granularity", "next_token", "query",
		"since_id", "start_time", "until_id",
	}

	if params == nil {
		params = make(map[string]interface{})
	}

	params["query"] = query

	route := "tweets/counts/all"

	response, err := c.get_request(route, OAuth_2, params, endpoint_parameters)

	if err != nil {
		return nil, err
	}

	return (&TweetsCountResponse{}).Parse(response)
}

// The recent Tweet counts endpoint returns count of Tweets from the last
// seven days that match a search query.
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/counts/api-reference/get-tweets-counts-recent
func (c *Client) GetRecentTweetsCount(query string, params map[string]interface{}) (*TweetsCountResponse, error) {
	endpoint_parameters := []string{
		"end_time", "granularity", "query",
		"since_id", "start_time", "until_id",
	}
	if params == nil {
		params = make(map[string]interface{})
	}
	params["query"] = query

	response, err := c.get_request("tweets/counts/recent", OAuth_2, params, endpoint_parameters)

	if err != nil {
		return nil, err
	}

	return (&TweetsCountResponse{}).Parse(response)
}

// ** Tweet lookup ** //

// Returns a variety of information about a single Tweet specified by
// the requested ID.
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/lookup/api-reference/get-tweets-id
func (c *Client) GetTweet(tweet_id string, params map[string]interface{}) (*TweetResponse, error) {
	endpoint_parameters := []string{
		"expansions", "media.fields", "place.fields",
		"poll.fields", "tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("tweets/%s", tweet_id)
	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&TweetResponse{}).Parse(response)
}

// Returns a variety of information about the Tweet specified by the
// requested ID or list of IDs.
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/lookup/api-reference/get-tweets
func (c *Client) GetTweets(tweet_ids []string, params map[string]interface{}) (*TweetsResponse, error) {
	endpoint_parameters := []string{
		"ids", "expansions", "media.fields", "place.fields",
		"poll.fields", "tweet.fields", "user.fields",
	}
	if params == nil {
		params = make(map[string]interface{})
	}
	params["ids"] = tweet_ids
	response, err := c.get_request("tweets", c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	// TODO: Needs pagination => It gets an [] of strings, not string like others!
	// But it seems this doesn't need pagination, it doesn't accept a pagination token.
	return (&TweetsResponse{}).Parse(response)
}

// ** Blocks ** //
func (c *Client) Block(target_user_id string) (*BlockResponse, error) {
	data := map[string]interface{}{
		"target_user_id": target_user_id,
	}

	route := fmt.Sprintf("users/%s/blocking", c.userID)

	response, err := c.request(
		"POST",
		route,
		data,
	)

	if err != nil {
		return nil, err
	}

	return (&BlockResponse{}).Parse(response)
}

// The request succeeds with no action when the user sends a request to a
// user they're not blocking or have already unblocked.
//
// https://developer.twitter.com/en/docs/twitter-api/users/blocks/api-reference/delete-users-user_id-blocking
func (c *Client) UnBlock(target_user_id string) (*BlockResponse, error) {
	route := fmt.Sprintf("users/%s/blocking/%s", c.userID, target_user_id)

	response, err := c.delete_request(route)

	if err != nil {
		return nil, err
	}

	return (&BlockResponse{}).Parse(response)
}

// Returns a list of users who are blocked by the authenticating user.
//
// https://developer.twitter.com/en/docs/twitter-api/users/blocks/api-reference/get-users-blocking
func (c *Client) GetBlocked(params map[string]interface{}) (*UsersResponse, error) {
	// TODO: PAGINATION PROBLEM
	endpoint_parameters := []string{
		"expansions", "max_results", "tweet.fields",
		"user.fields", "pagination_token",
	}

	route := fmt.Sprintf("users/%s/blocking", c.userID)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, OAuth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	return (&UsersResponse{}).Parse(response)
}

// ** Follows ** //

// Allows a user ID to follow another user.
//
// If the target user does not have public Tweets, this endpoint will send
// a follow request.
//
// The request succeeds with no action when the authenticated user sends a
// request to a user they're already following, or if they're sending a
// follower request to a user that does not have public Tweets.
//
// https://developer.twitter.com/en/docs/twitter-api/users/follows/api-reference/post-users-source_user_id-following
func (c *Client) FollowUser(target_user_id string, params map[string]interface{}) (*FollowResponse, error) {
	data := map[string]interface{}{
		"target_user_id": target_user_id,
	}

	route := fmt.Sprintf("users/%s/following", c.userID)

	response, err := c.request(
		"POST",
		route,
		data,
	)

	if err != nil {
		return nil, err
	}

	return (&FollowResponse{}).Parse(response)
}

// Allows a user ID to unfollow another user.
//
// The request succeeds with no action when the authenticated user sends a
// request to a user they're not following or have already unfollowed.
//
// https://developer.twitter.com/en/docs/twitter-api/users/follows/api-reference/delete-users-source_id-following
func (c *Client) UnfollowUser(target_user_id string) (*FollowResponse, error) {
	route := fmt.Sprintf("users/%s/following/%s", c.userID, target_user_id)

	response, err := c.delete_request(route)
	if err != nil {
		return nil, err
	}

	return (&FollowResponse{}).Parse(response)
}

// Returns a list of users who are followers of the specified user ID.
//
// https://developer.twitter.com/en/docs/twitter-api/users/follows/api-reference/get-users-id-followers
func (c *Client) GetUserFollowers(user_id string, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "tweet.fields",
		"user.fields", "pagination_token",
	}

	route := fmt.Sprintf("users/%s/followers", user_id)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: user_id, Params: params}
	users := &UsersResponse{Caller: c.GetUserFollowers, CallerData: caller_data}

	return users.Parse(response)
}

// Returns a list of users the specified user ID is following
//
// https://developer.twitter.com/en/docs/twitter-api/users/follows/api-reference/get-users-id-following
func (c *Client) GetUserFollowing(user_id string, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "tweet.fields",
		"user.fields", "pagination_token",
	}

	route := fmt.Sprintf("users/%s/following", user_id)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: user_id, Params: params}
	users := &UsersResponse{Caller: c.GetUserFollowing, CallerData: caller_data}

	return users.Parse(response)
}

// ** Mutes ** //
// Allows an authenticated user ID to mute the target user.
//
// https://developer.twitter.com/en/docs/twitter-api/users/mutes/api-reference/post-users-user_id-muting
func (c *Client) Mute(target_user_id string) (*MuteResponse, error) {
	data := map[string]interface{}{
		"target_user_id": target_user_id,
	}

	route := fmt.Sprintf("users/%s/muting", c.userID)

	response, err := c.request(
		"POST",
		route,
		data,
	)

	if err != nil {
		return nil, err
	}

	return (&MuteResponse{}).Parse(response)
}

// Allows an authenticated user ID to unmute the target user.
//
// The request succeeds with no action when the user sends a request to a
// user they're not muting or have already unmuted.
//
// https://developer.twitter.com/en/docs/twitter-api/users/mutes/api-reference/delete-users-user_id-muting
func (c *Client) UnMute(target_user_id string) (*MuteResponse, error) {
	route := fmt.Sprintf("users/%s/muting/%s", c.userID, target_user_id)

	response, err := c.delete_request(route)
	if err != nil {
		return nil, err
	}

	return (&MuteResponse{}).Parse(response)
}

// Returns a list of users who are muted by the authenticating user.
//
// https://developer.twitter.com/en/docs/twitter-api/users/mutes/api-reference/get-users-muting
func (c *Client) GetMuted(params map[string]interface{}) (*MutedUsersResponse, error) {
	// TODO: PAGINATION PROBLEM
	endpoint_parameters := []string{
		"expansions", "max_results", "tweet.fields",
		"user.fields", "pagination_token",
	}
	route := fmt.Sprintf("users/%s/muting", c.userID)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, OAuth_2, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: "", Params: params}
	users := &MutedUsersResponse{Caller: c.GetMuted, CallerData: caller_data}

	return users.Parse(response)
}

// ** User lookup ** //

// Returns information about an authorized user.
//
// https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-me
func (c *Client) GetMe(oauth_1a bool, params map[string]interface{}) (*UserResponse, error) {
	endpoint_parameters := []string{
		"expansions", "tweet.fields", "user.fields",
	}

	route := "users/me"

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	return (&UserResponse{}).Parse(response)
}

// Returns a variety of information about a single user specified by the
// requested ID.
//
// https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-id
func (c *Client) GetUserByID(user_id string, params map[string]interface{}) (*UserResponse, error) {
	endpoint_parameters := []string{
		"expansions", "tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("users/%s", user_id)
	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UserResponse{}).Parse(response)
}

// Returns a variety of information about a single user specified by the
// requested username.
//
// https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-by-username-username
func (c *Client) GetUserByUsername(username string, params map[string]interface{}) (*UserResponse, error) {
	endpoint_parameters := []string{
		"expansions", "tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("users/by/username/%s", username)
	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UserResponse{}).Parse(response)
}

// Returns a variety of information about one or more users specified by
// the requested IDs.
//
// https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users
func (c *Client) GetUsersByIDs(user_ids []string, params map[string]interface{}) (*UsersResponse, error) {
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

	response, err := c.get_request("users", c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UsersResponse{}).Parse(response)
}

// Returns a variety of information about one or more users specified by
// the requested usernames.
//
// https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-by
func (c *Client) GetUsersByUsernames(usernames []string, params map[string]interface{}) (*UsersResponse, error) {
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

	response, err := c.get_request("users/by", c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UsersResponse{}).Parse(response)
}

// ** Spaces ** //

// Return live or scheduled Spaces matching your specified search terms
//
// https://developer.twitter.com/en/docs/twitter-api/spaces/search/api-reference/get-spaces-search
func (c *Client) SearchSpaces(query string, params map[string]interface{}) (*SpacesResponse, error) {
	endpoint_parameters := []string{
		"query", "expansions", "max_results",
		"space.fields", "state", "user.fields",
	}
	route := "spaces/search"
	if params == nil {
		params = make(map[string]interface{})
	}
	params["query"] = query
	response, err := c.get_request(route, OAuth_2, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&SpacesResponse{}).Parse(response)
}

// Returns details about multiple live or scheduled Spaces.
//
// Up to 100 comma-separated Space IDs can be looked up using this method.
//
// https://developer.twitter.com/en/docs/twitter-api/spaces/lookup/api-reference/get-spaces
func (c *Client) GetSpacesBySpaceIDs(space_ids []string, params map[string]interface{}) (*SpacesResponse, error) {
	endpoint_parameters := []string{
		"ids", "user_ids", "expansions", "space.fields", "user.fields",
	}
	route := "spaces"
	if params == nil {
		params = make(map[string]interface{})
	}
	params["ids"] = space_ids
	response, err := c.get_request(route, OAuth_2, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&SpacesResponse{}).Parse(response)
}

// Returns details about multiple live or scheduled Spaces created by the
// specified user IDs.
// Up to 100 comma-separated user IDs can be looked up using this method.
//
// https://developer.twitter.com/en/docs/twitter-api/spaces/lookup/api-reference/get-spaces-by-creator-ids
func (c *Client) GetSpacesByCreatorIDs(creator_ids []string, params map[string]interface{}) (*SpacesResponse, error) {
	endpoint_parameters := []string{
		"ids", "user_ids", "expansions", "space.fields", "user.fields",
	}
	route := "spaces/by/creator_ids"
	if params == nil {
		params = make(map[string]interface{})
	}
	params["user_ids"] = creator_ids
	response, err := c.get_request(route, OAuth_2, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&SpacesResponse{}).Parse(response)
}

// Returns a variety of information about a single Space specified by the
// requested ID.
//
// https://developer.twitter.com/en/docs/twitter-api/spaces/lookup/api-reference/get-spaces-id
func (c *Client) GetSpace(space_id string, params map[string]interface{}) (*SpaceResponse, error) {
	endpoint_parameters := []string{
		"expansions", "space.fields", "user.fields",
	}
	route := fmt.Sprintf("spaces/%s", space_id)
	response, err := c.get_request(route, OAuth_2, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&SpaceResponse{}).Parse(response)
}

// Returns a list of user who purchased a ticket to the requested Space.
// You must authenticate the request using the Access Token of the creator
// of the requested Space.
//
// https://developer.twitter.com/en/docs/twitter-api/spaces/lookup/api-reference/get-spaces-id-buyers
func (c *Client) GetSpaceBuyers(space_id string, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"expansions", "media.fields", "place.fields",
		"poll.fields", "tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("spaces/%s/buyers", space_id)
	response, err := c.get_request(route, OAuth_2, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&UsersResponse{}).Parse(response)
}

// Returns Tweets shared in the requested Spaces.
//
// https://developer.twitter.com/en/docs/twitter-api/spaces/lookup/api-reference/get-spaces-id-tweets
func (c *Client) GetSpaceTweets(space_id string, params map[string]interface{}) (*TweetsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "media.fields", "place.fields",
		"poll.fields", "tweet.fields", "user.fields",
	}

	route := fmt.Sprintf("spaces/%s/tweets", space_id)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, OAuth_2, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	return (&TweetsResponse{}).Parse(response)
}

// ** List Tweets lookup ** //

// Returns a list of Tweets from the specified List.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/list-tweets/api-reference/get-lists-id-tweets
func (c *Client) GetListTweets(list_id string, params map[string]interface{}) (*TweetsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "pagination_token",
		"tweet.fields", "user.fields",
	}
	route := fmt.Sprintf("lists/%s/tweets", list_id)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: list_id, Params: params}
	tweets := &TweetsResponse{Caller: c.GetListTweets, CallerData: caller_data}

	return tweets.Parse(response)
}

// ** List follows ** //

// Enables the authenticated user to follow a List.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/list-follows/api-reference/post-users-id-followed-lists
func (c *Client) FollowList(list_id string) (*FollowResponse, error) {
	data := map[string]interface{}{
		"list_id": list_id,
	}

	route := fmt.Sprintf("users/%s/followed_lists", c.userID)

	response, err := c.request(
		"POST",
		route,
		data,
	)

	if err != nil {
		return nil, err
	}

	return (&FollowResponse{}).Parse(response)

}

// Enables the authenticated user to unfollow a List.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/list-follows/api-reference/delete-users-id-followed-lists-list_id
func (c *Client) UnfollowList(list_id string) (*FollowResponse, error) {
	route := fmt.Sprintf("users/%s/followed_lists/%s", c.userID, list_id)

	response, err := c.delete_request(route)
	if err != nil {
		return nil, err
	}

	return (&FollowResponse{}).Parse(response)
}

// Returns a list of users who are followers of the specified List.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/list-follows/api-reference/get-lists-id-followers
func (c *Client) GetListFollowers(list_id string, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "pagination_token",
		"tweet.fields", "user.fields",
	}

	route := fmt.Sprintf("lists/%s/followers", list_id)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: list_id, Params: params}
	users := &UsersResponse{Caller: c.GetListFollowers, CallerData: caller_data}

	return users.Parse(response)
}

// Returns all Lists a specified user follows.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/list-follows/api-reference/get-users-id-followed_lists
func (c *Client) GetFollowedLists(user_id string, params map[string]interface{}) (*ListsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "pagination_token",
		"list.fields", "user.fields",
	}

	route := fmt.Sprintf("users/%s/followed_lists", user_id)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: user_id, Params: params}
	lists := &ListsResponse{Caller: c.GetFollowedLists, CallerData: caller_data}

	return lists.Parse(response)
}

// ** List lookup ** //

// Returns the details of a specified List.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/list-lookup/api-reference/get-lists-id
func (c *Client) GetList(list_id string, params map[string]interface{}) (*ListResponse, error) {
	endpoint_parameters := []string{
		"expansions", "list.fields", "user.fields",
	}
	route := fmt.Sprintf("lists/%s", list_id)
	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&ListResponse{}).Parse(response)
}

// Returns all Lists owned by the specified user.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/list-lookup/api-reference/get-users-id-owned_lists
func (c *Client) GetOwnedLists(user_id string, params map[string]interface{}) (*ListsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "pagination_token",
		"list.fields", "user.fields",
	}

	route := fmt.Sprintf("users/%s/owned_lists", user_id)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: user_id, Params: params}
	lists := &ListsResponse{Caller: c.GetOwnedLists, CallerData: caller_data}

	return lists.Parse(response)
}

// ** List members ** //

// Enables the authenticated user to add a member to a List they own.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/list-members/api-reference/post-lists-id-members
func (c *Client) AddListMemeber(list_id, user_id string) (*ListMemberResponse, error) {
	data := map[string]interface{}{
		"user_id": user_id,
	}

	route := fmt.Sprintf("lists/%s/members", list_id)

	response, err := c.request(
		"POST",
		route,
		data,
	)

	if err != nil {
		return nil, err
	}

	return (&ListMemberResponse{}).Parse(response)
}

// Enables the authenticated user to remove a member from a List they own.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/list-members/api-reference/delete-lists-id-members-user_id
func (c *Client) RemoveListMember(list_id, user_id string) (*ListMemberResponse, error) {
	route := fmt.Sprintf("lists/%s/members/%s", list_id, user_id)

	response, err := c.delete_request(route)
	if err != nil {
		return nil, err
	}

	return (&ListMemberResponse{}).Parse(response)
}

// Returns a list of users who are members of the specified List.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/list-members/api-reference/get-lists-id-members
func (c *Client) GetListMembers(list_id string, params map[string]interface{}) (*UsersResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "pagination_token",
		"tweet.fields", "user.fields",
	}

	route := fmt.Sprintf("lists/%s/members", list_id)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: list_id, Params: params}
	users := &UsersResponse{Caller: c.GetListMembers, CallerData: caller_data}

	return users.Parse(response)
}

// Returns all Lists a specified user is a member of.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/list-members/api-reference/get-users-id-list_memberships
func (c *Client) GetListMemberships(user_id string, params map[string]interface{}) (*ListsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "pagination_token",
		"list.fields", "user.fields",
	}

	route := fmt.Sprintf("users/%s/list_memberships", user_id)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, c.oauth_type, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: user_id, Params: params}
	lists := &ListsResponse{Caller: c.GetListMemberships, CallerData: caller_data}

	return lists.Parse(response)
}

// ** Manage Lists ** //

// Enables the authenticated user to create a List.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/manage-lists/api-reference/post-lists
func (c *Client) CreateList(name string, description string, private bool, params map[string]interface{}) (*ListResponse, error) {
	data := map[string]interface{}{
		"name":        name,
		"description": description,
		"private":     private,
	}

	route := "lists"

	response, err := c.request(
		"POST",
		route,
		data,
	)

	if err != nil {
		return nil, err
	}

	return (&ListResponse{}).Parse(response)
}

// Enables the authenticated user to update the meta data of a
// specified List that they own.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/manage-lists/api-reference/put-lists-id
func (c *Client) UpdateList(list_id string, name string, description string, private bool, params map[string]interface{}) (*UpdateListResponse, error) {
	data := map[string]interface{}{
		"name":        name,
		"description": description,
		"private":     private,
	}

	route := fmt.Sprintf("lists/%s", list_id)

	response, err := c.request("PUT", route, data)

	if err != nil {
		return nil, err
	}

	return (&UpdateListResponse{}).Parse(response)
}

// Enables the authenticated user to delete a List that they own.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/manage-lists/api-reference/delete-lists-id
func (c *Client) DeleteList(list_id string) (*DeleteResponse, error) {
	route := fmt.Sprintf("lists/%s", list_id)

	response, err := c.delete_request(route)
	if err != nil {
		return nil, err
	}

	return (&DeleteResponse{}).Parse(response)
}

// ** Pinned Lists ** //

// Enables the authenticated user to pin a List.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/pinned-lists/api-reference/post-users-id-pinned-lists
func (c *Client) PinList(list_id string) (*PinResponse, error) {
	data := map[string]interface{}{
		"list_id": list_id,
	}

	route := fmt.Sprintf("users/%s/pinned_lists", c.userID)

	response, err := c.request(
		"POST",
		route,
		data,
	)

	if err != nil {
		return nil, err
	}

	return (&PinResponse{}).Parse(response)
}

// Enables the authenticated user to unpin a List.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/pinned-lists/api-reference/delete-users-id-pinned-lists-list_id
func (c *Client) UnpinList(list_id string) (*PinResponse, error) {
	route := fmt.Sprintf("users/%s/pinned_lists/%s", c.userID, list_id)

	response, err := c.delete_request(route)
	if err != nil {
		return nil, err
	}

	return (&PinResponse{}).Parse(response)
}

// Returns the Lists pinned by a specified user.
//
// https://developer.twitter.com/en/docs/twitter-api/lists/pinned-lists/api-reference/get-users-id-pinned_lists
func (c *Client) GetPinnedLists(params map[string]interface{}) (*ListsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "list.fields", "user.fields",
	}
	route := fmt.Sprintf("users/%s/pinned_lists", c.userID)
	response, err := c.get_request(route, OAuth_1a, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}
	return (&ListsResponse{}).Parse(response)
}

// ** Batch Compliance ** //

// Creates a new compliance job for Tweet IDs or user IDs.
//
// A compliance job will contain an ID and a destination URL. The
// destination URL represents the location that contains the list of IDs
// consumed by your app.
//
// You can run one batch job at a time.
//
// https://developer.twitter.com/en/docs/twitter-api/compliance/batch-compliance/api-reference/post-compliance-jobs
func (c *Client) CreateComplianceJob(job_type, name, resumable string) (*ComplianceJobResponse, error) {
	if job_type != "tweets" && job_type != "users" {
		return nil, fmt.Errorf("job_type must be either 'tweets' or 'users'")
	}

	data := map[string]interface{}{
		"type": job_type,
	}
	if name != "" {
		data["name"] = name
	}
	if resumable != "" {
		data["resumable"] = resumable
	}

	route := "compliance/jobs"

	response, err := c.request(
		"POST",
		route,
		data,
	)

	if err != nil {
		return nil, err
	}

	return (&ComplianceJobResponse{}).Parse(response)
}

// Get a single compliance job with the specified ID.
//
// https://developer.twitter.com/en/docs/twitter-api/compliance/batch-compliance/api-reference/get-compliance-jobs-id
func (c *Client) GetComplianceJob(job_id string) (*ComplianceJobResponse, error) {
	route := fmt.Sprintf("compliance/jobs/%s", job_id)

	response, err := c.get_request(route, OAuth_2, nil, nil)

	if err != nil {
		return nil, err
	}

	return (&ComplianceJobResponse{}).Parse(response)
}

// Returns a list of recent compliance jobs.
//
// https://developer.twitter.com/en/docs/twitter-api/compliance/batch-compliance/api-reference/get-compliance-jobs
func (c *Client) GetComplianceJobs(job_type string, params map[string]interface{}) (*ComplianceJobsResponse, error) {
	endpoint_parameters := []string{
		"type", "status",
	}
	if params == nil {
		params = map[string]interface{}{}
	}
	params["type"] = job_type

	route := "compliance/jobs"

	response, err := c.get_request(route, OAuth_2, params, endpoint_parameters)

	if err != nil {
		return nil, err
	}

	return (&ComplianceJobsResponse{}).Parse(response)
}

// ** Bookmarks ** //

// Causes the authenticating user to Bookmark the target Tweet provided
// in the request body.
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/bookmarks/api-reference/post-users-id-bookmarks
func (c *Client) BookmarkTweet(tweet_id string) (*BookmarkResponse, error) {
	data := map[string]interface{}{
		"tweet_id": tweet_id,
	}

	route := fmt.Sprintf("users/%s/bookmarks", c.userID)

	response, err := c.request(
		"POST",
		route,
		data,
	)

	if err != nil {
		return nil, err
	}

	return (&BookmarkResponse{}).Parse(response)
}

// Allows a user or authenticated user ID to remove a Bookmark of a
// Tweet.
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/bookmarks/api-reference/delete-users-id-bookmarks-tweet_id
func (c *Client) RemoveBookmark(tweet_id string) (*BookmarkResponse, error) {
	route := fmt.Sprintf("users/%s/bookmarks/%s", c.userID, tweet_id)

	response, err := c.delete_request(route)
	if err != nil {
		return nil, err
	}

	return (&BookmarkResponse{}).Parse(response)
}

// Allows you to get an authenticated user's 800 most recent bookmarked
// Tweets.
//
// https://developer.twitter.com/en/docs/twitter-api/tweets/bookmarks/api-reference/get-users-id-bookmarks
func (c *Client) GetBookmarkedTweets(params map[string]interface{}) (*BookmarkedTweetsResponse, error) {
	endpoint_parameters := []string{
		"expansions", "max_results", "media.fields",
		"pagination_token", "place.fields", "poll.fields",
		"tweet.fields", "user.fields",
	}

	route := fmt.Sprintf("users/%s/bookmarks", c.userID)

	if params == nil {
		params = make(map[string]interface{})
	}

	response, err := c.get_request(route, OAuth_2, params, endpoint_parameters)
	if err != nil {
		return nil, err
	}

	caller_data := CallerData{ID: "", Params: params}
	tweets := &BookmarkedTweetsResponse{Caller: c.GetBookmarkedTweets, CallerData: caller_data}

	return tweets.Parse(response)
}

// func QueryMaker() string
