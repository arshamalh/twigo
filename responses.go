package twigo

import (
	"encoding/json"
	"net/http"
	"time"
)

type LikeResponse struct {
	Data struct {
		Liked bool `json:"liked"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *LikeResponse) Parse(raw_response *http.Response) (*LikeResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type HideReplyResponse struct {
	Data struct {
		Hidden bool `json:"hidden"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *HideReplyResponse) Parse(raw_response *http.Response) (*HideReplyResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type RetweetResponse struct {
	Data struct {
		Retweeted bool `json:"retweeted"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *RetweetResponse) Parse(raw_response *http.Response) (*RetweetResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type BlockResponse struct {
	Data struct {
		Blocking bool `json:"blocking"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *BlockResponse) Parse(raw_response *http.Response) (*BlockResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type FollowResponse struct {
	Data struct {
		Following bool `json:"following"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *FollowResponse) Parse(raw_response *http.Response) (*FollowResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type MuteResponse struct {
	Data struct {
		Muting bool `json:"muting"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *MuteResponse) Parse(raw_response *http.Response) (*MuteResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type ListMemberResponse struct {
	Data struct {
		IsMember bool `json:"is_member"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *ListMemberResponse) Parse(raw_response *http.Response) (*ListMemberResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type DeleteResponse struct {
	Data struct {
		Deleted bool `json:"deleted"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *DeleteResponse) Parse(raw_response *http.Response) (*DeleteResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type PinResponse struct {
	Data struct {
		Pinned bool `json:"pinned"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *PinResponse) Parse(raw_response *http.Response) (*PinResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type BookmarkResponse struct {
	Data struct {
		Bookmarked bool `json:"bookmarked"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *BookmarkResponse) Parse(raw_response *http.Response) (*BookmarkResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type UpdateListResponse struct {
	Data struct {
		Updated bool `json:"updated"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *UpdateListResponse) Parse(raw_response *http.Response) (*UpdateListResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type TweetsCountResponse struct {
	Data struct {
		End        time.Time `json:"end"`
		Start      time.Time `json:"start"`
		TweetCount int       `json:"tweet_count"`
	}
	Meta struct {
		// TODO: Also there is a "meta" field in the response that is a object
		TotalTweetCount int    `json:"total_tweet_count"`
		NextToken       string `json:"next_token"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	RateLimits RateLimits
}

func (r *TweetsCountResponse) Parse(raw_response *http.Response) (*TweetsCountResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type ComplianceJobResponse struct {
	Data       ComplianceJob
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *ComplianceJobResponse) Parse(raw_response *http.Response) (*ComplianceJobResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type ComplianceJobsResponse struct {
	Data       []ComplianceJob
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *ComplianceJobsResponse) Parse(raw_response *http.Response) (*ComplianceJobsResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}
