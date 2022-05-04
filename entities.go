package twigo

import (
	"net/http"
	"strconv"
	"time"
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
	ID               string         `json:"id"`
	Text             string         `json:"text"`
	CreatedAt        time.Time      `json:"created_at"`
	AuthorID         string         `json:"author_id"`
	ConversationID   string         `json:"conversation_id"`
	InReplyToUserID  string         `json:"in_reply_to_user_id"`
	ReferencedTweets []string       `json:"referenced_tweets"`
	Lang             string         `json:"lang"`
	ReplySettings    string         `json:"reply_settings"`
	Source           string         `json:"source"`
	PublicMetrics    map[string]int `json:"public_metrics"`
	// And more...
}

type User struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	CreatedAt   time.Time `json:"created_at"`
	Protected   bool      `json:"protected"`
	Location    string    `json:"location"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	Verified    bool      `json:"verified"`
	// And more...
}

type Space struct {
	ID               string    `json:"id"`
	State            string    `json:"state"` // It's a enum actually, not a string, so maybe we should parse it
	HostIDs          []string  `json:"host_ids"`
	CreatedAt        time.Time `json:"created_at"`
	CreatorID        string    `json:"creator_id"`
	EndedAt          string    `json:"ended_at"`
	Lang             string    `json:"lang"`
	IsTicketed       bool      `json:"is_ticketed"`
	InvitedUserIDs   []string  `json:"invited_user_ids"`
	ParticipantCount int       `json:"participant_count"`
	ScheduledStart   string    `json:"scheduled_start"`
	SpeakerIDs       []string  `json:"speaker_ids"`
	StartedAt        time.Time `json:"started_at"`
	SubscriberCount  int       `json:"subscriber_count"`
	// And more...
}
type List struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	Private       bool      `json:"private"`
	FollowerCount int       `json:"follower_count"`
	MemberCount   int       `json:"member_count"`
	OwnerID       string    `json:"owner_id"`
	Description   string    `json:"description"`
}

type ComplianceJob struct {
	ID                string    `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	Status            string    `json:"status"`
	Type              string    `json:"type"`
	Resumable         bool      `json:"resumable"`
	DownloadExpiresAt time.Time `json:"download_expires_at"`
	UploadUrl         string    `json:"upload_url"`
	DownloadUrl       string    `json:"download_url"`
	UploadExpiresAt   time.Time `json:"upload_expires_at"`
}

// *** Request struct *** //
type CallerData struct {
	ID       string
	OAuth_1a bool
	Params   map[string]interface{}
}
