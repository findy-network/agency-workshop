# Track 1.1 - Task 3: Prepare for issuing credentials

We have learned how to command our SSI cloud agent. We have also got the first
touch of the chatbot state-machine language. The following tasks will teach us
how to issue and verify credentials using these bots.

But before issuing credentials, we have to prepare our agent for it. It means we
must have a suitable schema and credential definition available.

Schema describes the contents of the verifiable credential (VC): which data
attributes it contains. The credential definition is needed for the [CL signature scheme](https://asecuritysite.com/encryption/cl)
that offers zero-knowledge-proof capabilities. It also combines an ID to
reference all the credentials issued by this issuer from a specific schema.
Other parties can verify these credential's validity against the credential
definition and ensure that your and only your agent has issued the credential.

## 0. Create a new schema

In the terminal window 2:
```shell
# execute next 2 commands only if you have to.
cd "$FCLI_PATH/findy-agent-cli/scripts/fullstack"
cd "play/world/<UUID-from-task1>"
cli agent ping # to check that you are the 'world'
source new-schema email # == one string attribute: 'email'
```

## 1. Create a new schema

Continue in the terminal window 2:
```shell
source new-cred-def # takes time, even tens of seconds in a slower machine
```

## 2. Continue with task 4

Congratulations, you have completed task 3 and you know now how to create schemas
and credential definitions (needed for ZKP and CL signing scheme).

You can now continue with [task 4](../task4/README.md).
