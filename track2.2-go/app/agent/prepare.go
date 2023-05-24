package agent

import (
	"context"
	"log"
	"os"
	"time"

	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

func (a *AgencyClient) PrepareIssuing() (credDefID string, err error) {
	defer err2.Handle(&err)

	const credDefIDFileName = "CRED_DEF_ID"
	const schemaName = "email"
	schemaAttributes := []string{"email"}

	credDefIDBytes, credDefReadErr := os.ReadFile(credDefIDFileName)
	if credDefReadErr == nil {
		credDefID = string(credDefIDBytes)
		log.Printf("Credential definition %s exists already", credDefID)
		return
	}

	schemaRes := try.To1(a.AgentClient.CreateSchema(
		context.TODO(),
		&agency.SchemaCreate{
			Name:       schemaName,
			Version:    "1.0",
			Attributes: schemaAttributes,
		},
	))

	// wait for schema to be readable before creating cred def
	schemaGet := &agency.Schema{
		ID: schemaRes.ID,
	}
	schemaFound := false
	for !schemaFound {
		if _, err := a.AgentClient.GetSchema(context.TODO(), schemaGet); err == nil {
			schemaFound = true
		} else {
			time.Sleep(1 * time.Second)
		}
	}

	log.Printf("Schema %s created successfully", schemaRes.ID)

	res := try.To1(a.AgentClient.CreateCredDef(
		context.TODO(),
		&agency.CredDefCreate{
			SchemaID: schemaRes.ID,
			Tag:      os.Getenv("FCLI_USER"),
		},
	))
	credDefID = res.GetID()

	log.Printf("Credential definition %s created successfully", res.ID)
	try.To(os.WriteFile(credDefIDFileName, []byte(credDefID), 0666))

	return
}
