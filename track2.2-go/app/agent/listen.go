package agent

import (
	"context"
	"log"

	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	"github.com/google/uuid"
	"github.com/lainio/err2/try"
)

type Listener interface {
	HandleNewConnection(*agency.Notification, *agency.ProtocolStatus_DIDExchangeStatus)
}

func (agencyClient *AgencyClient) Listen(listeners []Listener) {

	// Listening callback handles agent events
	ch := try.To1(
		agencyClient.Conn.ListenStatus(
			context.TODO(),
			&agency.ClientID{ID: uuid.New().String()},
		),
	)

	// Go routine for listening event channel
	go func() {
		for {
			chRes, ok := <-ch
			if !ok {
				panic("Listening failed")
			}

			notification := chRes.GetNotification()
			log.Printf("Received agent notification %v\n", notification)

			// Fetch detailed status information for notification
			status := try.To1(
				agencyClient.ProtocolClient.Status(
					context.TODO(),
					&agency.ProtocolID{
						ID:     notification.ProtocolID,
						TypeID: notification.ProtocolType,
					},
				),
			)

			// Notify listeners of protocol events
			switch notification.GetTypeID() {
			case agency.Notification_STATUS_UPDATE:
				if status.State.State == agency.ProtocolState_OK {
					switch notification.GetProtocolType() {
					case agency.Protocol_DIDEXCHANGE:
						for _, listener := range listeners {
							listener.HandleNewConnection(notification, status.GetDIDExchange())
						}
					default:
						log.Printf("No handler for protocol message %s\n", notification.GetProtocolType())
					}
				} else {
					log.Printf("Status NOK %v for %s\n", status, notification.GetProtocolType())
				}
			default:
				log.Printf("No handler for notification %s\n", notification.GetTypeID())
			}

		}
	}()

}
