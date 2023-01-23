package executor

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type EcoWatExecutor struct{}

type Auth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type EcoWatt struct {
	Signals []struct {
		GenerationFichier time.Time `json:"GenerationFichier"`
		Jour              time.Time `json:"jour"`
		Dvalue            int       `json:"dvalue"`
		Message           string    `json:"message"`
		Values            []struct {
			Pas    int `json:"pas"`
			Hvalue int `json:"hvalue"`
		} `json:"values"`
	} `json:"signals"`
}

func (e EcoWatExecutor) Execute() error {
	client := &http.Client{}

	authReq, err := http.NewRequest("POST", "https://digital.iservices.rte-france.com/token/oauth/", nil)
	authReq.Header.Set("Authorization", "Basic NGZlZTBhMjYtNzMyMS00YjEwLWFhMzUtN2NkODUxY2RkYjA4OmE3OTIxZmNiLWVhZDQtNDZhZS1hYWI0LWYyYWM2NTFhODUyMQ==")
	token, err := client.Do(authReq)

	if err != nil {
		fmt.Println("Unable to authenticate", err)
		return err
	}

	defer token.Body.Close()

	tokenBody, err := io.ReadAll(token.Body)
	if err != nil {
		fmt.Println("Unable to read auth response body")
	}

	var auth Auth
	if err := json.Unmarshal(tokenBody, &auth); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON from authentication")
	}

	req, err := http.NewRequest("GET", "https://digital.iservices.rte-france.com/open_api/ecowatt/v4/sandbox/signals", nil)

	if err != nil {
		fmt.Println("Unable to create request for ecowatt", err)
		return err
	}

	req.Header.Set("Authorization", auth.TokenType+" "+auth.AccessToken)
	response, err := client.Do(req)

	if err != nil {
		fmt.Println("Unable to request ecowatt", err)
		return err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body) // response body is []byte
	if err != nil {
		return err
	}

	var ecowatt EcoWatt
	if err := json.Unmarshal(body, &ecowatt); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON from authentication")
	}

	log.Println("response: ", ecowatt.Signals[0].Message)
	return nil
}
