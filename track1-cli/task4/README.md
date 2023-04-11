# Track 1.1 - Task 4: Issue credential

Now that we have completed setting up the basic bells and whistles we can start
the real fun with issuing and verifying credentials.

Let's first issue a dummy credential to the web wallet user and verify it after
that. In a real world application the issuer would naturally know something
about the user and would issue a credential only with validated information. But
for this example case, we issue a dummy credential to all users that connect
with us.

Agents use [the issue credential
protocol](https://github.com/hyperledger/aries-rfcs/blob/main/features/0036-issue-credential/README.md)
when handling the issuing process. Luckily, Findy Agency handles execution of
this complex protocol for us (similarly as with other Hyperledger Aries
protocols).

## 0. Check your 'holder' aka credential receiver

You should have following command still running in 'terminal 1 chat':
`cli bot chat`

And, you should have following command still running in 'terminal 1 read':
`cli bot read`

## 1. Start issuer chatbot

In the terminal window 2:
```shell
# execute next 2 commands only if you have to.
cd $FCLI_PATH/findy-agent-cli/scripts/fullstack
cd play/world/<UUID-from-task1>
cli agent ping # to check that you are the 'world'
cli bot start ../../../email-issuer-bot.yaml -v=1 # verbose lvl 1, we want to know!
```

## 1. Receive the credential

Continue in the terminal window 1 chat by following instructions in
read-terminal **and the secret PIN-code is printed to terminal window 2 where
chatbot is runnnig.

Tip:
> Please read the issuer bot's FSM file thru and try to figure out how you
> could change it. You are welcoming to test those changes.

## 2. Continue with task 5

Congratulations, you have completed task 4 and you know now something.. :-D

You can now continue with [task 5](../task5/README.md).

