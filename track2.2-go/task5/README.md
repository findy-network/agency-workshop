# Track 2.2 - Task 5: Verify credential

## Progress

* [Task 0: Setup environment](../README.md#task-0-setup-environment)
* [Task 1: Create a new connection](../task1/README.md#track-21---task-1-create-a-new-connection)
* [Task 2: Send greetings](../task2/README.md#track-21---task-2-send-greetings)
* [Task 3: Prepare for issuing credentials](../task3/README.md#track-21---task-3-prepare-for-issuing-credentials)
* [Task 4: Issue credential](../task4/README.md#track-21---task-4-issue-credential)
* **Task 5: Verify credential**
* [Task 6: Issue credential for verified information](../task6/README.md#track-21---task-6-issue-credential-for-verified-information)
* [Task 7: Additional tasks](../task7/README.md#track-21---task-7-additional-tasks)

## Description

Your web wallet user should now have their first credential in their wallet.
Now we can build the functionality that will verify that credential.

In a real-world implementation, we would naturally have two applications and two separate
agents, one for issuing and one for verifying. The wallet user would first acquire a credential
using the issuer application and then use the credential, i.e., prove the data,
in another application.

For simplicity, we build the verification functionality into the same application
we have been working on. The underlying protocol for requesting and presenting proofs is
[the present proof protocol](https://github.com/hyperledger/aries-rfcs/blob/main/features/0037-present-proof/README.md).

## 1. Listen to present proof protocol

Open file `agent/listen.go`.
Add new methods `HandlePresentProofPaused` and `HandlePresentProofDone` to listener interface:

```go
type Listener interface {
  HandleNewConnection(*agency.Notification, *agency.ProtocolStatus_DIDExchangeStatus)
  HandleBasicMesssageDone(*agency.Notification, *agency.ProtocolStatus_BasicMessageStatus)
  HandleIssueCredentialDone(*agency.Notification, *agency.ProtocolStatus_IssueCredentialStatus)
  // Send notification to listener when present proof protocol is paused
  HandlePresentProofPaused(*agency.Notification, *agency.ProtocolStatus_PresentProofStatus)
  // Send notification to listener when present proof protocol is completed
  HandlePresentProofDone(*agency.Notification, *agency.ProtocolStatus_PresentProofStatus)
}
```

When receiving notification for the issue credential protocol, notify listeners via the new method.
Edit `Listen`-function:

```go

  ...

func (agencyClient *AgencyClient) Listen(listeners []Listener) {

  ...

    // Notify listeners of protocol events
    switch notification.GetTypeID() {
    case agency.Notification_STATUS_UPDATE:
      if status.State.State == agency.ProtocolState_OK {
        switch notification.GetProtocolType() {
        case agency.Protocol_DIDEXCHANGE:
          for _, listener := range listeners {
            listener.HandleNewConnection(notification, status.GetDIDExchange())
          }
        case agency.Protocol_BASIC_MESSAGE:
          for _, listener := range listeners {
            listener.HandleBasicMesssageDone(notification, status.GetBasicMessage())
          }
        case agency.Protocol_ISSUE_CREDENTIAL:
          for _, listener := range listeners {
            listener.HandleIssueCredentialDone(notification, status.GetIssueCredential())
          }
          // Notify listener when present proof protocol is completed
        case agency.Protocol_PRESENT_PROOF:
          for _, listener := range listeners {
            listener.HandlePresentProofDone(notification, status.GetPresentProof())
          }
        default:
          log.Printf("No handler for protocol message %s\n", notification.GetProtocolType())
        }
      } else {
        log.Printf("Status NOK %v for %s\n", status, notification.GetProtocolType())
      }
      // Notify listener when present proof protocol is paused
    case agency.Notification_PROTOCOL_PAUSED:
      for _, listener := range listeners {
        listener.HandlePresentProofPaused(notification, status.GetPresentProof())
      }
    default:
      log.Printf("No handler for notification %s\n", notification.GetTypeID())
    }

  ...

}
```

## 1. Add code for verifying logic

Create a new file `src/verify.ts`.

Add the following content to the new file:

```go
package handlers

import (
  "context"
  "log"
  "sync"

  "github.com/findy-network/agency-workshop/agent"
  "github.com/findy-network/findy-common-go/agency/client"
  "github.com/findy-network/findy-common-go/agency/client/async"
  agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
  "github.com/lainio/err2"
  "github.com/lainio/err2/try"
)

type Verifier struct {
  agent.DefaultListener
  conn        client.Conn
  connections sync.Map
  credDefID   string
}

func NewVerifier(conn client.Conn, credDefID string) *Verifier {
  return &Verifier{
    conn:      conn,
    credDefID: credDefID,
  }
}

func (v *Verifier) getConnection(id string) *connection {
  if anyConn, ok := v.connections.Load(id); ok {
    if conn, ok := anyConn.(*connection); ok {
      return conn
    }
  }
  return nil
}

func (v *Verifier) AddInvitation(id string) {
  v.connections.Store(id, &connection{id})
}

func (v *Verifier) HandleNewConnection(
  notification *agency.Notification,
  status *agency.ProtocolStatus_DIDExchangeStatus,
) {
  defer err2.Catch(func(err error) {
    log.Printf("Error handling new connection: %v", err)
  })

  conn := v.getConnection(notification.ConnectionID)

  if conn == nil {
    // Connection was not for verifying, skip
    return
  }

  // Create proof request content
  attributes := make([]*agency.Protocol_Proof_Attribute, 1)
  attributes[0] = &agency.Protocol_Proof_Attribute{
    CredDefID: v.credDefID,
    Name:      "foo",
  }

  log.Printf("Request proof, conn id: %s, attrs: %v", notification.ConnectionID, attributes)

  // Send the proof request
  pw := async.NewPairwise(v.conn, notification.ConnectionID)
  res := try.To1(pw.ReqProofWithAttrs(context.TODO(), &agency.Protocol_Proof{
    Attributes: attributes,
  }))

  log.Printf("Proof request sent: %s", res.GetID())
}

// This function is called after proof is verified cryptographically.
// The application can execute its business logic and reject the proof
// if the attribute values are not valid.
func (v *Verifier) HandlePresentProofPaused(
  notification *agency.Notification,
  status *agency.ProtocolStatus_PresentProofStatus,
) {

  pw := async.NewPairwise(v.conn, notification.ConnectionID)

  // we have no special logic here - accept all received values
  res := try.To1(pw.Resume(
    context.TODO(),
    notification.ProtocolID,
    agency.Protocol_PRESENT_PROOF,
    agency.ProtocolState_ACK,
  ))

  log.Printf("Proof continued: %s", res.GetID())
}

func (v *Verifier) HandlePresentProofDone(
  notification *agency.Notification,
  status *agency.ProtocolStatus_PresentProofStatus,
) {
  conn := v.getConnection(notification.ConnectionID)

  if conn == nil {
    // Connection was not for issuing, skip
    return
  }

  log.Printf(
    "Proof verified from: %s, with id: %s",
    notification.ConnectionID,
    notification.ProtocolID,
  )

  v.connections.Delete(notification.ConnectionID)
}

```

## 3. Implement the `/verify`-endpoint

Open file `main.go`.

Add new field `verifier` to `app` state struct:

```go
type app struct {
  agencyClient *agent.AgencyClient
  greeter      *handlers.Greeter
  issuer       *handlers.Issuer
  // Verifier handles the verifying logic
  verifier     *handlers.Verifier
}
```

Modify function `main`.
Create the `verifier` and give it as a parameter on listener initialization:

```go
func main() {

  ...

  // Create handlers
  myApp := app{
    agencyClient: agencyClient,
    greeter:      handlers.NewGreeter(agencyClient.Conn),
    issuer:       handlers.NewIssuer(agencyClient.Conn, credDefId),
    // Handler for verifying logic
    verifier:     handlers.NewVerifier(agencyClient.Conn, credDefId),
  }

// Start listening
  myApp.agencyClient.Listen([]agent.Listener{
    myApp.greeter,
    myApp.issuer,
    // Add verifier to listener array
    myApp.verifier,
  })

  ...
}
```

Replace the implementation in the `/verify`-endpoint with the following:

```go
// Show pairwise invitation. Once connection is established, verify credential.
func (a *app) verifyHandler(response http.ResponseWriter, r *http.Request) {
  defer err2.Catch(func(err error) {
    log.Println(err)
    http.Error(response, err.Error(), http.StatusInternalServerError)
  })
  id, html := try.To2(createInvitationPage(a.agencyClient.AgentClient, "Verify"))
  a.verifier.AddInvitation(id)
  try.To1(response.Write([]byte(html)))
}
```

## 4. Test the `/verify`-endpoint

Make sure the server is restarted (`go run .`).
Open your browser to <http://localhost:3001/verify>

*You should see a simple web page with a QR code and a text input with a prefilled string.*

![Verify page](./docs/verify-page.png)

## 5. Read the QR code with the web wallet

Add the connection in the same way as in [task 1](../task1/README.md#6-read-the-qr-code-with-the-web-wallet):
Tap the "Add connection" button in your web wallet and read the QR code with your mobile device. Alternatively,
you can copy-paste the invitation string to the "Add connection"-dialog.

## 6. Ensure proof request is received in the web wallet

Accept proof request.

![Accept proof request](./docs/accept-proof-web-wallet.png)

## 7. Check server logs

Ensure that server logs display the success for the proof protocol:

![Server logs](./docs/server-logs-verify-proof.png)

## 8. Continue with task 6

Congratulations, you have completed task 6, and now know how to verify
credentials!

You can now continue with [task 6](../task6/README.md).
