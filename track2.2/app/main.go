package main

import (
	"log"
	"net/http"
	"time"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

type appState struct{}

// Routes
func homeHandler(app appState) http.HandlerFunc {
	return func(response http.ResponseWriter, r *http.Request) {
		defer err2.Catch(func(err error) {
			log.Println(err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
		})
		try.To1(response.Write([]byte("Go example")))
	}
}

// Show pairwise invitation. Once connection is established, send greeting.
func greetHandler(app appState) http.HandlerFunc {
	return func(response http.ResponseWriter, r *http.Request) {
		defer err2.Catch(func(err error) {
			log.Println(err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
		})
		try.To1(response.Write([]byte("IMPLEMENT ME")))
	}
}

// Show pairwise invitation. Once connection is established, issue credential.
func issueHandler(app appState) http.HandlerFunc {
	return func(response http.ResponseWriter, r *http.Request) {
		defer err2.Catch(func(err error) {
			log.Println(err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
		})
		try.To1(response.Write([]byte("IMPLEMENT ME")))
	}
}

// Show pairwise invitation. Once connection is established, verify credential.
func verifyHandler(app appState) http.HandlerFunc {
	return func(response http.ResponseWriter, r *http.Request) {
		defer err2.Catch(func(err error) {
			log.Println(err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
		})
		try.To1(response.Write([]byte("IMPLEMENT ME")))
	}
}

func main() {
	defer err2.Catch(func(err error) {
		log.Fatal(err)
	})

	app := appState{}

	router := http.NewServeMux()

	router.HandleFunc("/", homeHandler(app))
	router.HandleFunc("/greet", greetHandler(app))
	router.HandleFunc("/issue", issueHandler(app))
	router.HandleFunc("/verify", verifyHandler(app))

	addr := ":3001"
	log.Printf("Starting server at %s", addr)

	server := http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	try.To(server.ListenAndServe())
}
