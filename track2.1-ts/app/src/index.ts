import express, { Express, Request, Response } from 'express';
import { agencyv1, AgentClient, createAcator, openGRPCConnection } from '@findy-network/findy-common-ts'
import QRCode from 'qrcode'

import listenAgent from './listen'

const app: Express = express();
const port = process.env.PORT || 3001;

const createInvitationPage = async (agentClient: AgentClient, header: string) => {
  // Invitation input
  const msg = new agencyv1.InvitationBase()
  // Whichever name we want to expose from ourselves to the other end
  msg.setLabel(process.env.FCLI_USER!)

  // Agency API call for creating the DIDComm connection invitation
  const invitation = await agentClient.createInvitation(msg)

  console.log(`Created invitation with Findy Agency: ${invitation.getUrl()}`)
  // Convert invitation string to QR code
  const qrData = await QRCode.toDataURL(invitation.getUrl())

  // Create HTML payload
  const payload = `<html>
  <h1>${header}</h1>
  <p>Read the QR code with the wallet application:</p>
  <img src="${qrData}"/>
  <p>or copy-paste the invitation:</p>
  <textarea onclick="this.focus();this.select()" readonly="readonly" rows="10" cols="60">${invitation.getUrl()
    }</textarea></html>`;

  // Return invitation id and the HTML payload
  return { id: JSON.parse(invitation.getJson())['@id'], payload }
}

const setupAgentConnection = async () => {
  const acatorProps = {
    authUrl: process.env.FCLI_URL!,
    authOrigin: process.env.FCLI_ORIGIN!,
    userName: process.env.FCLI_USER!,
    key: process.env.FCLI_KEY!,
  }
  // Create authenticator
  const authenticator = createAcator(acatorProps)

  const serverAddress = process.env.AGENCY_API_SERVER!
  const certPath = process.env.FCLI_TLS_PATH!
  const grpcProps = {
    serverAddress,
    serverPort: parseInt(process.env.AGENCY_API_SERVER_PORT!, 10),
    // NOTE: we currently assume that we do not need certs for cloud installation
    // as the cert is issued by a trusted issuer
    certPath: serverAddress === 'localhost' ? certPath : ''
  }

  // Open gRPC connection to agency using authenticator
  return openGRPCConnection(grpcProps, authenticator)
}

const runApp = async () => {
  const { createAgentClient, createProtocolClient } = await setupAgentConnection()

  // Create API clients using the connection
  const agentClient = await createAgentClient()
  const protocolClient = await createProtocolClient()

  // Start listening to agent notifications
  await listenAgent(agentClient, protocolClient)

  app.get('/greet', async (req: Request, res: Response) => {
    const { payload } = await createInvitationPage(agentClient, 'Greet')
    res.send(payload)
  });

  app.get('/issue', async (req: Request, res: Response) => {
    throw 'IMPLEMENT ME!'
  });

  app.get('/verify', async (req: Request, res: Response) => {
    throw 'IMPLEMENT ME!'
  });

  app.get('/', (req: Request, res: Response) => {
    res.send('Typescript example');
  });

  app.listen(port, async () => {
    console.log(`⚡️[server]: Server is running at http://localhost:${port}`);
  });
}

runApp()