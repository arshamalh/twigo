package twigo

import (
	"encoding/json"
	"net/http"
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
