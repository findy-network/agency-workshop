# Track 1.1 - Task 2: Send greetings

## Progress

* [Task 0: Setup environment](../README.md)
* [Task 1: Create a new connection](../task1/README.md)
* [**Task 2: Send greetings**](../task2/README.md)
* [Task 2.5: Chatbot and FSM Language](../task2.5/README.md)
* [Task 3: Prepare for issuing credentials](../task3/README.md)
* [Task 4: Issue credential](../task4/README.md)
* [Task 5: Verify credential](../task5/README.md)
* [Task 6: Issue credential for verified information](../task6/README.md)

In the previous task, we learned how to establish e2e-encrypted messaging
between agents. Now we will send our first messages using the communication
pairwise.

Agents interact over DIDComm using a specific Hyperledger Aries protocols. There
are different protocols for different purposes. Agents send text messages to
each other using [basic message
protocol](https://github.com/hyperledger/aries-rfcs/blob/main/features/0095-basic-message/README.md).

## 0. Receive a text message

In the 'terminal window 1 read' (**we have this open already**):
NOTE. **That if you haven't moved, you are in the *almost* correct place**.
```shell
cd ./<PRESS_TAB_IN_HERE_AND_ENTER_WHEN_YOU_SEE_CONN_ID>
```
Otherwise, move to right place and setup the environment thru the findy-network 
root:
```shell
cd "$FCLI_PATH/findy-agent-cli/scripts/fullstack"
source ./recover-names.sh
cd "play/$hello/<UUID-from-task1>"
```
In the agent's connection directory (that's the `UUID` in path), let's do the
health check with then `ping` command first.
```shell
$FCLI agent ping
```

On then start the reader to show what other party is saying. The difference
between the command `agent listen` and `bot read` is that first is agent level
and pretty technical, and the second simulates chat applications history stream.

```shell
$FCLI bot read
```
When you want to stop reading, just hit `ctrl-C`. Note. These listeners are
event-based and don't consume CPU cycles when idling.

## 1. Send a text message

In the terminal window 2 for `$world` agent, we are in the right place.

Let's do the health check with then `ping` command first.
```shell
$FCLI agent ping
```
Let's send one text line thru the DIDComm to the agent `$hello`:
```shell
echo 'Hello world!' | $FCLI bot chat
```
> Tip. If you give just `$FCLI bot chat`, you can enter and send many text lines
> to other end. To stop that give `ctrl-D`

Now, remember, your `$world` agent and its terminal, they are in the right
directory. Stay there it makes next step easier.

You can output the connection ID need in the next step with:
```shell
basename `pwd`
```
Copy it to clipboard.

## 2. Continue with task 2.5

Congratulations, you have completed the task and you know now how to send a basic
text message over a DIDComm connection.

You can now continue with [task 2.5](../task2.5/README.md).
