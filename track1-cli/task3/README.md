# Track 1.1 - Task 3: Prepare for issuing credentials

## Progress

* [Task 0: Setup environment](../README.md)
* [Task 1: Create a new connection](../task1/README.md)
* [Task 2: Send greetings](../task2/README.md)
* [Task 2.5: Chatbot and FSM Language](../task2.5/README.md)
* [**Task 3: Prepare for issuing credentials**](../task3/README.md)
* [Task 4: Issue credential](../task4/README.md)
* [Task 5: Verify credential](../task5/README.md)
* [Task 6: Issue credential for verified information](../task6/README.md)

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

In the terminal window 2 (`$world` agent), we should be in the right place
already. To check that all of these three lines should output `$world` agents
unique name:
```shell
basename `pwd`
echo $world
```

Let's feel the pulse, are we live and kicking?
```shell
$FCLI agent ping
```
Next, we will create a schema with one attribute `email`.
```shell
source new-schema email
```
It should take too much time.

## 1. Create a new credential definition

Continue in the same terminal window 2 (`$world` agent) to create a new CredDef.
```shell
source new-cred-def
```

#### Bonus Task:
> Keeping your self this same terminal session, try these `$FCLI agent
> get-schema` and `$FCLI agent get-cred-def`. They should work because the
> actual IDs are transported thru environment variables.

## 2. Continue with task 4

Congratulations, you have completed task 3 and you know now how to create schemas
and credential definitions (needed for ZKP and CL signing scheme).

You can now continue with [task 4](../task4/README.md).
