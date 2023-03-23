import express, { Express, Request, Response } from 'express';
import { createAcator, openGRPCConnection, agencyv1 } from '@findy-network/findy-common-ts'
import QRCode from 'qrcode'

import prepareIssuer from './prepare';
import { existsSync, readFileSync, writeFileSync } from 'fs';
import mailer from '@sendgrid/mail'


const app: Express = express();
const port = process.env.PORT || 3001;
const invitations: {
  issue: string[],
  verify: string[],
  waitingForEmail: string[],
  waitingForVerification: Map<string, string>
} = {
  issue: [],
  verify: [],
  waitingForEmail: [],
  waitingForVerification: new Map<string, string>()
}
const userName = process.env.FCLI_USER || 'ts-example'

mailer.setApiKey(process.env.SENDGRID_API_KEY || '')

const setupFindyAgency = async () => {
  const acatorProps = {
    authUrl: process.env.FCLI_URL || 'http://localhost:8088',
    authOrigin: process.env.FCLI_ORIGIN || 'http://localhost:3000',
    userName,
    key: process.env.FCLI_KEY || '15308490f1e4026284594dd08d31291bc8ef2aeac730d0daf6ff87bb92d4336c',
  }
  const authenticator = createAcator(acatorProps)

  const serverAddress = process.env.AGENCY_API_SERVER || 'localhost'
  const certPath = process.env.FCLI_TLS_PATH || ''
  const grpcProps = {
    serverAddress,
    serverPort: parseInt(process.env.AGENCY_API_SERVER_PORT || '50052', 10),
    // NOTE: we currently assume that we do not need certs for cloud installation
    // as the cert is issued by a trusted issuer
    certPath: serverAddress === 'localhost' ? certPath : ''
  }

  // Authenticate and open GRPC connection to agency
  return openGRPCConnection(grpcProps, authenticator)
}

