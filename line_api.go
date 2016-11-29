package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
)

type ProfileInfo struct {
	DisplayName      string `json:"displayName,omitempty"`
	Mid              string `json:"mid,omitempty"`
	PictureUrl       string `json:"pictureUrl,omitempty"`
	StatusMessage    string `json:"statusMessage,omitempty"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_descrition,omitempty"`
}

type APICredentials struct {
	AccessToken  string
	RefreshToken string
	Expire       int64
}

type RefreshResponse struct {
	Mid              string `json:"mid,omitempty"`
	AccessToken      string `json:"accessToken,omitempty"`
	RefreshToken     string `json:"refreshToken,omitempty"`
	Expire           int64  `json:"expire,omitempty"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_descrition,omitempty"`
}

type VerifyResponse struct {
	Mid              string `json:"mid,omitempty"`
	ChannelId        int64  `json:"channelId,omitempty"`
	Expire           int64  `json:"expire,omitempty"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_descrition,omitempty"`
}

type MessageTextParameters struct {
	UserName string `json:"user_name,omitempty"`
}

type MessageContent struct {
	To        []string `json:"to,omitempty"`
	ToChannel string   `json:"toChannel,omitempty"`
	EventType string   `json:"eventType,omitempty"`
}

type MessageRequest struct {
	To        []string `json:"to,omitempty"`
	ToChannel string   `json:"toChannel,omitempty"`
	EventType string   `json:"eventType,omitempty"`
}

func LineLogout() {

	credentials := GetCredentials()

	client := &http.Client{}

	urlStr := os.Getenv("LINE_API_BASE_URL") + "/oauth/logout"

	req, err := http.NewRequest("DELETE", urlStr, nil)

	log.Println("Logout URL: " + urlStr)
	log.Println("Logout Access Token: " + credentials.AccessToken)

	req.Header.Add("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

}

func RefreshAccessToken(credentials APICredentials) error {

	client := &http.Client{}

	var refreshResponse RefreshResponse

	data := url.Values{}
	data.Set("refreshToken", credentials.RefreshToken)

	urlStr := os.Getenv("LINE_API_BASE_URL") + "/oauth/accessToken"

	log.Println("Refresh Token to be used: " + credentials.RefreshToken)
	log.Println("Access Token to be used: " + credentials.AccessToken)

	req, err := http.NewRequest("POST", urlStr+"?"+data.Encode(), nil)

	log.Println("request: " + req.URL.String())

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Line-ChannelToken", credentials.AccessToken)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&refreshResponse)

	log.Println("Access Token from Response: " + refreshResponse.AccessToken)
	log.Println("Refresh Token from Response: " + refreshResponse.RefreshToken)

	log.Println("Expire from Response: ", refreshResponse.Expire)

	log.Println("Error: " + refreshResponse.Error)

	newCredentials := APICredentials{}

	if refreshResponse.Error == "" {

		newCredentials.AccessToken = refreshResponse.AccessToken
		newCredentials.RefreshToken = refreshResponse.RefreshToken
		newCredentials.Expire = refreshResponse.Expire

		WriteCredentials(newCredentials)

		// If there are no errors, refresh the credentials in Redis
		return nil

	} else {

		log.Println("Error in refreshresponse: " + refreshResponse.AccessToken)

		returnError := APIError{
			ErrorCode:        refreshResponse.Error,
			ErrorDescription: refreshResponse.ErrorDescription,
		}

		return returnError
	}

	return err

}

func VerifyCredentials(credentials APICredentials) error {

	client := &http.Client{}

	var verifyResponse VerifyResponse

	urlStr := os.Getenv("LINE_API_BASE_URL") + "/oauth/verify"

	req, err := http.NewRequest("GET", urlStr, nil)

	req.Header.Add("Authorization", " Bearer "+credentials.AccessToken)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&verifyResponse)

	log.Println("Mid: " + verifyResponse.Mid)
	log.Println("ChannelId: ", verifyResponse.ChannelId)
	log.Println("Expire: ", verifyResponse.Expire)
	log.Println("Error: " + verifyResponse.Error)

	log.Println("Error Descrption: " + verifyResponse.ErrorDescription)

	// Return errors
	if verifyResponse.Error != "" {

		errorObject := APIError{
			ErrorCode:        verifyResponse.Error,
			ErrorDescription: verifyResponse.ErrorDescription,
		}

		return &errorObject
	}

	// If the API does not return an error, we can return a success
	return nil

}

func GetProfile() (ProfileInfo, error) {

	//Get Credentials from Redis
	credentials := GetCredentials()

	//Verify Credentials
	err := VerifyCredentials(credentials)

	if err != nil {

		switch err.(APIError).ErrorCode {
		case "412":
			log.Println("Token has expired. Attempting to refresh")
			err = RefreshAccessToken(credentials)

			if err != nil {
				log.Println("ERROR: Token Renewal has failed")
				return ProfileInfo{}, err
			}

		default:

			return ProfileInfo{}, err

		}

	}

	client := &http.Client{}

	var profileInfo ProfileInfo

	urlStr := os.Getenv("LINE_API_BASE_URL") + "/profile"

	req, err := http.NewRequest("GET", urlStr, nil)

	req.Header.Add("Authorization", " Bearer "+credentials.AccessToken)

	resp, err := client.Do(req)

	if err != nil {

		returnError := APIError{
			ErrorCode:        profileInfo.Error,
			ErrorDescription: profileInfo.ErrorDescription,
		}

		return ProfileInfo{}, &returnError

	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&profileInfo)

	return profileInfo, nil

}

func SendLinkMessage(mid string, username string) {

	//credentials := GetCredentials()

	log.Println("Entered SendMessage")

}
