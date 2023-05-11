package handlers

import (
	"log"

	"github.com/findy-network/agency-workshop/agent"
	"github.com/findy-network/findy-common-go/agency/client"
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
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
	// Print their label to logs
	log.Printf("New connection %s with id %s", status.TheirLabel, notification.ConnectionID)
}
