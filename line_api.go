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
}

type RefreshResponse struct {
	Mid              string `json:"mid,omitempty"`
	AccessToken      string `json:"accessToken,omitempty"`
	RefreshToken     string `json:"refreshToken,omitempty"`
	Expire           int64  `json:"expire,omitempty"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_descrition"`
}

type VerifyResponse struct {
	Mid              string `json:"mid,omitempty"`
	ChannelId        int64  `json:"channelId,omitempty"`
	Expire           int64  `json:"expire,omitempty"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_descrition"`
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

	req, err := http.NewRequest("POST", urlStr, nil)

	req.Header.Add("Content-Type", "x-www-form-urlencoded")
	req.Header.Add("X-Line-ChannelToken", credentials.AccessToken)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&refreshResponse)

	newCredentials := APICredentials{}

	if refreshResponse.Error != "" {

		newCredentials.AccessToken = refreshResponse.AccessToken
		newCredentials.RefreshToken = refreshResponse.RefreshToken

		return nil

	} else {

		returnError := APIError{
			ErrorCode:        refreshResponse.Error,
			ErrorDescription: refreshResponse.ErrorDescription,
		}

		return returnError
	}

	// If there are no errors, refresh the credentials in Redis

	WriteCredentials(newCredentials)

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
