# Track 1.1 - Task 2.5: Chatbot and FSM language

## Progress

* [Task 0: Setup environment](../README.md)
* [Task 1: Create a new connection](../task1/README.md)
* [Task 2: Send greetings](../task2/README.md)
* [**Task 2.5: Chatbot and FSM Language**](../task2.5/README.md)
* [Task 3: Prepare for issuing credentials](../task3/README.md)
* [Task 4: Issue credential](../task4/README.md)
* [Task 5: Verify credential](../task5/README.md)
* [Task 6: Issue credential for verified information](../task6/README.md)

In the previous task, we learned how to send simple text messages between
agents. Now we build our first hello world application with Findy Agency's [FSM
language
(YAML).](https://findy-network.github.io/blog/2023/03/13/no-code-ssi-chatbots-part-i/)

We have two agents `$hello` and `$world`. We want `$world` to be a chatbot and
`$hello` be us, a human. As a human, we need to see both input and output
messages. We already had `$FCLI bot read` command in terminal window 1. You should
split that or bring new terminal near to it. Let's call this new terminal as
'terminal window 1 chat'. And for clarity, the previous terminal window 1 to
'terminal window 1 read'

## 0. Open the human side terminals

In the 'terminal window 1 read' (**we have this open already**):
```shell
cd "<CONNETION_ID_IS_YOUR_CLIPBOARD/OR_PRESS_TAB>"
$FCLI agent ping
```
If last `ping` is OK, everything is ready. Let's start the reader that we hear
what other end is saying to us.
```shell
$FCLI bot read
```
Remember that this can be stopped with `ctrl-C`.

In the 'terminal window 3 chat' (**open/split a new**):
```shell
cd "$FCLI_PATH/findy-agent-cli/scripts/fullstack"
source ./recover-names.sh
cd play/$hello
```
Next, use blow or use connection ID from your clipboard:
```shell
cd ./<PRESS_TAB_AND_WHEN_YOU_SEE_CONN_ID_PRESS_ENTER>
```
Let's test all is good:
```shell
$FCLI agent ping
```
If last `ping` is OK, everything is ready for our other terminal for the same
agent `$hello`. Let's start the chatter that we can send our text lines to other
end.
```shell
$FCLI bot chat
```
Remember that this can be stopped with `ctrl-D` or `ctrl-C`.

## 1. Open chatbot terminal and start the FSM

Before we are ready to start our chat bot we should save it to a file (for
convenience and potential future needs).

```yaml
initial:
  target: INITIAL
states:
  INITIAL:
    transitions:
    - trigger:
        protocol: basic_message
      sends:
      - data: Hello! I'm Hello-World bot.
        protocol: basic_message
      target: INITIAL
```
Save above YAML file to
`$FCLI_PATH/findy-agent-cli/scripts/fullstack/play/$world/hello-world.yaml`

In the terminal window 2 (we have this already for agent `$world`):
```shell
cd "$FCLI_PATH/findy-agent-cli/scripts/fullstack"
source ./recover-names.sh
cd "play/$world"
$FCLI agent ping
```
Next command will start the chatbot state-machine. Don't worry if you see
state-transitions. That can happen if there has been some buffered protocol
events for then agent like `trust_ping`, `did_exchange`, etc.

```shell
$FCLI bot start hello-world.yaml -v=1
```
Dev tip:
> You could open more agent listeners (`$FCLI agent listen`) for both agents:
> `$hello` and `$world`. This helps you keep track what's going on. This is
> especially handy with complex FSM chatbots.

## 2. Be human and communicate with your chatbot

Go back to 'terminal 1 chat' and write something and press `enter`. You should
see the chatbot's reply in your 'terminal window 1 read'.

## 3. Continue with task 3 ...

Congratulations, you have (almost) completed the task 2.5 and you (almost) know
how to write chatbot state-machines with Findy FSM language.

**Bonus Task**:
> Stop (C-c) the current chatbot and modify its declaration so that it also
> echoes the message it receives from the user.

The documentation for Findy FSM is in the previously mentioned [blog
post](https://findy-network.github.io/blog/2023/03/13/no-code-ssi-chatbots-part-i/)
We recommend you to used it as a reference manual. It's written for that kind of
use in mind.

Tip:
> Only two (2) new lines are needed, plus altering the one `data` line with %s.

## 4. Finally, continue with task 3

You can now continue with [task 3](../task3/README.md).
