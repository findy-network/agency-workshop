# Track 1.1: Chatbot Application

In this track, you will learn how to build a chatbot application that utilizes
Findy Agency FSM language. The assumption is that you are working in a guided
workshop with the default tooling. In this case you can skip the sections with
symbol 🤠.

Follow the instructions carefully and execute the tasks in order. Good luck!
It's needed, because both CLI and FSM chatbot language are still under
development.

## Task 0: Setup environment

### 1. Clone this repository to your workspace

```shell
git clone https://github.com/findy-network/agency-workshop.git
```

<details>
<summary>🤠 Extra scripts</summary>

Also clone `findy-agent-cli` to get more helper scripts for your use.
```shell
git clone https://github.com/findy-network/findy-agent-cli.git
```
The scripts are located `scripts/fullstack/`. The directory contains `README.md`
where some of the scripts are documented. It presents few examples as well.

</details><br/>

### 2. Install The Tooling

The prerequisites for the FSM track are:
1. You should have a Unix and terminal access. The shell can be what ever, but
  `bash` is preferred because it's mostly used by the dev team.
2. [`findy-agent-cli`](https://github.com/findy-network/findy-agent-cli#installation)
   installed.
3. And most importantly, you should have access to [Findy
   Agency](https://findy-network.github.io). After CLI is successfully installed
   you need to setup it's execution environment, i.e. bind it to the Findy Agency.

Everything presented here can be executed just by using `findy-agent-cli` and
Unix terminal & shell.

<details>
<summary>🤠 Other options</summary>

You can also set up some extra tools. However, these instructions describe only
how to work with the recommended tooling.

If you still wish to go to the dark side, make sure you have these tools available:

- tmux or screen
- [direnv](https://direnv.net/) (*optional*, but useful)
- FIDO2 capable Web browser of your choice

</details><br/>

### 3. 🤠 Install Findy Agency

If you are participating in a guided workshop, you will likely have a cloud
installation of Findy Agency available. Skip this step.

<details>
<summary>🤠 Local setup</summary>

Start local agency instance if you do not have cloud installation available.
See instructions [here](../agency-local/README.md).

</details><br/>

### 4. Set environment variables

The agency environment provides a script for setting up the needed environment
variables automatically.

Run following script in the dev container terminal:

```bash
source <(curl <agency_url>/set-env.sh)
```

The agency URL is provided for you in the guided workshop. e.g.
`https://agency.example.com`

The script will export the needed environment variables. It will also create
file `.envrc` that contains these variables. Typing `direnv allow` will ensure
that the variables are automatically exported when you open a new terminal
window in this folder.

If you don't have direnv installed, you can export the variables by typing `source .envrc`.

<details>
<summary>🤠 Local setup (WebServer&docker)</summary>

For [local agency
installation](https://github.com/findy-network/findy-wallet-pwa/blob/master/tools/env/README.md#agency-setup-for-local-development),
use the web wallet URL `http://localhost:3000`:

```bash
source <(curl http://localhost:3000/set-env.sh)
```

</details><br/>

<details>
<summary>🤠 Local setup (Native from sources)</summary>

You need to have Go 1.20 installed to run needed Agency services from sources:
**but you don't need docker and network access**.

Clone the needed Agency service source repos:
```console
git clone https://github.com/findy-network/findy-wrapper-go.git
git clone https://github.com/findy-network/findy-agent-auth.git
git clone https://github.com/findy-network/findy-agent.git
git clone https://github.com/findy-network/findy-agent-cli.git
```

Start the FIDO2 Server:
```console
cd <findy-agent-auth-repo>
cd scripts; ./mem-dev-server.sh
```

Start the Agency Core Server:
```console
cd <findy-agent-repo>
cd scripts/test
no_clean=1 enclave=MEMORY_enclave.bolt ./mem-server --reset-register --grpc-cert-path ../../grpc/cert
```

Start the Findy Agent CLI to command your local agency (in a new terminal/window/tab):
```console
cd <findy-agent-cli-repo>
make cli # builds and installs binary named cli in your path
cd scripts/fullstack
source ./setup-cli-env-local.sh
admin/register && . admin/login
cli agency count # tells how many cloud agent/wallet is running/onboarded
```

After you have verified that everything above works with:
```console
# continue in findy-agent-cli/scripts/fullstack 
./make-play-agent.sh test-alice test-bob
pushd test-alice
cli agent ping
# do something else with `test-alice` and `test-bob` like:
cd $(./invitation | ../test-bob/connect)
cli connection trustping
popd
./rm-play-agent.sh test-alice test-bob # cleanup wallets and client stores
# typically you shutdown FIDO2 and Core servers at this point
```

Now you can build whole new starting script to bring your local Findy Agency up
and open some windows to play with it. Below is a `tmuxinator` example to use
with `tmux`:

```yaml
# tmuxinator file to start local Findy Agency playground
name: play
# In this example findy-network is in GOPATH but you should use the common root
# of your previously cloned 3 findy repos.
root: ~/go/src/github.com/findy-network

windows:
  - devops:
      layout: main-vertical
      panes:
        - # empty shell
        - # empty shell
  - running:
      layout: tiled
      panes:
        - cd findy-agent-auth/scripts; ./mem-dev-run.sh
        - cd findy-agent/scripts/test; no_clean=1 enclave=MEMORY_enclave.bolt ./mem-server --reset-register
        - cd findy-agent-cli/scripts/fullstack; source ./setup-cli-env-local.sh
        - cd findy-agent-cli/scripts/fullstack; source ./setup-cli-env-local.sh
        - cd findy-agent-cli/scripts/fullstack; source ./setup-cli-env-local.sh
        - cd findy-agent-cli/scripts/fullstack; source ./setup-cli-env-local.sh
```
</details><br/>

<details>
<summary>🤠 No direnv?</summary>


</details><br/>

![Script output](./docs/environment-direnv.png)

*Note! By default, the script will generate a generic username for your client.
If you wish to use a more descriptive name for your app, define it before running the script:*

```bash
export FCLI_USER="my-fancy-issuer-service"

source <(curl <agency_url>/set-env.sh)
```

*The username needs to be unique in the agency context.*

### 6. Start the application for the first time

  When starting the application for the first time, run following commands:

  ```bash
  nvm use
  npm install
  npm run build
  npm run dev     # start server in watch mode
  ```

  After the first start, you can use just `npm run dev`.

  When the server is started, VS Code displays a dialog telling where to find the app.

  ![Application running](./docs/application-running.png)

  Click "Open in Browser". The browser should open to URL <http://localhost:3001>
  and display the text "Typescript example".

  Now you have a simple express server running in port 3001 with four endpoints:
  `/`, `/greet`, `/issue` and `/verify`. Next step is to start adding some actual code
  to the server skeleton.

### 7. Create connection to the agency

Add a new dependency to your project:

```bash
npm install @findy-network/findy-common-ts --save
```

`findy-common-ts` library has functionality that helps us authenticate to the agency
or use the agent capabilities.

Open file `src/index.ts`.

Add following row to imports:

```ts
import { createAcator, openGRPCConnection } from '@findy-network/findy-common-ts'
```

Create new function `setupAgentConnection`:

```ts
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
```

This function will open a connection to our agent. Through this connection, we can control
the agent and listen for any events the agent produces while handling our credential protocol
flows.

We authenticate the client using a headless FIDO2 authenticator provided by the agency helper
library. When opening the connection for the first time, the underlying functionality
automatically registers the authenticator to our agent.

The `FCLI_KEY` variable contains the master key to your authenticator. It is generated during
the development environment setup. (In production the key should be naturally generated and
handled in a secure manner as a secret). If someone gets access to the key,
they can control your agent.

Add call to `setupAgentConnection` to existing `runApp` function:

```ts
const runApp = async () => {

  await setupAgentConnection()

  ...
}
```

As you can see from the logs, the authentication fails at first as the client is not yet registered.
With further server starts, this error should disappear.

Verify that you see logs similar to this:
![First login log](./docs/log-first-login.png)

### 8. Continue with task 1

Congratulations, you have completed task 0 and have
a working agency client development environment available!

You can now continue with [task 1](./task1/README.md).
