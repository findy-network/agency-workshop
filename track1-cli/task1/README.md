# Track 1.1 - Task 1: Create a New Connection

## Progress

* [Task 0: Setup environment](../README.md)
* [**Task 1: Create a new connection**](../task1/README.md)
* [Task 2: Send greetings](../task2/README.md)
* [Task 2.5: Chatbot and FSM Language](../task2.5/README.md)
* [Task 3: Prepare for issuing credentials](../task3/README.md)
* [Task 4: Issue credential](../task4/README.md)
* [Task 5: Verify credential](../task5/README.md)
* [Task 6: Issue credential for verified information](../task6/README.md)

An agent's primary capability is peer-to-peer communication, which allows for
exchanging messages between agents. These interactions can range from simple
plaintext messages to more complex tasks such as negotiating the issuance of a
credential or presenting proof. The peer-to-peer messaging mechanism is called
[DIDComm](https://identity.foundation/didcomm-messaging/spec/), which is short
for DID communication and operates based on the exchange and use of DIDs.

Establishing a DIDComm connection requires one agent to generate an invitation
and transfer the invitation to the other agent, or the first agent could have a
public DID.

Typically the invitation is displayed as a QR code that the other agent can read
using a mobile device. The connection negotiation can then begin using the
information in the invitation. Eventually, the agents have a secure,
e2e-encrypted communication pipeline that they can use to transmit other
protocol messages.

## 0. Allocate Cloud Agents

Go to your 'playground' root:
```shell
cd "$FCLI_PATH/findy-agent-cli/scripts/fullstack"
```

Check that your environment is ready:
```shell
alias pf='printenv | grep FCLI | sort'
pf
```
It should output:
```shell
FCLI_CONFIG=./cfg.yaml
FCLI=findy-agent-cli
FCLI_KEY=1cb85f............cea..............addb7..............0c6122a340
FCLI_ORIGIN=https://f...net
FCLI_PATH=/....s/your-god/work/temp
FCLI_SERVER=fi......................i:50051
FCLI_TLS_PATH=/..............-techlab/work/.....cert
FCLI_URL=https://...................fi
```
If there is something extra like `FCLI_USER`, `FCLI_JWT`, etc. unset them.

> Tip. You have now a `pf` alias in your session. Use it for problem-solving for
> then environment.

When your environment is ready and you are in
`$FCLI_PATH/findy-agent-cli/scripts/fullstack`, execute the following: 
```shell
source agent-name.sh hello world
```
The script generates two unique agent name and stores them variables `$helllo`
and `$world`. It also generates `recover-names.sh` script to easily get those
variables in use in new terminals.

Next is the playground magic. Magic because clients must store something to be
able to keep contact to Findy Agency. The ClI playground does that by using just
file system and OS environment variables. Underneath it executes commands like
`findy-agent-cli auth register --name 'hello-<some-uuid-or-whater>'`. Rest of
the command flags come from `FCLI_URL, etc.` and `./cfg.yaml` file. If you are
interested to see what's the whole command structure, please execute
`findy-agent-cli tree`. It gives the big picture.

Back to business, `make-play-agent.sh` on-boards (allocates) cloud agents from
Findy Agency. It also log in all of the allocated agents.

```shell
./make-play-agent.sh $hello $world
```
That registers and logins two agents for you. The exact names for the agents are
behind the variables: `$hello` and `$world`.

## 1. Create A Listener To Monitor Agent Notifications

For now we think that your are in playground root dir
(`$FCLI_PATH/findy-agent-cli/scripts/fullstack`).

In the terminal window 1, we move our `$hello` agent's home directory:
```shell
cd play/$hello
```
Let's make sure we really are in the right place and everything is OK:
```shell
$FCLI agent ping
```
If it is, we did see response that showed something like:
```shell
Agent registered by name: hello...
```
Next, let's start a listener for our agent to be able to monitor what's going
on. It will, e.g., let us know when other parties are connecting us thru the
`invitation` generated from this agent:
```shell
$FCLI agent listen
```
#### Summary of Our First Findy Agent

After `$FCLI agent ping` you should see the agent's actual name-handle. `listen`
command starts to listen your agent, so that you'll see what's happens when
other agent connects to it.

This the agent who's invitation will be used. For the convenience we'll execute
both invitation creation and invitation connection in the same command line in
the next step.

## 1. Create A Pairwise Connection

In the terminal window 2 (and remember follow terminal 1):
```shell
cd "$FCLI_PATH/findy-agent-cli/scripts/fullstack"
source ./recover-names.sh
cd play/$world
cd $(../$hello/invitation | ./connect)
```
Notice what your current working directory is. The same UUID is printed to
`$hello` agents listen output. Copy that, you'll need it later.

## 2. Verify The Pairwise Connection

And still in the terminal window 2 (and look at terminal 1):
```shell
$FCLI connection trustping
```
The `trustping` verifies that the pairwise connection is well working properly.

## 3. Continue With Task 2

Congratulations, you have completed the task and you know now how to establish
DIDComm connections between agents for message exchange!

You can now continue with [task 2](../task2/README.md).
