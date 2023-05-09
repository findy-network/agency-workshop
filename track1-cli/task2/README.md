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

In the terminal window 1:
```shell
cd "$FCLI_PATH/findy-agent-cli/scripts/fullstack"
cd "play/$hello/<UUID-from-task1>"
cli agent ping # you should see the message: Agent register by name: $hello
cli bot read # now, leave this here until ALL the tasks are done
```

## 1. Send a text message

In the terminal window 2:
```shell
# execute next 2 commands only if you have to.
cd "$FCLI_PATH/findy-agent-cli/scripts/fullstack"
cd "play/$world/<UUID-from-task1>"
echo "Hello world!" | cli bot chat    # look at terminal 1
```

## 2. Continue with task 2.5

Congratulations, you have completed the task and you know now how to send a basic
text message over a DIDComm connection.

You can now continue with [task 2.5](../task2.5/README.md).
