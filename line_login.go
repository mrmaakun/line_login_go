package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
)

type LoginModel struct {
	LoginUrl string
}

func LoginPage(w http.ResponseWriter, r *http.Request) {

	//Get the Access Code from the http request

	log.Println("Entered LoginPage")

	redirectUrl := os.Getenv("BASE_URL") + "/postlogin"

	escapedRedirectUrl := url.QueryEscape(redirectUrl)

	loginUrl := "https://access.line.me/dialog/oauth/weblogin?response_type=code&client_id=" + os.Getenv("CHANNEL_ID") + "&redirect_uri=" + escapedRedirectUrl

	loginModel := LoginModel{LoginUrl: loginUrl}

	log.Println("Login URL: " + loginModel.LoginUrl)

	t, _ := template.ParseFiles("templates/login.html")
	t.Execute(w, &loginModel)

}

func main() {

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/postlogin", PostLogin)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/verify", VerifyToken)
	http.HandleFunc("/refresh", RefreshToken)
	http.HandleFunc("/message", Message)

	http.Handle("/", http.FileServer(http.Dir("./web")))

	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
