# Track 1.1 - Task 4: Issue credential

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
`cli bot chat`

And you should have following command still running in 'terminal 1 read':
`cli bot read`.  **Note that this is very important.** The `read` command also
accept important protocol requests coming from the agency.

## 1. Start credential issuer chatbot

In the terminal window 2:
```shell
# execute next 2 commands only if you have to.
cd "$FCLI_PATH/findy-agent-cli/scripts/fullstack"
cd "play/world/<UUID-from-task1>"
cli agent ping # to check that you are the 'world'
printenv | grep FCLI # check that FCLI_CRED_DEF_ID is defined and correct *Tip
cli bot start ../../../email-issuer-bot.yaml -v=1 # verbose lvl 1, we want to know!
```
Tip:
> Defining `alias pf='printenv | grep FCLI'` is very useful idea when working
> with Findy Agency and its CLI. If you haven't done this already.

## 1. Receive the credential

Continue in the terminal window 1 chat by following instructions in
read-terminal **and the secret PIN-code is printed to terminal window 2 where
chatbot is running**.

Note that it's very important that you have reader window for 'hello' agent
running because it accepts important protocol request. There is
`auto-accept-mode` for testing purposes, but we don't need it here.

Sub-task:
> Please read the issuer bot's FSM file thru and try to figure out how you
> could change it. You are welcoming to test those changes. Maybe just change
> what it says?

## 2. Continue with task 5

Congratulations, you have completed task 4 and you now know kung fu.

You can now continue with [task 5](../task5/README.md).

