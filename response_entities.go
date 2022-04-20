package twigo

import (
	"encoding/json"
	"net/http"
)

// *** Basic Entities *** //
type MetaEntity struct {
	ResultCount int    `json:"result_count"`
	NewestID    string `json:"newest_id"`
	OldestID    string `json:"oldest_id"`
	PreviousToken string `json:"previous_token"`
	NextToken   string `json:"next_token"`
}

type ErrorEntity struct {
	Parameters map[string]interface{} `json:"parameters"`
	Message    string                 `json:"message"`
}

type IncludesEntity struct {
	Users []User `json:"users"`
}

type Tweet struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	CreatedAt string `json:"created_at"`
	AuthorID string `json:"author_id"`
	ConversationID string `json:"conversation_id"`
	InReplyToUserID string `json:"in_reply_to_user_id"`
	ReferencedTweets []string `json:"referenced_tweets"`
	Lang string `json:"lang"`
	ReplySettings string `json:"reply_settings"`
	Source string `json:"source"`
	// And more...
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	CreatedAt string `json:"created_at"`
	Protected bool `json:"protected"`
	Location string `json:"location"`
	URL string `json:"url"`
	Description string `json:"description"`
	Verified bool `json:"verified"`
	// And more...
}

type Space struct{
	ID string `json:"id"`
	State string `json:"state"` // It's a enum actually, not a string, so maybe we should parse it
	HostIDs []string `json:"host_ids"`
	CreatedAt string `json:"created_at"`
	CreatorID string `json:"creator_id"`
	EndedAt string `json:"ended_at"`
	Lang string `json:"lang"`
	IsTicketed bool `json:"is_ticketed"`
	InvitedUserIDs []string `json:"invited_user_ids"`
	ParticipantCount int `json:"participant_count"`
	ScheduledStart string `json:"scheduled_start"`
	SpeakerIDs []string `json:"speaker_ids"`
	StartedAt string `json:"started_at"`
	SubscriberCount int `json:"subscriber_count"`
	// And more...
}
type List struct{
	ID string `json:"id"`
	Name string `json:"name"`
	CreatedAt string `json:"created_at"`
	Private bool `json:"private"`
	FollowerCount int `json:"follower_count"`
	MemberCount int `json:"member_count"`
	OwnerID string `json:"owner_id"`
	Description string `json:"description"`
}

// *** Response Entities *** //
type TweetResponse struct {
	Data     Tweet
	Includes interface{}
	Errors   []ErrorEntity
	Meta     MetaEntity
}

func (r *TweetResponse) Parse(raw_response *http.Response) (*TweetResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	return r, err
}

type TweetsResponse struct {
	Data     []Tweet
	Includes interface{}
	Errors   []ErrorEntity
	Meta     MetaEntity
}

func (r *TweetsResponse) Parse(raw_response *http.Response) (*TweetsResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	return r, err
}

type UserResponse struct {
	Data     User
	Includes interface{}
	Errors   []ErrorEntity
	Meta     MetaEntity
}

func (r *UserResponse) Parse(raw_response *http.Response) (*UserResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	return r, err
}

type UsersResponse struct {
	Data     []User
	Includes interface{}
	Errors   []ErrorEntity
	Meta     MetaEntity
}

func (r *UsersResponse) Parse(raw_response *http.Response) (*UsersResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	return r, err
}

type SpaceResponse struct {
	Data     Space
	Includes interface{}
	Errors   []ErrorEntity
	Meta     MetaEntity
}

func (r *SpaceResponse) Parse(raw_response *http.Response) (*SpaceResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	return r, err
}

type ListResponse struct {
	Data     List
	Includes interface{}
	Errors   []ErrorEntity
	Meta     MetaEntity
}

func (r *ListResponse) Parse(raw_response *http.Response) (*ListResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	return r, err
}

type ListsResponse struct {
	Data     List
	Includes interface{}
	Errors   []ErrorEntity
	Meta     MetaEntity
}

func (r *ListsResponse) Parse(raw_response *http.Response) (*ListsResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	return r, err
}
