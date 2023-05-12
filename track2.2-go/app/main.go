package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/findy-network/agency-workshop/agent"
	"github.com/findy-network/agency-workshop/handlers"
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	qrcode "github.com/skip2/go-qrcode"
)

type app struct {
	agencyClient *agent.AgencyClient
	// Store greeter handler to app state
	greeter *handlers.Greeter
	// Issuer handles the issuing logic
	issuer *handlers.Issuer
}

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
	// Create HTML payload
	_, html := try.To2(createInvitationPage(a.agencyClient.AgentClient, "Greet"))
	// Render HTML
	try.To1(response.Write([]byte(html)))
}

// Show pairwise invitation. Once connection is established, issue credential.
func (a *app) issueHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	id, html := try.To2(createInvitationPage(a.agencyClient.AgentClient, "Issue"))
	a.issuer.AddInvitation(id)
	try.To1(response.Write([]byte(html)))
}

// Show pairwise invitation. Once connection is established, verify credential.
func (a *app) verifyHandler(response http.ResponseWriter, r *http.Request) {
	defer err2.Catch(func(err error) {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	})
	try.To1(response.Write([]byte("IMPLEMENT ME")))
}

func createInvitationPage(
	agentClient agency.AgentServiceClient,
	header string,
) (html, invitationID string, err error) {
	defer err2.Handle(&err)

	// Agency API call for creating the DIDComm connection invitation
	res := try.To1(agentClient.CreateInvitation(
		context.TODO(),
		// Whichever name we want to expose from ourselves to the other end
		&agency.InvitationBase{Label: os.Getenv("FCLI_USER")},
	))

	var invitationMap map[string]any
	try.To(json.Unmarshal([]byte(res.GetJSON()), &invitationMap))

	url := res.URL
	log.Printf("Created invitation\n %s\n", url)

	// Convert invitation string to QR code
	png, err := qrcode.Encode(url, qrcode.Medium, 512)
	imgSrc := "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(png))

	// Create HTML payload
	html = `<html>
      <h1>` + header + `</h1>
      <p>Read the QR code with the wallet application:</p>
      <img src="` + imgSrc + `"/>
      <p>or copy-paste the invitation:</p>
      <textarea onclick="this.focus();this.select()" readonly="readonly" rows="10" cols="60">` +
		url + `</textarea></html>`

	// Return invitation id and the HTML payload
	return invitationMap["@id"].(string), html, nil
}

func main() {
	defer err2.Catch(func(err error) {
		log.Fatal(err)
	})

	// Login agent
	agencyClient := try.To1(agent.LoginAgent())

	// Create credential definition
	credDefId := try.To1(agencyClient.PrepareIssuing())

	// Create handlers
	myApp := app{
		agencyClient: agencyClient,
		greeter:      handlers.NewGreeter(agencyClient.Conn),
		issuer:       handlers.NewIssuer(agencyClient.Conn, credDefId),
	}

	// Start listening
	myApp.agencyClient.Listen([]agent.Listener{
		myApp.greeter,
		// Issuer handles the issuing logic
		myApp.issuer,
	})

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
