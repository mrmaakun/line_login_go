package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("./web")))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
