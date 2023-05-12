import { agencyv1, ProtocolClient, ProtocolInfo } from '@findy-network/findy-common-ts'
import { ProtocolStatus } from '@findy-network/findy-common-ts/dist/idl/protocol_pb'

export interface Issuer {
  addInvitation: (id: string) => void
  handleNewConnection: (info: ProtocolInfo, didExchange: ProtocolStatus.DIDExchangeStatus) => Promise<void>
  handleIssueDone: (info: ProtocolInfo, issueCredential: ProtocolStatus.IssueCredentialStatus) => void
}

interface Connection {
  id: string
}

export default (protocolClient: ProtocolClient, credDefId: string) => {
  const connections: Connection[] = []

  const addInvitation = (id: string) => {
    connections.push({ id })
  }

  const handleNewConnection = async (info: ProtocolInfo, didExchange: ProtocolStatus.DIDExchangeStatus) => {
    // Skip if this connection was not for issuing
    const connection = connections.find(({ id }) => id === info.connectionId)
    if (!connection) {
      return
    }

    // Create credential content
    const attributes = new agencyv1.Protocol.IssuingAttributes()
    const attr = new agencyv1.Protocol.IssuingAttributes.Attribute()
    // Attribute name
    attr.setName('foo')
    // Attribute value
    attr.setValue('bar')
    attributes.addAttributes(attr)

    const credential = new agencyv1.Protocol.IssueCredentialMsg()
    credential.setCredDefid(credDefId)
    credential.setAttributes(attributes)

    // Send credential offer to the other agent
    console.log(`Sending credential offer\n${JSON.stringify(credential.toObject())}\nto ${info.connectionId}`)
    await protocolClient.sendCredentialOffer(connection.id, credential)
  }

  const handleIssueDone = (info: ProtocolInfo, issueCredential: ProtocolStatus.IssueCredentialStatus) => {
    console.log(`Credential\n${JSON.stringify(issueCredential.toObject())}\nwith protocol id ${info.protocolId} issued to ${info.connectionId}`)

    // Remove connection from cache
    const index = connections.findIndex(({ id }) => id === info.connectionId)
    connections.splice(index, 1)
  }

  return {
    addInvitation,
    handleNewConnection,
    handleIssueDone
  }
}