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

	r := mux.NewRouter()
	r.HandleFunc("/", helloWorld)

	fmt.Printf("Rangda listening on %s\n", url)
	log.Fatal(http.ListenAndServe(url, r))
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
