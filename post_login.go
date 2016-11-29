package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type PostLoginModel struct {
	DisplayName    string
	Mid            string
	AccessToken    string
	ProfilePicture string
	ProfileMessage string
	RefreshToken   string
	Expiry         int64
	LogoutLink     string
}

type GetAccessTokenResponse struct {
	Mid              string `json:"mid,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int64  `json:"expires_in,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	Scope            string `json:"scope,omitempty"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_desciption"` //This is a json format bug. Description is mispelleld
}

type ModalModel struct {
	TokenIsValid bool
}

func GetAccessToken(authorizationCode string) GetAccessTokenResponse {

	log.Println("Entered GetAccessToken")

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Add("client_id", os.Getenv("CHANNEL_ID"))
	data.Add("client_secret", os.Getenv("CHANNEL_SECRET"))
	data.Add("code", authorizationCode)
	data.Add("redirect_uri", os.Getenv("BASE_URL")+"postlogin")

	urlStr := os.Getenv("LINE_API_BASE_URL") + "/oauth/accessToken"

	log.Println("Parameters: " + data.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr+"?"+data.Encode(), nil)

	if err != nil {
		panic(err)
	}

	log.Println("Base urlStr: " + urlStr)

	log.Println("Authorize Token API Request URL: " + req.URL.String())

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)

	log.Println("Response Status:" + resp.Status)
	log.Println("Response URL" + resp.Request.RequestURI)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var accessTokenResponse GetAccessTokenResponse

	err = json.Unmarshal(body, &accessTokenResponse)

	if err != nil {
		panic(err)
	}

	return accessTokenResponse

}

// Change modal that appears when the user clicks the "Verify Token" modal
// depending on whether the token is valid or not

func VerifyToken(w http.ResponseWriter, r *http.Request) {

	model := ModalModel{TokenIsValid: false}

	credentials := GetCredentials()

	err := VerifyCredentials(credentials)

	if err == nil {

		model.TokenIsValid = true

	}

	log.Println("Token Model: ", model.TokenIsValid)

	t, _ := template.ParseFiles("templates/verify.html")
	t.Execute(w, &model)

}

func RefreshToken(w http.ResponseWriter, r *http.Request) {

	err := RefreshAccessToken(GetCredentials())

	if err != nil {
		panic(err)
	}

	//Get New Credentials

	credentials := GetCredentials()

	model := PostLoginModel{AccessToken: credentials.AccessToken, RefreshToken: credentials.RefreshToken, Expiry: credentials.Expire}

	t, _ := template.ParseFiles("templates/refresh.html")
	t.Execute(w, &model)

}

func Message(w http.ResponseWriter, r *http.Request) {

	// Get Profile info for user name
	profile, err := GetProfile()

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(r.Body)

	log.Println(string(body))

	SendLinkMessage("12312312", profile.DisplayName)

}

func Logout(w http.ResponseWriter, r *http.Request) {

	//Log out of LINE
	LineLogout()

	// Clear Credentials in Redis
	ClearCredentials()

	//fmt.Fprintf(w, "Logging out and redirecting to the login page...")

	http.Redirect(w, r, os.Getenv("BASE_URL")+"/login", 302)

}

func PostLogin(w http.ResponseWriter, r *http.Request) {

	//Get the Access Code from the http request

	log.Println("Entered Postlogin")

	requestUrl, err := url.Parse(r.RequestURI)

	if err != nil {
		fmt.Println(w, "Could not parse query")
	}

	queryParameters := requestUrl.Query()

	accessCode := queryParameters.Get("code")

	log.Println(accessCode)

	accessTokenResponse := GetAccessToken(accessCode)

	log.Printf("Mid: %s\n", accessTokenResponse.Mid)
	log.Printf("Access Token: %s\n", accessTokenResponse.AccessToken)
	log.Printf("Token Type: %s\n", accessTokenResponse.TokenType)
	log.Printf("Expires in: %s\n", accessTokenResponse.ExpiresIn)
	log.Printf("Refresh Token: %s\n", accessTokenResponse.RefreshToken)
	log.Printf("Scope: %s\n", accessTokenResponse.Scope)
	log.Printf("Error: %s\n", accessTokenResponse.Error)
	log.Printf("ErrorDescription: %s\n", accessTokenResponse.ErrorDescription)

	// Write the access token to redis
	WriteCredentials(APICredentials{AccessToken: accessTokenResponse.AccessToken, RefreshToken: accessTokenResponse.RefreshToken})

	// Get Profile
	profileInfo, err := GetProfile()

	if err != nil {
		panic(err)
	}

	postLoginModel := PostLoginModel{
		DisplayName:    profileInfo.DisplayName,
		Mid:            profileInfo.Mid,
		AccessToken:    accessTokenResponse.AccessToken,
		ProfilePicture: profileInfo.PictureUrl,
		ProfileMessage: profileInfo.StatusMessage,
		Expiry:         accessTokenResponse.ExpiresIn,
		RefreshToken:   accessTokenResponse.RefreshToken,
		LogoutLink:     os.Getenv("BASE_URL") + "/logout",
	}

	t, _ := template.ParseFiles("templates/postlogin.html")
	t.Execute(w, &postLoginModel)

}