const runApp = async () => {
  const agencyConnection = await setupFindyAgency()
  const { createAgentClient, createProtocolClient } = agencyConnection
  const agentClient = await createAgentClient()
  const protocolClient = await createProtocolClient()
  // Credential definition is created on server startup.
  // We need it to be able to issue credentials.
  const credDefIdFilePath = "CRED_DEF_ID"
  const credDefCreated = existsSync(credDefIdFilePath)
  const credDefId = credDefCreated ? readFileSync(credDefIdFilePath).toString() : await prepareIssuer(agentClient, userName)
  if (!credDefCreated) {
    // store id in order to avoid unnecessary creations
    writeFileSync(credDefIdFilePath, credDefId)
  }
  console.log(`Credential definition available: ${credDefId}`)

  const askForEmail = async (connectionId: string) => {
    const msg = new agencyv1.Protocol.BasicMessageMsg()
    msg.setContent("Hi there 👋!\n Please enter your email to get started.")
    await protocolClient.sendBasicMessage(connectionId, msg)
  }

  // Listening callback handles agent events
  await agentClient.startListeningWithHandler(
    {
      // New connection is established
      DIDExchangeDone: async (info) => {
        console.log(`New connection: ${info.connectionId}`)

        // If connection was for issuing, continue by querying user's email
        if (invitations.issue.includes(info.connectionId)) {
          await askForEmail(info.connectionId)
          invitations.waitingForEmail = [info.connectionId, ...invitations.waitingForEmail]

        } else {
          const attributes = new agencyv1.Protocol.Proof()
          const attr = new agencyv1.Protocol.Proof.Attribute()
          attr.setName("foo")
          attr.setCredDefid(credDefId)
          attributes.addAttributes(attr)

          const proofRequest = new agencyv1.Protocol.PresentProofMsg()
          proofRequest.setAttributes(attributes)

          await protocolClient.sendProofRequest(info.connectionId, proofRequest)
        }
      },
      BasicMessageDone: async (info, basicMessage) => {
        if (!basicMessage.getSentByMe() && invitations.waitingForEmail.includes(info.connectionId)) {
          const msg = basicMessage.getContent()
          // dummy validation
          if (msg.split(' ').length === 1 && msg.indexOf('@') >= 0) {
            invitations.waitingForEmail = invitations.waitingForEmail.filter(item => item !== info.connectionId)
            // send verification mail
            const content = `Please verify your email by clicking the following link:\n http://localhost:3001/email/${info.connectionId}`
            const emailMsg = {
              to: msg,
              from: {
                email: process.env.SENDGRID_SENDER || '',
                name: "Issuer example",
              },
              subject: 'Email verification',
              text: content,
              html: content
            }
            console.log("Sending:", emailMsg)
            const ok = await mailer.send(emailMsg)
            invitations.waitingForVerification.set(info.connectionId, msg)
          } else {
            await askForEmail(info.connectionId)
          }
        }
      },
      IssueCredentialDone: (info) => {
        console.log(`Credential issued: ${info.protocolId}`)
        invitations.issue = invitations.issue.filter(item => item !== info.connectionId)
      },

      // This function is called after proof is verified cryptographically.
      // The application can execute its business logic and reject the proof
      // if the attribute values are not valid.
      PresentProofPaused: async (info, presentProof) => {
        console.log(`Proof paused: ${info.protocolId}`)
        presentProof.getProof()?.getAttributesList().forEach((value, index) => {
          console.log(`Proof attribute ${index} ${value.getName()}: ${value.getValue()}`)
        })
        const protocolID = new agencyv1.ProtocolID()
        protocolID.setId(info.protocolId)
        protocolID.setTypeid(agencyv1.Protocol.Type.PRESENT_PROOF)
        protocolID.setRole(agencyv1.Protocol.Role.RESUMER)
        const msg = new agencyv1.ProtocolState()
        msg.setProtocolid(protocolID)
        // we have no special logic here - accept all received values
        msg.setState(agencyv1.ProtocolState.State.ACK)
        await protocolClient.resume(msg)
      },
      PresentProofDone: (info) => {
        console.log(`Proof verified: ${info.protocolId}`)
        invitations.verify = invitations.verify.filter(item => item !== info.connectionId)
      },
    },
    {
      protocolClient,
      retryOnError: true,
    }
  )

  const renderInvitation = async (header: string, res: Response) => {
    const msg = new agencyv1.InvitationBase()
    msg.setLabel(userName)

    const invitation = await agentClient.createInvitation(msg)

    console.log(`Created invitation with Findy Agency: ${invitation.getUrl()}`)
    const qrData = await QRCode.toDataURL(invitation.getUrl())

    res.send(`<html>
    <h1>${header}</h1>
    <p>Read the QR code with the wallet application:</p>
    <img src="${qrData}"/>
    <p>or copy-paste the invitation:</p>
    <textarea onclick="this.focus();this.select()" readonly="readonly" rows="10" cols="60">${invitation.getUrl()}</textarea>
</html>`);

    return JSON.parse(invitation.getJson())["@id"]
  }

  // Show pairwise invitation. Once connection is established, verify credential.
  app.get('/verify', async (req: Request, res: Response) => {
    const id = await renderInvitation("Verify proof", res)
    invitations.verify = [...invitations.verify, id]
  });

  // Show pairwise invitation. Once connection is established, issue credential.
  app.get('/issue', async (req: Request, res: Response) => {
    const id = await renderInvitation("Issue credential", res)
    invitations.issue = [...invitations.issue, id]
  });

  // Verify email
  app.get('/email/:connectionId', async (req: Request, res: Response) => {
    const { connectionId } = req.params
    const item = invitations.waitingForVerification.get(connectionId)
    if (item) {
      invitations.waitingForVerification.delete(connectionId)
      const attributes = new agencyv1.Protocol.IssuingAttributes()
      const attr = new agencyv1.Protocol.IssuingAttributes.Attribute()
      attr.setName("email")
      attr.setValue(item || " ")
      attributes.addAttributes(attr)

      const credential = new agencyv1.Protocol.IssueCredentialMsg()
      credential.setCredDefid(credDefId)
      credential.setAttributes(attributes)

      await protocolClient.sendCredentialOffer(connectionId, credential)
      res.send(`<html>
      <h1>Offer sent!</h1>
      <p>Please open your wallet application and accept the credential.</p>
      <p>You can close this window.</p>
  </html>`);

    } else {
      res.send(`<html><h1>Error</h1></html>`);
    }
  });

  app.get('/', (req: Request, res: Response) => {
    res.send('Typescript example');
  });

  app.listen(port, async () => {
    console.log(`⚡️[server]: Server is running at http://localhost:${port}`);
  });
}

runApp()