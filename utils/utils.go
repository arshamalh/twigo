package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Does an array of strings contain an especial string?
func Contains(arrayOfStrings []string, string_item string) bool {
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

func QueryValue(params []string) string {
	if len(params) == 0 {
		return ""
	}

	return strings.Join(params, ",")
}

func QueryMaker(params map[string]interface{}, endpoint_parameters []string) string {
	parameters := url.Values{}
	for param_name, param_value := range params {
		if new_param_name := strings.Replace(param_name, "_", ".", 1); Contains(endpoint_parameters, new_param_name) {
			param_name = new_param_name
		} else if !Contains(endpoint_parameters, param_name) {
			fmt.Printf("it seems endpoint parameter '%s' is not supported", param_name)
		}
		switch param_valt := param_value.(type) {
		case int:
			parameters.Add(param_name, strconv.Itoa(param_valt))
		case string:
			parameters.Add(param_name, param_valt)
		case []string:
			parameters.Add(param_name, strings.Join(param_valt, ","))
		case time.Time:
			parameters.Add(param_name, param_valt.Format(time.RFC3339))
		default:
			fmt.Printf("%s with value of %s is not supported, please contact us", param_name, param_value)
		}
	}
	return parameters.Encode()
}
