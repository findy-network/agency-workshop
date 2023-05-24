package handlers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/findy-network/agency-workshop/agent"
	"github.com/findy-network/findy-common-go/agency/client"
	"github.com/findy-network/findy-common-go/agency/client/async"
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Issuer struct {
	agent.DefaultListener
	conn        client.Conn
	connections sync.Map
	credDefID   string
}

type connection struct {
	id       string
	email    string
	verified bool
}

var (
	sgClient = sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
)

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
	i.connections.Store(id, &connection{id: id})
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

	if conn.email == "" {
		i.askForEmail(conn.id)
	}
}

func (i *Issuer) HandleBasicMesssageDone(
	notification *agency.Notification,
	status *agency.ProtocolStatus_BasicMessageStatus,
) {
	defer err2.Catch(func(err error) {
		log.Printf("Error handling basic message: %v", err)
	})

	conn := i.getConnection(notification.ConnectionID)

	// Skip handling if
	// 1. Connection was not for issuing
	// 2. Message was sent by us
	// 3. Email has been already asked
	if conn == nil || status.SentByMe || conn.email != "" {
		return
	}

	msg := status.Content
	msgValid := len(strings.Split(msg, " ")) == 1 && strings.Contains(msg, "@")

	log.Printf("Basic message %s with protocol id %s completed with %s",
		msg, notification.ProtocolID, conn.id)

	if msgValid {
		i.connections.Store(conn.id, &connection{id: conn.id, email: msg})

		// Create simple verification link
		// Note: in real-world we should use some random value instead of the connection id
		content := fmt.Sprintf("Please verify your email by clicking the following link:\n http://localhost:3001/email?value=%s", conn.id)
		i.sendEmail(content, msg)

		// Send confirmation via basic message
		pw := async.NewPairwise(i.conn, conn.id)
		_ = try.To1(pw.BasicMessage(context.TODO(), "Email is on it's way! Please check your mailbox 📫."))

	} else {
		// If email is invalid, ask again
		i.askForEmail(conn.id)
	}
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

func (i *Issuer) SetEmailVerified(connectionID string) (err error) {
	defer err2.Handle(&err)

	conn := i.getConnection(connectionID)

	// Skip handling if
	// 1. Connection was not for issuing
	// 2. Email has not been saved
	// 3. Credential has been already issued
	if conn == nil || conn.email == "" || conn.verified {
		return
	}

	i.connections.Store(conn.id, &connection{id: conn.id, email: conn.email, verified: true})

	// Create credential content
	attributes := make([]*agency.Protocol_IssuingAttributes_Attribute, 1)
	attributes[0] = &agency.Protocol_IssuingAttributes_Attribute{
		Name:  "email",
		Value: conn.email,
	}

	log.Printf(
		"Offer credential, conn id: %s, credDefID: %s, attrs: %v",
		conn.id,
		i.credDefID,
		attributes,
	)

	// Send credential offer to the other agent
	pw := async.NewPairwise(i.conn, conn.id)
	res := try.To1(pw.IssueWithAttrs(
		context.TODO(),
		i.credDefID,
		&agency.Protocol_IssuingAttributes{
			Attributes: attributes,
		}),
	)

	log.Printf("Credential offered: %s", res.GetID())
	return nil
}

func (i *Issuer) askForEmail(connectionID string) (err error) {
	defer err2.Handle(&err)

	// Ask for user email via basic message
	pw := async.NewPairwise(i.conn, connectionID)
	_ = try.To1(pw.BasicMessage(context.TODO(), "Please enter your email to get started."))

	return err
}

func (i *Issuer) sendEmail(content, email string) (err error) {
	defer err2.Handle(&err)

	from := mail.NewEmail("Issuer example", os.Getenv("SENDGRID_SENDER"))
	subject := "Email verification"
	to := mail.NewEmail(email, email) // Change to your recipient
	message := mail.NewSingleEmail(from, subject, to, content, content)

	log.Printf("Sending email %s to %s", content, email)
	_ = try.To1(sgClient.Send(message))

	return err
}
