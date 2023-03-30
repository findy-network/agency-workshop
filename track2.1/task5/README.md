# Track 2.1 - Task 5: Verify credential

Your web wallet user should now have their first credential in their wallet.
Now we can build functionality that will verify that credential.

In a real world implementation we would naturally have two applications and two separate
agents, one for issuing and one for verifying. The wallet user would first acquire a credential
using the issuer application and after that use the credential, i.e. prove the data,
in another application.

For simplicity we build the verification functionality into the same application
we have been working on. The underlying protocol for requesting and presenting proofs is
[the present proof protocol](https://github.com/hyperledger/aries-rfcs/blob/main/features/0037-present-proof/README.md).

## 1
