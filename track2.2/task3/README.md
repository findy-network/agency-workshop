# Track 2.2 - Task 3: Prepare for issuing credentials

## Progress

* [Task 0: Setup environment](../README.md#task-0-setup-environment)
* [Task 1: Create a new connection](../task1/README.md#track-21---task-1-create-a-new-connection)
* [Task 2: Send greetings](../task2/README.md#track-21---task-2-send-greetings)
* **Task 3: Prepare for issuing credentials**
* [Task 4: Issue credential](../task4/README.md#track-21---task-4-issue-credential)
* [Task 5: Verify credential](../task5/README.md#track-21---task-5-verify-credential)
* [Task 6: Issue credential for verified information](../task6/README.md#track-21---task-6-issue-credential-for-verified-information)
* [Task 7: Additional tasks](../task7/README.md#track-21---task-7-additional-tasks)

## Description

In the previous task, we learned how to start Hyperledger Aries protocol interactions and
react to the protocol notifications utilizing the Findy Agency agent API. In the following
tasks, we will learn how to issue and verify credentials using similar APIs.

But before issuing credentials, we have to prepare our agent for it.
It means that we must have a suitable schema and credential definition available.

A schema describes the contents of the verifiable credential: which data attributes it
contains. The credential definition is like a public key published against that schema.
Other parties can verify the credential's validity against the credential definition and
ensure that your and only your agent has issued the credential.

## 1. Add code for creating credential definition

The creation of the credential definition is only needed then when we start
to issue new types of credentials. So we don't need to do it too often.

Create a new file `agent/prepare.go`.

Add the following content to the new file:

```go
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
  const schemaName = "foobar"
  schemaAttributes := []string{"foo"}

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
```

## 2. Create credential definition on server start

Open file `main.go`.

Modify the `main`-function to create the credential definition on server start.

```go
func main() {

  ...

  // Login agent
  agencyClient := try.To1(agent.LoginAgent())

  // Create credential definition
  _ = try.To1(agencyClient.PrepareIssuing())

  ...

}
```

## 3. Ensure the credential definition is created from logs

Note! It will take a while for the agency to create a new credential definition.
Wait patiently.

![Server logs](./docs/server-logs-cred-def.png)

## 4. Check `CRED_DEF_ID`-file

When the credential definition is created, the logic in your server code will store
the credential definition id in a text file  `CRED_DEF_ID` in the workspace root. If this file exists,
the creation of the credential definition is skipped on server start.

![CRED_DEF_ID file](./docs/cred-def-file.png)

By default, the logic will create a schema `foobar` with one attribute `foo`. If you wish to change either
the schema name, the schema attributes or the tag of the credential definition, you need to delete
the `CRED_DEF_ID` file from the workspace root (so that the credential definition gets recreated).

For now, you can use the default values and proceed to the next tasks with the current credential
definition.

## 5. Continue with task 4

Congratulations, you have completed task 3, and you now know how to prepare your agent
to issue credentials!

You can now continue with [task 4](../task4/README.md).
