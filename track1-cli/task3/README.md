# Track 1.1 - Task 3: Prepare for issuing credentials

We have learned how to command our SSI cloud agent. In following tasks we will
learn how to issue and verify credentials using similar APIs.

But before we can start issuing credentials, we have to prepare our agent for
it. It means that we have to have suitable schema and credential definition
available.

Schema describes the contents of the verifiable credential: which data
attributes it contains. Credential definition is like a public key published
against that schema. Other parties can verify the validity of the credential
against the credential definition, and make sure that your and only your agent
has issued the credential.

## 0. Create a new schema

In the terminal window 2:
```shell
# execute next 2 commands only if you have to.
cd $FCLI_PATH/findy-agent-cli/scripts/fullstack
cd play/world/<UUID-from-task1>
cli agent ping # to check that you are the 'world'
source new-schema foo
```

## 1. Create a new schema

Continue in the terminal window 2:
```shell
source new-cred-def # this can take time, even tens of seconds in a slower machine
```

## 2. Continue with task 4

Congratulations, you have completed task 3 and you know now how to create schemas
and credential definitions (needed for ZKP and CL signing scheme).

You can now continue with [task 4](../task4/README.md).
