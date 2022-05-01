package twigo

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Does an array of strings contain an especial string?
func contains(arrayOfStrings []string, string_item string) bool {
	for _, val := range arrayOfStrings {
		if val == string_item {
			return true
		}
	}
	return false
}

type BearerToken struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

// Utility function to get the bearer token using the client_id and client_secret
func BearerFinder(ConsumerKey, ConsumerSecret string) (string, error) {
	credentials := ConsumerKey + ":" + ConsumerSecret
	credentialsBase64Encoded := base64.StdEncoding.EncodeToString([]byte(credentials))

	request, err := http.NewRequest(
		"POST",
		"https://api.twitter.com/oauth2/token",
		strings.NewReader("grant_type=client_credentials"),
	)

	if err != nil {
		return "", err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", credentialsBase64Encoded))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)

	if response.StatusCode != 200 || err != nil {
		return "", fmt.Errorf("error code: %d, error: %s", response.StatusCode, err)
	}

	defer response.Body.Close()

	bearer_token := &BearerToken{}
	err = json.NewDecoder(response.Body).Decode(bearer_token)

	if err != nil {
		return "", err
	}

	return bearer_token.AccessToken, nil
}
