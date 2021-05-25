package main

import (
	"fmt"
	"log"
	"net/http"
	"rangda"

	"github.com/gorilla/mux"
)

func main() {
	s, err := rangda.GetSecrets("./secrets.json")

	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("%s:%d", s.Host, s.Port)

	ran := rangda.New(s.ApiKey)

	r := mux.NewRouter()
	r.HandleFunc("/github/review", ran.ReviewEventHandler)

	fmt.Printf("Rangda listening on %s\n", url)
	log.Fatal(http.ListenAndServe(url, r))
}
