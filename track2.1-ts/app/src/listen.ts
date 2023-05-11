import { AgentClient, ProtocolClient } from '@findy-network/findy-common-ts'

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
        console.log(`New connection ${didExchange.getTheirLabel()} with id ${info.connectionId}`)
      },
    },
    options
  )
}