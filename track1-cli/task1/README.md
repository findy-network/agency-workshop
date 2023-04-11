# Track 1.1 - Task 1: Create a new connection

An agent's primary capability is peer-to-peer messaging, which allows for
exchanging messages between agents. These interactions can range from simple
plaintext messages to more complex tasks such as negotiating the issuance of a
credential or presenting proof. The peer-to-peer messaging mechanism is called
DIDComm, which is short for DID communication and operates based on the exchange
and use of DIDs.

Establishing a DIDComm connection requires one agent to generate an invitation
and transfer the invitation to the other agent, or the first agent could have a
public DID instead of invitation.

Typically the invitation is displayed as a QR code that the other agent can read
using a mobile device. The connection negotiation can then begin using the
information in the invitation. Eventually, the agents have a secure,
e2e-encrypted communication pipeline that they can use to transmit other
protocol messages.

## 0. Allocate a cloud agents

```shell
export FCLI_PATH=<set_your_findy_network_dir_here> # this will be used later
cd $FCLI_PATH/findy-agent-cli/scripts/fullstack
./make-play-agent.sh hello world
```

## 1. Create a pairwise connection

In the terminal window 1:
```shell
cd $FCLI_PATH/findy-agent-cli/scripts/fullstack
./make-play-agent.sh hello world
cd play/hello
cli agent ping # you should see the message: Agent register by name: hello
cli agent listen # terminate with C-c when step 2 is finished
```

## 2. Verify the pairwise connection

In the terminal window 2:
```shell
cd $FCLI_PATH/findy-agent-cli/scripts/fullstack
cd play/world
cd $(../hello/invitation | ./connect)    # look at terminal 1
cli connection trustping                 # look at terminal 1
```
## 3. Continue with task 2

Congratulations, you have completed task 1 and you know now how to establish
DIDComm connections between agents for message exchange!

You can now continue with [task 2](../task2/README.md).
