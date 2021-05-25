package main

import (
	"fmt"
	"net/http"
	"rangda"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	s, err := rangda.GetSecrets("./secrets.json")

	if err != nil {
		panic(err)
	}

	address := fmt.Sprintf("%s:%d", s.Host, s.Port)

	ran := rangda.New(s.ApiKey)

	r := mux.NewRouter()
	r.HandleFunc("/github/review", ran.ReviewEventHandler)

	srv := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
