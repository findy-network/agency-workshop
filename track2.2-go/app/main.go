package main

import (
	"log"
	"net/http"
	"time"

	"github.com/findy-network/agency-workshop/agent"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

type app struct{}

func (a *app) homeHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	try.To1(response.Write([]byte("Go example")))
}

// Show pairwise invitation. Once connection is established, send greeting.
func (a *app) greetHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	try.To1(response.Write([]byte("IMPLEMENT ME")))
}

// Show pairwise invitation. Once connection is established, issue credential.
func (a *app) issueHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	try.To1(response.Write([]byte("IMPLEMENT ME")))
}

// Show pairwise invitation. Once connection is established, verify credential.
func (a *app) verifyHandler(response http.ResponseWriter, r *http.Request) {
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

	// Login agent
	_ = try.To1(agent.LoginAgent())

	myApp := app{}

	router := http.NewServeMux()

	router.HandleFunc("/", myApp.homeHandler)
	router.HandleFunc("/greet", myApp.greetHandler)
	router.HandleFunc("/issue", myApp.issueHandler)
	router.HandleFunc("/verify", myApp.verifyHandler)

	addr := ":3001"
	log.Printf("Starting server at %s", addr)

	server := http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	try.To(server.ListenAndServe())
}
