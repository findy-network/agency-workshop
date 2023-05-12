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

type Issuer struct {
	agent.DefaultListener
	conn        client.Conn
	connections sync.Map
	credDefID   string
}

type connection struct {
	id string
}

func NewIssuer(conn client.Conn, credDefID string) *Issuer {
	return &Issuer{
		conn:      conn,
		credDefID: credDefID,
	}
}

func (i *Issuer) getConnection(id string) *connection {
	if anyConn, ok := i.connections.Load(id); ok {
		if conn, ok := anyConn.(*connection); ok {
			return conn
		}
	}
	return nil
}

func (i *Issuer) AddInvitation(id string) {
	i.connections.Store(id, &connection{id})
}

func (i *Issuer) HandleNewConnection(
	notification *agency.Notification,
	status *agency.ProtocolStatus_DIDExchangeStatus,
) {
	defer err2.Catch(func(err error) {
		log.Printf("Error handling new connection: %v", err)
	})

	conn := i.getConnection(notification.ConnectionID)

	if conn == nil {
		// Connection was not for issuing, skip
		return
	}

	// Create credential content
	attributes := make([]*agency.Protocol_IssuingAttributes_Attribute, 1)
	attributes[0] = &agency.Protocol_IssuingAttributes_Attribute{
		Name:  "foo",
		Value: "bar",
	}

	log.Printf(
		"Offer credential, conn id: %s, credDefID: %s, attrs: %v",
		notification.ConnectionID,
		i.credDefID,
		attributes,
	)

	// Send credential offer to the other agent
	pw := async.NewPairwise(i.conn, notification.ConnectionID)
	res := try.To1(pw.IssueWithAttrs(
		context.TODO(),
		i.credDefID,
		&agency.Protocol_IssuingAttributes{
			Attributes: attributes,
		}),
	)

	log.Printf("Credential offered: %s", res.GetID())
}

func (i *Issuer) HandleIssueCredentialDone(
	notification *agency.Notification,
	status *agency.ProtocolStatus_IssueCredentialStatus,
) {
	conn := i.getConnection(notification.ConnectionID)

	if conn == nil {
		// Connection was not for issuing, skip
		return
	}

	log.Printf(
		"Credential issued to: %s, with id: %s",
		notification.ConnectionID,
		notification.ProtocolID,
	)

	i.connections.Delete(notification.ConnectionID)
}
