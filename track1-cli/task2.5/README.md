# Track 1.1 - Task 2.5: Chat bot state machines and FSM language

In the previous task, we learned how to send simple text messages between
agents. Now we build our first hello world application with Findy Agency's [FSM
language
(YAML).](https://findy-network.github.io/blog/2023/03/13/no-code-ssi-chatbots-part-i/)

We have two agents `hello` and `world`. We want `world` to be a chatbot and
`hello` be us, a human. As a human, we need to see both input and output
messages. We already had `cli bot read` command in terminal window 1. You should
split that or bring new terminal near to it. Let's call this new terminal as
'terminal window 1 chat'. And for clarity, the previous terminal window 1 to
'terminal window 1 read'

## 0. Open the human side terminals

In the 'terminal window 1 read' (**we have this already**):
```shell
cd $FCLI_PATH/findy-agent-cli/scripts/fullstack
cd play/hello/<UUID-from-task1>
cli agent ping # you should see the message: Agent register by name: hello
cli bot read # leave this here until ALL the tasks are done
```

In the 'terminal window 1 chat' (**open/split a new**):
```shell
cd $FCLI_PATH/findy-agent-cli/scripts/fullstack
cd play/hello/<UUID-from-task1>
cli agent ping # you should see the message: Agent register by name: hello
cli bot chat # leave this here until ALL the tasks are done
```

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
`$FCLI_PATH/findy-agent-cli/scripts/fullstack/play/world/hello-world.yaml`

In the terminal window 2:
```shell
# execute next 2 commands only if you have to.
cd $FCLI_PATH/findy-agent-cli/scripts/fullstack
cd play/world/<UUID-from-task1>
cli agent ping # we're the world!
cli bot start hello-world.yaml -v=1 # verbose lvl 1, we want to know!
```
Dev tip:
> You could open more agent listeners (`cli agent listen`) for both agents:
> `hello` and `world`. This helps you keep track what's going on. This is
> especially handy with complex FSM chatbots.

## 2. Be human and communicate with your chatbot

Go back to 'terminal 1 chat' and write something and press `enter`. You should
see the chatbot's reply in your 'terminal window 1 read'.

## 3. Continue with ~task 3~ ...

Congratulations, you have (almost) completed the task 2.5 and you (almost) know
how to write chatbot state-machines with Findy FSM language.

**Super-User** Task:
> Stop (C-c) the current chatbot and modify its declaration so that it also
> echoes the message it receives from the user.

Unfortunately, the only documentation for Findy FSM is in the previously
mentioned [blog
post](https://findy-network.github.io/blog/2023/03/13/no-code-ssi-chatbots-part-i/)
We recommend you to used it as a reference manual. It's written for that kind of
use in mind.

Tip:
> Only two (2) new lines are needed, plus altering the one `data` line with %s.

## 4. Finally, continue with task 3

You can now continue with [task 3](../task3/README.md).