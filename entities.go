package twigo

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/arshamalh/twigo/entities"
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

type SpecialError struct {
	Title  string `json:"title"`
	Type   string `json:"type"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func (e *SpecialError) Error() error {
	return fmt.Errorf("%s - %d - %s - %s", e.Title, e.Status, e.Type, e.Detail)
}

type IncludesEntity struct {
	Users  []entities.User  `json:"users"`
	Tweets []entities.Tweet `json:"tweets"`
	Polls  []Poll           `json:"polls"` // Do you remember if "poll" in the client?
	Places []Place          `json:"places"`
	Media  []Media          `json:"media"`
}

type RateLimits struct {
	Limit          int   `json:"x-rate-limit"`
	Remaining      int   `json:"x-rate-limit-remaining"`
	ResetTimestamp int64 `json:"x-rate-limit-reset"` // Isn't this a time.Time?
}

func (r *RateLimits) Set(header http.Header) {
	r.Limit, _ = strconv.Atoi(header.Get("X-Rate-Limit-Limit"))
	r.Remaining, _ = strconv.Atoi(header.Get("X-Rate-Limit-Remaining"))
	r.ResetTimestamp, _ = strconv.ParseInt(header.Get("X-Rate-Limit-Reset"), 10, 64)
}

type Place struct {
	FullName        string     `json:"full_name"`
	ID              string     `json:"id"`
	ContainedWithin string     `json:"contained_within,omitempty"`
	Country         string     `json:"country,omitempty"`
	CountryCode     string     `json:"country_code,omitempty"`
	Geo             IncludeGeo `json:"geo,omitempty"`
	Name            string     `json:"name,omitempty"`
	PlaceType       string     `json:"place_type,omitempty"`
}

type Media struct {
	MediaKey         string         `json:"media_key"`
	Type             string         `json:"type"`
	Url              string         `json:"url"`
	DurationMs       int            `json:"duration_ms,omitempty"`
	Height           int            `json:"height,omitempty"`
	NonPublicMetrics map[string]int `json:"non_public_metrics,omitempty"`
	OrganicMetrics   map[string]int `json:"organic_metrics,omitempty"`
	PreviewImageUrl  string         `json:"preview_image_url,omitempty"`
	PromotedMetrics  map[string]int `json:"promoted_metrics,omitempty"`
	PublicMetrics    map[string]int `json:"public_metrics,omitempty"`
	Width            int            `json:"width,omitempty"`
	AltText          int            `json:"alt_text,omitempty"`
}

type Poll struct {
	ID              string       `json:"id"`
	Options         []PollOption `json:"options"`
	DurationMinutes int          `json:"duration_minutes,omitempty"`
	EndDatetime     time.Time    `json:"end_datetime,omitempty"`
	VotingStatus    string       `json:"voting_status,omitempty"`
}

type IncludeGeo struct {
	Type       string      `json:"type"`
	BBox       [4]float64  `json:"bbox"`
	Properties interface{} `json:"properties"`
}

type PollOption struct {
	Position int    `json:"position"`
	Label    string `json:"label"`
	Votes    int    `json:"votes"`
}

type List struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	Private       bool      `json:"private,omitempty"`
	FollowerCount int       `json:"follower_count,omitempty"`
	MemberCount   int       `json:"member_count,omitempty"`
	OwnerID       string    `json:"owner_id,omitempty"`
	Description   string    `json:"description,omitempty"`
}
