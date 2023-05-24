import { agencyv1, AgentClient, ProtocolClient } from '@findy-network/findy-common-ts'
import { Issuer } from './issue'
import { Verifier } from './verify'

export default async (
  agentClient: AgentClient,
  protocolClient: ProtocolClient,
  issuer: Issuer,
  verifier: Verifier,
) => {

  // Options for listener
  const options = {
    protocolClient,
    retryOnError: true,
  }

  // Listening callback handles agent events
  await agentClient.startListeningWithHandler(
    {
      // New connection is established
      DIDExchangeDone: async (info, didExchange) => {
        console.log(`New connection: ${didExchange.getTheirLabel()} with id ${info.connectionId}`)

        // Greet each new connection with basic message
        const msg = new agencyv1.Protocol.BasicMessageMsg()
        msg.setContent('Hi there 👋!')
        await protocolClient.sendBasicMessage(info.connectionId, msg)

        // Notify issuer of new connection
        issuer.handleNewConnection(info, didExchange)

        // Notify verifier of new connection
        verifier.handleNewConnection(info, didExchange)
      },

      BasicMessageDone: async (info, basicMessage) => {
        // Print out greeting sent from the other agent
        if (!basicMessage.getSentByMe()) {
          const msg = basicMessage.getContent()
          console.log(`Received basic message ${msg} from ${info.connectionId}`)
        }

        // Notify issuer
        issuer.handleBasicMessageDone(info, basicMessage)
      },

      IssueCredentialDone: (info, issueCredential) => {
        // Notify issuer of issue protocol success
        issuer.handleIssueDone(info, issueCredential)
      },

      PresentProofPaused: (info, presentProof) => {
        verifier.handleProofPaused(info, presentProof)
      },

      PresentProofDone: (info, presentProof) => {
        verifier.handleProofDone(info, presentProof)
      },
    },
    options
  )
}