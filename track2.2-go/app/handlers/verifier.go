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
