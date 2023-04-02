# Track 2.1 - Task 6: Issue credential for verified information

## Progress

* [Task 0: Setup environment](../README.md#task-0-setup-environment)
* [Task 1: Create a new connection](../task1/README.md#track-21---task-1-create-a-new-connection)
* [Task 2: Send greetings](../task2/README.md#track-21---task-2-send-greetings)
* [Task 3: Prepare for issuing credentials](../task3/README.md#track-21---task-3-prepare-for-issuing-credentials)
* [Task 4: Issue credential](../task4/README.md#track-21---task-4-issue-credential)
* [Task 5: Verify credential](../task5/README.md#track-21---task-5-verify-credential)
* **Task 6: Issue credential for verified information**
* [Task 7: Additional tasks](../task7/README.md#track-21---task-7-additional-tasks)

## Description

In our previous issuing example, we issued a `foobar` credential to anyone who connects with us.
However, this is not a likely real-world scenario. Probably the issuer wishes to issue some
meaningful data that it knows to be valid.

Let's change our issuer so that it issues credentials for a verified email. The issuer displays
a QR code as before, but when the connection is established, it will ask for the user's email address.
It sends an email with a verification URL to this address. Only when a user opens the verification URL,
issuer will send the credential offer.

We need to create new schema and credential definition to issue email credentials. We also
need to create logic for asking the email address. In addition, a new endpoint needs to be added for
the verification URL.

In this task, we will utilize SendGrid API for sending the emails. You need an API key
to access the SendGrid API. You will be provided one in the guided workshop.

<details>
<summary>🤠 Acquire SendGrid API key</summary></br>

Create a free account to SendGrid and acquire the API key:
<https://sendgrid.com/solutions/email-api/>

Configure and verify also a sender identity for your email address.

</details></br>

## 1. Install SendGrid dependency

Add a new dependency to your project:

```bash
npm install @sendgrid/mail --save
```

## 2. Export environment variables for SendGrid API access

Open file `.envrc`. Add two new environment variables there:

```bash
export SENDGRID_API_KEY='<this_value_will_be_provided_for_you_in_the_workshop>'
export SENDGRID_SENDER='<this_value_will_be_provided_for_you_in_the_workshop>'
```

Save the file and type `direnv allow`. Restart your server `npm run dev`.

<details>
<summary>🤠 Configure own SendGrid account</summary></br>

Create API key with SendGrid UI and replace the value to `SENDGRID_API_KEY` variable.
Configure the verified sender email to `SENDGRID_SENDER` variable.

</details></br>

## 3. Create new credential definition

In [task 3](../task3/README.md) we created a schema and credential definition for `foobar`-credentials.
Now we need another schema and credential definition for our email credential.

Let's modify our code for creating the schema and the credential definition.

Open file `src/prepare.ts`.

Modify function `prepareIssuing`. Change the `schemaName` to `'email'` and attributes list to `['email']`:

```ts
  const prepareIssuing = async (): Promise<string> => {
    // Schema name
    const schemaName = 'email'
    console.log(`Creating schema ${schemaName}`)

    const schemaMsg = new agencyv1.SchemaCreate()
    schemaMsg.setName(schemaName)
    schemaMsg.setVersion('1.0')
    // Schema attribute list
    schemaMsg.setAttributesList(['email'])

    const schemaId = (await agentClient.createSchema(schemaMsg)).getId()
    return await createCredDef(schemaId)
  }
```

Then delete (or rename) file `CRED_DEF_ID` from the workspace root.
This ensures that the schema and credential definition creation code is executed on server startup
as there is no cached credential definition id.

```bash
mv CRED_DEF_ID foobar_CRED_DEF_ID
```

## 4. Ensure credential definition for email schema is created

<<screencapture>>

## 5. Modify issuer for email changes

Open file `src/issuer.ts`.

Add following row to imports:

```ts
import mailer from '@sendgrid/mail'
```

Add new functions `handleBasicMessageDone` and `setEmailVerified` to `Issuer`-interface:

```ts
export interface Issuer {
  addInvitation: (id: string) => void
  handleNewConnection: (info: ProtocolInfo, didExchange: ProtocolStatus.DIDExchangeStatus) => Promise<void>
  handleBasicMessageDone: (info: ProtocolInfo, basicMessage: ProtocolStatus.BasicMessageStatus) => Promise<void>
  handleIssueDone: (info: ProtocolInfo, issueCredential: ProtocolStatus.IssueCredentialStatus) => void
  setEmailVerified: (connectionId: string) => Promise<boolean>
}
```

Add new fields to `email` and `verified` to `Connection`-interface:

```ts
interface Connection {
  id: string
  email?: string
  verified?: boolean
}
```

Configure SendGrid API and add new utility functions `askForEmail` and `sendEmail`
for sending messages to default exported function:

```ts
export default (protocolClient: ProtocolClient, credDefId: string) => {
  const connections: Connection[] = []

  // Configure API key for SendGrid API
  mailer.setApiKey(process.env.SENDGRID_API_KEY!)

  // Ask for user email via basic message
  const askForEmail = async (connectionId: string) => {
    const msg = new agencyv1.Protocol.BasicMessageMsg()
    msg.setContent('Please enter your email to get started.')
    await protocolClient.sendBasicMessage(connectionId, msg)
  }

  // Send email via SendGrid API
  const sendEmail = async (email: string, content: string) => {
    const emailMsg = {
      to: email,
      from: {
        email: process.env.SENDGRID_SENDER!,
        name: 'Issuer example',
      },
      subject: 'Email verification',
      text: content,
      html: content
    }
    console.log(`Sending email '${content}' to ${email}`)
    await mailer.send(emailMsg)
  }

...

}
```

Instead of issuing the credential when new connection is established,
we want to ask user for their email, and send the credential offer
only when they have verified their email. So when a new connection is established,
we send a basic message asking the email.
Replace contents of `handleNewConnection` to following:

```ts
  const handleNewConnection = async (info: ProtocolInfo, didExchange: ProtocolStatus.DIDExchangeStatus) => {
    // Skip if this connection was not for issuing
    const connection = connections.find(({ id }) => id === info.connectionId)
    if (!connection) {
      return
    }

    // Ask for email from the other end
    if (!connection.email) {
      await askForEmail(info.connectionId)
    }
  }
```

Add new function `handleBasicMessageDone`. This function will handle the basic messages user
is sending from the other end. If user replies with an email address, a verification email is sent
to the provided address.

```ts
  const handleBasicMessageDone = async (info: ProtocolInfo, basicMessage: ProtocolStatus.BasicMessageStatus) => {
    const connection = connections.find(({ id }) => id === info.connectionId)
    // Skip handling if message was sent by us or
    // the verification is already done
    if (basicMessage.getSentByMe() || !connection || connection.email) {
      return
    }
    console.log(`Basic message\n${JSON.stringify(basicMessage.toObject())}\nwith protocol id ${info.protocolId} completed with ${info.connectionId}`)

    // Some sanity checking
    const email = basicMessage.getContent()
    const emailValid = email.split(' ').length === 1 && email.indexOf('@') >= 0

    if (emailValid) {
      // Valid email, do verification
      connection.email = email
      // Create simple verification link
      // Note: in real-world we should use some random value instead of the connection id
      const content = `Please verify your email by clicking the following link:\n http://localhost:3001/email/${connection.id}`
      // Send verification mail
      await sendEmail(connection.email, content)

      // Send confirmation via basic message
      const msg = new agencyv1.Protocol.BasicMessageMsg()
      msg.setContent(`Email is on it's way! Please check your mailbox 📫.`)
      await protocolClient.sendBasicMessage(connection.id, msg)

    } else {
      // Email invalid, ask again
      await askForEmail(info.connectionId)
    }
  }
```

Add new function `setEmailVerified`. This function will send a credential offer
of a verified email when the user has clicked the verification link.

```ts
  const setEmailVerified = async (connectionId: string) => {
    const connection = connections.find(({ id }) => id === connectionId)

    if (!connection || !connection.email || connection.verified) {
      return false
    }

    connection.verified = true

    // Send credential offer for verified email
    const attributes = new agencyv1.Protocol.IssuingAttributes()
    const attr = new agencyv1.Protocol.IssuingAttributes.Attribute()
    attr.setName('email')
    attr.setValue(connection.email)
    attributes.addAttributes(attr)

    const credential = new agencyv1.Protocol.IssueCredentialMsg()
    credential.setCredDefid(credDefId)
    credential.setAttributes(attributes)

    console.log(`Sending credential offer\n${JSON.stringify(credential.toObject())}\nto ${connectionId}`)
    await protocolClient.sendCredentialOffer(connectionId, credential)

    return true
  }
```

And finally, update the list of returned functions with the newly created functions
`handleBasicMessageDone` and `setEmailVerified`:

```ts
...

export default (protocolClient: ProtocolClient, credDefId: string) => {

  ...

  return {
    addInvitation,
    handleNewConnection,
    handleBasicMessageDone,
    handleIssueDone,
    setEmailVerified
  }
}
```

## 5. Send new notifications to issuer from listener

Open file `src/listen.ts`.

We need to notify also issuer for the received basic messages so that
the issuer can handle user provider email address.

Replace the contents for `BasicMessageDone` to following:

```ts
      ...

      BasicMessageDone: async (info, basicMessage) => {
        // Print out greeting sent from the other agent
        if (!basicMessage.getSentByMe()) {
          const msg = basicMessage.getContent()
          console.log(`Received basic message ${msg} from ${info.connectionId}`)
        }

        // Notify issuer
        issuer.handleBasicMessageDone(info, basicMessage)
      },

      ...
```

## 6. Add endpoint for email verification

Open file `src/index.ts`.

Add a new endpoint that handles email URL clicks.
The function asks issuer to send a credential offer if the connection is valid and found.

```ts
  app.get('/email/:connectionId', async (req: Request, res: Response) => {
    const { connectionId } = req.params
    // Ask issuer to send credential offer for verified email
    if (await issuer.setEmailVerified(connectionId)) {
      res.send(`<html>
    <h1>Offer sent!</h1>
    <p>Please open your wallet application and accept the credential.</p>
    <p>You can close this window.</p></html>`);
    } else {
      res.send(`<html><h1>Error</h1></html>`);
    }
  });
```

## 7. Testing

<<screencapture here>>

## 8. Continue with task 7

Congratulations, you have completed task 7 and you know now a little bit more
how to build the application logic for issuers!

You can now continue with [task 7](../task7/README.md).
