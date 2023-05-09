# Track 1.1 - Task 5: Verify credential

## Progress

* [Task 0: Setup environment](../README.md)
* [Task 1: Create a new connection](../task1/README.md)
* [Task 2: Send greetings](../task2/README.md)
* [Task 2.5: Chatbot and FSM Language](../task2.5/README.md)
* [Task 3: Prepare for issuing credentials](../task3/README.md)
* [Task 4: Issue credential](../task4/README.md)
* [**Task 5: Verify credential**](../task5/README.md)
* [Task 6: Issue credential for verified information](../task6/README.md)

Your human user should now have their first credential in their wallet. Now
we can build functionality that will verify that credential.

In a real world implementation we would naturally have two applications and two
separate agents, one for issuing and one for verifying. The wallet user would
first acquire a credential using the issuer application and after that use the
credential, i.e. prove the data, in another application.

For simplicity we build the verification functionality into the same application
we have been working on. The underlying protocol for requesting and presenting
proofs is [the present proof
protocol](https://github.com/hyperledger/aries-rfcs/blob/main/features/0037-present-proof/README.md).

## 0. Check your 'holder' aka credential receiver

You should have following command still running in 'terminal 1 chat':
`cli bot chat`

And, you should have following command still running in 'terminal 1 read':
`cli bot read`. **Note that this is very important.**

## 1. Start credential verifier chatbot

In the terminal window 2:
```shell
# execute next 2 commands only if you have to.
cd "$FCLI_PATH/findy-agent-cli/scripts/fullstack"
cd "play/XX-world/<UUID-from-task1>"
cli agent ping # to check that you are the 'XX-world'
printenv | grep FCLI # check that FCLI_CRED_DEF_ID is defined and correct *Tip
cli bot start ../../../email-verifier-bot.yaml -v=1 # verbose lvl 1, we want to know!
```

Tip:
> Defining `alias pf='printenv | grep FCLI'` is very useful idea when working
> with Findy Agency and its CLI.

## 1. Present a proof of the credential

Continue in the terminal window 1 chat by following instructions in
read-terminal.

Note that it's very important that you have reader window for 'XX-hello' agent
running because it accepts important protocol request. There is
`auto-accept-mode` for testing purposes, but we don't need it here.

Sub-task:
> Once again, please read the issuer bot's FSM file thru and try to figure out
> how you could change it. You are welcoming to test those changes. Maybe just
> change what it says?

## 2. Continue with task 6

Congratulations, you have completed task 5 and you now know a lot. You should
have the basic understanding what is core purpose of SSI/DID technology.

For this specific task it should be quite clear how powerful the `proof
presentation` is for cases like authentication. Especially from user's point of
view. And it's suits very well for these kind of chatbot applications.

You can now continue with [task 6](../task6/README.md).

