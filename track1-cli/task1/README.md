# Track 1.1 - Task 1: Create a new connection

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

## 0. Allocate cloud agents

Go to your 'playground' root:
```shell
cd "$FCLI_PATH/findy-agent-cli/scripts/fullstack" # playground root
```
or
```shell
cd ./findy-agent-cli/scripts/fullstack" # playground root
```

Check that your env is ready:
```shell
printenv | grep FCLI
```
It should output:
```shell
FCLI_URL=https://...fi
FCLI_SERVER=f...fi:50051
FCLI_KEY=d92f5847..............1897900457da599d
FCLI=cli
FCLI_TLS_PATH=/dir/path/to/cert
FCLI_CONFIG=./cfg.yaml
```
If there is something extra like `FCLI_USER`, `FCLI_JWT`, etc. unset them.

When your env is ready and you are in
`$FCLI_PATH/findy-agent-cli/scripts/fullstack`, execute the following:
```shell
# --- try to think something uniq for the XX in the class room :-)
./make-play-agent.sh XX-hello XX-world
```
That onboards two agents for you.

## 1. Create a listener to monitor agent notifications

For now we think that your are in playground root dir
(`$FCLI_PATH/findy-agent-cli/scripts/fullstack`).

In the terminal window 1:
```shell
cd play/XX-hello
cli agent ping # you should see the message: Agent register by name: XX-hello
cli agent listen # terminate with C-c when step 2 is finished
```
This the agent who's invitation will be used. For the convenience we'll execute
both invitation creation and invitation connection in the same command line in
the next step.

## 1. Create a pairwise connection

In the terminal window 2:
```shell
cd "$FCLI_PATH/findy-agent-cli/scripts/fullstack"
cd play/XX-world
cd $(../XX-hello/invitation | ./connect)    # look at terminal 1
```

## 2. Verify the pairwise connection

And still in the terminal window 2:
```shell
cli connection trustping                 # look at terminal 1
```
The `trustping` verifies that the pairwise connection is well working properly.

## 3. Continue with task 2

Congratulations, you have completed the task and you know now how to establish
DIDComm connections between agents for message exchange!

You can now continue with [task 2](../task2/README.md).
