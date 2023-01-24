package authentication

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type AuthResponse struct {
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   time.Duration `json:"expires_in"`
	ExpiresAt   time.Time
}

type Authenticator struct {
	URI  string `default:"https://digital.iservices.rte-france.com/token/oauth/"`
	code string `default:"NGZlZTBhMjYtNzMyMS00YjEwLWFhMzUtN2NkODUxY2RkYjA4OmE3OTIxZmNiLWVhZDQtNDZhZS1hYWI0LWYyYWM2NTFhODUyMQ=="`
}

const TokenExpirationSafetyDelta time.Duration = 10

var once sync.Once
var instance *Authenticator
var authResponse *AuthResponse

func getInstance() *Authenticator {
	once.Do(getSingleInstance)
	return instance
}

func authenticate() error {
	client := &http.Client{}

	request, err := http.NewRequest("POST", getInstance().URI, nil)
	request.Header.Set("Authorization", "Basic "+getInstance().code)
	token, err := client.Do(request)
	if err != nil {
		fmt.Println("Unable to authenticate", err)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Unable to close Body reader ", err)
		}
	}(token.Body)

	body, err := io.ReadAll(token.Body)
	if err != nil {
		fmt.Println("Unable to read auth response body")
	}

	var response AuthResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Can not unmarshal JSON from authentication")
		return err
	}

	response.ExpiresAt = time.Now().Add(time.Second * response.ExpiresIn)
	authResponse = &response
	return nil
}

func GetToken() (string, error) {
	if authResponse == nil || isTokenExpired() {
		err := authenticate()
		if err != nil {
			return "", err
		}
	}
	return authResponse.AccessToken, nil
}

func isTokenExpired() bool {
	return authResponse.ExpiresAt.Before(time.Now().Add(-time.Second * TokenExpirationSafetyDelta))
}

func getSingleInstance() {
	instance = &Authenticator{
		URI:  "https://digital.iservices.rte-france.com/token/oauth/",
		code: "NGZlZTBhMjYtNzMyMS00YjEwLWFhMzUtN2NkODUxY2RkYjA4OmE3OTIxZmNiLWVhZDQtNDZhZS1hYWI0LWYyYWM2NTFhODUyMQ==",
	}
}
