package twigo

import (
	"encoding/json"
	"net/http"
)

type UnLikeResponse struct {
	Data struct {
		Liked bool `json:"liked"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *UnLikeResponse) Parse(raw_response *http.Response) (*UnLikeResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type UnRetweetResponse struct {
	Data struct {
		Retweeted bool `json:"retweeted"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *UnRetweetResponse) Parse(raw_response *http.Response) (*UnRetweetResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type UnBlockResponse struct {
	Data struct {
		Blocking bool `json:"blocking"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *UnBlockResponse) Parse(raw_response *http.Response) (*UnBlockResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type UnFollowResponse struct {
	Data struct {
		Following bool `json:"following"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *UnFollowResponse) Parse(raw_response *http.Response) (*UnFollowResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type UnMuteResponse struct {
	Data struct {
		Muting bool `json:"muting"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *UnMuteResponse) Parse(raw_response *http.Response) (*UnMuteResponse, error) {
	err := json.NewDecoder(raw_response.Body).Decode(&r)
	defer raw_response.Body.Close()
	r.RateLimits.Set(raw_response.Header)
	return r, err
}

type RemoveListMemberResponse struct {
	Data struct {
		IsMember bool `json:"is_member"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *RemoveListMemberResponse) Parse(raw_response *http.Response) (*RemoveListMemberResponse, error) {
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

type UnPinResponse struct {
	Data struct {
		Pinned bool `json:"pinned"`
	}
	Includes   IncludesEntity
	Errors     []ErrorEntity
	Meta       MetaEntity
	RateLimits RateLimits
}

func (r *UnPinResponse) Parse(raw_response *http.Response) (*UnPinResponse, error) {
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
