# Track 1.1 - Task 4: Issue credential

## Progress

* [Task 0: Setup environment](../README.md)
* [Task 1: Create a new connection](../task1/README.md)
* [Task 2: Send greetings](../task2/README.md)
* [Task 2.5: Chatbot and FSM Language](../task2.5/README.md)
* [Task 3: Prepare for issuing credentials](../task3/README.md)
* [**Task 4: Issue credential**](../task4/README.md)
* [Task 5: Verify credential](../task5/README.md)
* [Task 6: Issue credential for verified information](../task6/README.md)

Now that we have completed setting up the basic bells and whistles we can start
the real fun with issuing and verifying credentials.

Let's first issue our email credential to the human user and verify it after
that. In a real world application the issuer would naturally know something
about the user and would issue a credential only with validated information.
In this example we have pseudo-verification, i.e, we simulate the case where a
bot generates 6 digit PIN-code that is printed in the log output. In the
real-world case the actual email should be used.

Agents use [the issue credential
protocol](https://github.com/hyperledger/aries-rfcs/blob/main/features/0036-issue-credential/README.md)
when handling the issuing process. Luckily, Findy Agency handles execution of
this complex protocol for us (similarly as with other Hyperledger Aries
protocols).

## 0. Check your 'holder' aka credential receiver

You should have following command still running in 'terminal 1 chat':
`$FCLI bot chat`

And you should have following command still running in 'terminal 1 read':
`$FCLI bot read`.  **Note that this is very important.** The `read` command also
accept important protocol requests coming from the agency.

## 1. Start credential issuer chatbot

In the terminal window 2, where you already are in the correct place, in the
agent `$world` home directory:
```shell
$FCLI bot start ../../email-issuer-bot.yaml -v=1
```
This bot is also passive and waits first message from other party, i.e., please
start the conversation from `$hello` agent's `chat` terminal.

## 1. Receive the credential

Continue in the terminal window 1 chat by following instructions in
read-terminal **and the secret PIN-code is printed to terminal window 2 where
chatbot is running**. Enter the PIN-code when asked (`read` window). After that
you should see something like `Thank you` and `Bye Bye`. And it was all. Now
your `$hello` agent has `email` credential.

Note that it's very important that you have reader window for '$hello' agent
running because it accepts important protocol request. There is
`auto-accept-mode` for testing purposes, but we don't need it here.

Sub-task:
> Please read the issuer bot's FSM file thru and try to figure out how you
> could change it. You are welcoming to test those changes. Maybe just change
> what it says?

## 2. Continue with task 5

Congratulations, you have completed task 4 and you now know kung fu.

You can now continue with [task 5](../task5/README.md).

