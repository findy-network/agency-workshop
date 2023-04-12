# Track 1.1: CLI & Chatbot Application

In this track, you will learn how to build a chatbot application that utilizes
Findy Agency FSM language. The assumption is that you are working in a guided
workshop with the default tooling. In this case you can skip the sections with
symbol 🤠.

Follow the instructions and execute the tasks in order. Good luck!
It's needed. Both, the CLI and the FSM chatbot language are still under
development.

## Task 0: Setup environment

### 1. Clone these repositories into your GOPATH/src or into what you prefer

```shell
git clone https://github.com/findy-network/agency-workshop.git
git clone https://github.com/findy-network/findy-agent-cli.git
```

As you noticed, you are cloning the `findy-agent-cli` repo as well to get *the
actual helper scripts* for your use.

The scripts are located `scripts/fullstack/`. The directory contains `README.md`
where some of the scripts are documented. It presents few examples as well.

### 2. Install The Tooling

The prerequisites for the FSM track are:
1. You should have a Unix and terminal access. The shell can be whatever, but
  `bash` is preferred because it's mostly tested by the dev team.
2. [`findy-agent-cli`](https://github.com/findy-network/findy-agent-cli#installation)
   installed.
3. And most importantly, you should have access to [Findy
   Agency](https://findy-network.github.io). After CLI is successfully installed
   you need to setup it's execution environment, i.e. bind it to the Findy Agency.

Everything presented here can be executed just by using `findy-agent-cli` (later
just `cli`) and Unix terminal & shell. (Terminal multiplexers and tiling window 
managers might help you during these tasks.)

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

Run following script in the terminal:

```bash
source <(curl <agency_url>/set-env.sh)
```

The agency URL is provided for you in the guided workshop. e.g.
`https://agency.example.com`.

The script will export the needed environment variables. It will also create
file `.envrc` that contains these variables. Typing `direnv allow` will ensure
that the variables are automatically exported when you open a new terminal
window in this folder.

If you don't have `direnv` installed, you can export the variables by typing
`source .envrc`.
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
**but you don't need docker or network access**.

In the case you want to play with the sources or you want to get touch of how the
whole system feels to run locally from sources, 
see instructions [here](../agency-native/README.md). There is a script
(`setup.sh`) which installs all the needed repos and a tmuxinator script to start
the system playground. The script targets a Debian Linux.

Here's the summary what should be done:

Clone the needed Agency service source repos:
```shell
git clone https://github.com/findy-network/findy-agent-auth.git
git clone https://github.com/findy-network/findy-agent.git
git clone https://github.com/findy-network/findy-agent-cli.git
```

Start the FIDO2 Server:
```shell
cd <findy-agent-auth-repo>
cd scripts; ./mem-dev-server.sh
```

Start the Agency Core Server:
```shell
cd <findy-agent-repo>
make cli # builds fa named binary
cd scripts/test
fa ledger steward create --config create-steward-to-mem-ledger.yaml
agency=fa register=findy.json no_clean=1 enclave=MEMORY_enclave.bolt ./mem-server --reset-register --grpc-cert-path ../../grpc/cert
```

Start the Findy Agent CLI to command your local agency (in a new terminal/window/tab):
```shell
cd <findy-agent-cli-repo>
make cli # builds and installs binary named cli in your path
cd scripts/fullstack
source ./setup-cli-env-local.sh
admin/register && . admin/login
cli agency count # tells how many cloud agent/wallet is running/onboarded
```

After you have verified that everything above works, you can allocate two
separate SSI agents:
```shell
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
# like (tmux kill-session)
```

If you want to use tmux and tmuxinator the previously mentioned `setup.sh`
script includes tmuxinator configuration that is installed by it with the name
`play`.
```shell
tmuxinator play
```
Tip:
> You can use tmuxinator configurations even when using cloud version of the
> agency. Just check from where the environment variables are loaded.

</details><br/>
### 5. Continue with task 1

Congratulations, you have completed the task and have a working agency CLI
development environment available!

You can now continue with [task 1](./task1/README.md).
