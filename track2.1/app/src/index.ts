import express, { Express, Request, Response } from 'express';

const app: Express = express();
const port = process.env.PORT || 3001;

const runApp = async () => {
  app.get('/greet', async (req: Request, res: Response) => {
    throw "IMPLEMENT ME!"
  });

  app.get('/issue', async (req: Request, res: Response) => {
    throw "IMPLEMENT ME!"
  });

  app.get('/verify', async (req: Request, res: Response) => {
    throw "IMPLEMENT ME!"
  });

  app.get('/', (req: Request, res: Response) => {
    res.send('Typescript example');
  });

  app.listen(port, async () => {
    console.log(`⚡️[server]: Server is running at http://localhost:${port}`);
  });
}

runApp()