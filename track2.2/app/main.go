package main

import (
	"log"
	"net/http"
	"time"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

// Routes
func homeHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	try.To1(response.Write([]byte("Go example")))
}

// Show pairwise invitation. Once connection is established, send greeting.
func greetHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	try.To1(response.Write([]byte("IMPLEMENT ME")))
}

// Show pairwise invitation. Once connection is established, issue credential.
func issueHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	try.To1(response.Write([]byte("IMPLEMENT ME")))
}

// Show pairwise invitation. Once connection is established, verify credential.
func verifyHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	try.To1(response.Write([]byte("IMPLEMENT ME")))
}

func main() {
	defer err2.Catch(func(err error) {
		log.Fatal(err)
	})

	router := http.NewServeMux()

	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/greet", greetHandler)
	router.HandleFunc("/issue", issueHandler)
	router.HandleFunc("/verify", verifyHandler)

	addr := ":3001"
	log.Printf("Starting server at %s", addr)

	server := http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	try.To(server.ListenAndServe())
}
