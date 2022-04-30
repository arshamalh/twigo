package twigo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// *** Basic Entities *** //
type MetaEntity struct {
	ResultCount   int    `json:"result_count"`
	NewestID      string `json:"newest_id"`
	OldestID      string `json:"oldest_id"`
	PreviousToken string `json:"previous_token"`
	NextToken     string `json:"next_token"`
}

type ErrorEntity struct {
	Parameters map[string]interface{} `json:"parameters"`
	Message    string                 `json:"message"`
}

type IncludesEntity struct {
	Users  []User  `json:"users"`
	Tweets []Tweet `json:"tweets"`
	Polls  []Poll  `json:"polls"` // Do you remember if "poll" in the client?
	Places []Place `json:"places"`
	Media  []Media `json:"media"`
}

type RateLimits struct {
	Limit          int   `json:"x-rate-limit"`
	Remaining      int   `json:"x-rate-limit-remaining"`
	ResetTimestamp int64 `json:"x-rate-limit-reset"`
}

func (r *RateLimits) Set(header http.Header) {
	r.Limit, _ = strconv.Atoi(header.Get("X-Rate-Limit-Limit"))
	r.Remaining, _ = strconv.Atoi(header.Get("X-Rate-Limit-Remaining"))
	r.ResetTimestamp, _ = strconv.ParseInt(header.Get("X-Rate-Limit-Reset"), 10, 64)
}

type Poll struct{}
type Place struct{}
type Media struct{}
type Tweet struct {
	ID               string   `json:"id"`
	Text             string   `json:"text"`
	CreatedAt        string   `json:"created_at"`
	AuthorID         string   `json:"author_id"`
	ConversationID   string   `json:"conversation_id"`
	InReplyToUserID  string   `json:"in_reply_to_user_id"`
	ReferencedTweets []string `json:"referenced_tweets"`
	Lang             string   `json:"lang"`
	ReplySettings    string   `json:"reply_settings"`
	Source           string   `json:"source"`
	// And more...
}

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	CreatedAt   string `json:"created_at"`
	Protected   bool   `json:"protected"`
	Location    string `json:"location"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Verified    bool   `json:"verified"`
	// And more...
}

type Space struct {
	ID               string   `json:"id"`
	State            string   `json:"state"` // It's a enum actually, not a string, so maybe we should parse it
	HostIDs          []string `json:"host_ids"`
	CreatedAt        string   `json:"created_at"`
	CreatorID        string   `json:"creator_id"`
	EndedAt          string   `json:"ended_at"`
	Lang             string   `json:"lang"`
	IsTicketed       bool     `json:"is_ticketed"`
	InvitedUserIDs   []string `json:"invited_user_ids"`
	ParticipantCount int      `json:"participant_count"`
	ScheduledStart   string   `json:"scheduled_start"`
	SpeakerIDs       []string `json:"speaker_ids"`
	StartedAt        string   `json:"started_at"`
	SubscriberCount  int      `json:"subscriber_count"`
	// And more...
}
type List struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	CreatedAt     string `json:"created_at"`
	Private       bool   `json:"private"`
	FollowerCount int    `json:"follower_count"`
	MemberCount   int    `json:"member_count"`
	OwnerID       string `json:"owner_id"`
	Description   string `json:"description"`
}

// *** Request struct *** //
type CallerData struct {
	ID       string
	OAuth_1a bool
	Params   map[string]interface{}
}

// *** Response Entities *** //
type TweetResponse struct {
	Data       Tweet
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *TweetResponse) Parse(raw_response *http.Response) (*TweetResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type TweetsResponse struct {
	Data       []Tweet
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
	CallerData CallerData
	Caller     func(string, bool, map[string]interface{}) (*TweetsResponse, error)
}

func (r *TweetsResponse) Parse(raw_response *http.Response) (*TweetsResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

func (r *TweetsResponse) NextPage() (*TweetsResponse, error) {
	if r.Meta.NextToken == "" {
		return nil, fmt.Errorf("no next page")
	}
	r.CallerData.Params["pagination_token"] = r.Meta.NextToken
	return r.Caller(r.CallerData.ID, r.CallerData.OAuth_1a, r.CallerData.Params)
}

type UserResponse struct {
	Data       User
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *UserResponse) Parse(raw_response *http.Response) (*UserResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type UsersResponse struct {
	Data       []User
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
	CallerData CallerData
	Caller     func(string, bool, map[string]interface{}) (*UsersResponse, error)
}

func (r *UsersResponse) Parse(raw_response *http.Response) (*UsersResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

func (r *UsersResponse) NextPage() (*UsersResponse, error) {
	if r.Meta.NextToken == "" {
		return nil, fmt.Errorf("no next page")
	}
	r.CallerData.Params["pagination_token"] = r.Meta.NextToken
	return r.Caller(r.CallerData.ID, r.CallerData.OAuth_1a, r.CallerData.Params)
}

type MutedUsersResponse struct {
	Data       []User
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
	CallerData CallerData
	Caller     func(map[string]interface{}) (*MutedUsersResponse, error)
}

func (r *MutedUsersResponse) Parse(raw_response *http.Response) (*MutedUsersResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

func (r *MutedUsersResponse) NextPage() (*MutedUsersResponse, error) {
	if r.Meta.NextToken == "" {
		return nil, fmt.Errorf("no next page")
	}
	r.CallerData.Params["pagination_token"] = r.Meta.NextToken
	return r.Caller(r.CallerData.Params)
}

type SpaceResponse struct {
	Data       Space
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *SpaceResponse) Parse(raw_response *http.Response) (*SpaceResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type SpacesResponse struct {
	Data       Space
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *SpacesResponse) Parse(raw_response *http.Response) (*SpacesResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type ListResponse struct {
	Data       List
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *ListResponse) Parse(raw_response *http.Response) (*ListResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type ListsResponse struct {
	Data       []List
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
	CallerData CallerData
	Caller     func(string, bool, map[string]interface{}) (*ListsResponse, error)
}

func (r *ListsResponse) Parse(raw_response *http.Response) (*ListsResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

func (r *ListsResponse) NextPage() (*ListsResponse, error) {
	if r.Meta.NextToken == "" {
		return nil, fmt.Errorf("no next page")
	}
	r.CallerData.Params["pagination_token"] = r.Meta.NextToken
	return r.Caller(r.CallerData.ID, r.CallerData.OAuth_1a, r.CallerData.Params)
}
