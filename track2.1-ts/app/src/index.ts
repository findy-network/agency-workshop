import express, { Express, Request, Response } from 'express';
import { createAcator, openGRPCConnection } from '@findy-network/findy-common-ts'

const app: Express = express();
const port = process.env.PORT || 3001;

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
  await setupAgentConnection()

  app.get('/greet', async (req: Request, res: Response) => {
    throw 'IMPLEMENT ME!'
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