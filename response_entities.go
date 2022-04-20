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
	NextToken   string `json:"next_token"`
}

type ErrorEntity struct {
	Parameters map[string]interface{} `json:"parameters"`
	Message    string                 `json:"message"`
}

type Tweet struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Space struct{}
type List struct{}

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
