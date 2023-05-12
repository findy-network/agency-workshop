import { agencyv1, AgentClient, ProtocolClient } from '@findy-network/findy-common-ts'

export default async (
  agentClient: AgentClient,
  protocolClient: ProtocolClient,
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
      },

      BasicMessageDone: async (info, basicMessage) => {
        // Print out greeting sent from the other agent
        if (!basicMessage.getSentByMe()) {
          const msg = basicMessage.getContent()
          console.log(`Received basic message ${msg} from ${info.connectionId}`)
        }
      },
    },
    options
  )
}