package handlers

import (
	"context"
	"log"

	"github.com/findy-network/agency-workshop/agent"
	"github.com/findy-network/findy-common-go/agency/client"
	"github.com/findy-network/findy-common-go/agency/client/async"
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

type Greeter struct {
	agent.DefaultListener
	conn client.Conn
}

func NewGreeter(conn client.Conn) *Greeter {
	return &Greeter{conn: conn}
}

func (g *Greeter) HandleNewConnection(
	notification *agency.Notification,
	status *agency.ProtocolStatus_DIDExchangeStatus,
) {
	defer err2.Catch(func(err error) {
		log.Printf("Error handling new connection: %v", err)
	})

	log.Printf("New connection %s with id %s", status.TheirLabel, notification.ConnectionID)

	// Greet each new connection with basic message
	pw := async.NewPairwise(g.conn, notification.ConnectionID)
	_ = try.To1(pw.BasicMessage(context.TODO(), "Hi there 👋!"))
}

func (g *Greeter) HandleBasicMesssageDone(
	notification *agency.Notification,
	status *agency.ProtocolStatus_BasicMessageStatus,
) {
	// Print out greeting sent from the other agent
	if !status.SentByMe {
		log.Printf("Received basic message %s from %s", status.Content, notification.ConnectionID)
	}
}
