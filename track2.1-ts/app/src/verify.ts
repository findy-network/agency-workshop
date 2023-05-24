import { agencyv1, ProtocolClient, ProtocolInfo } from '@findy-network/findy-common-ts'
import { ProtocolStatus } from '@findy-network/findy-common-ts/dist/idl/protocol_pb'

export interface Verifier {
  addInvitation: (id: string) => void
  handleNewConnection: (info: ProtocolInfo, didExchange: ProtocolStatus.DIDExchangeStatus) => Promise<void>
  handleProofPaused: (info: ProtocolInfo, presentProof: ProtocolStatus.PresentProofStatus) => void
  handleProofDone: (info: ProtocolInfo, presentProof: ProtocolStatus.PresentProofStatus) => void
}

export default (protocolClient: ProtocolClient, credDefId: string) => {
  const invitations: string[] = []

  const addInvitation = (id: string) => {
    invitations.push(id)
  }

  const handleNewConnection = async (info: ProtocolInfo, didExchange: ProtocolStatus.DIDExchangeStatus) => {
    // Skip if this connection was not for verifying
    if (!invitations.includes(info.connectionId)) {
      return
    }

    // Create proof request
    const attributes = new agencyv1.Protocol.Proof()
    const attr = new agencyv1.Protocol.Proof.Attribute()
    attr.setName("foo")
    attr.setCredDefid(credDefId)
    attributes.addAttributes(attr)

    const proofRequest = new agencyv1.Protocol.PresentProofMsg()
    proofRequest.setAttributes(attributes)

    // Send proof request to the other agent
    console.log(`Sending proof request\n${JSON.stringify(proofRequest.toObject())}\nto ${info.connectionId}`)
    await protocolClient.sendProofRequest(info.connectionId, proofRequest)
  }

  const handleProofPaused = async (info: ProtocolInfo, presentProof: ProtocolStatus.PresentProofStatus) => {
    console.log(`Proof\n${JSON.stringify(presentProof.toObject())}\nwith protocol id ${info.protocolId} paused from ${info.connectionId}`)

    // This function is called after proof is verified cryptographically.
    // The application can execute its business logic and reject the proof
    // if the attribute values are not valid.
    const protocolID = new agencyv1.ProtocolID()
    protocolID.setId(info.protocolId)
    protocolID.setTypeid(agencyv1.Protocol.Type.PRESENT_PROOF)
    protocolID.setRole(agencyv1.Protocol.Role.RESUMER)
    const msg = new agencyv1.ProtocolState()
    msg.setProtocolid(protocolID)

    // We have no special logic here - accept all received values
    msg.setState(agencyv1.ProtocolState.State.ACK)
    console.log(`Resuming proof with for protocol ${info.protocolId} with payload ${JSON.stringify(msg.toObject())}`)
    await protocolClient.resume(msg)
  }

  const handleProofDone = (info: ProtocolInfo, presentProof: ProtocolStatus.PresentProofStatus) => {
    console.log(`Proof\n${JSON.stringify(presentProof.toObject())}\nwith protocol id ${info.protocolId} verified from ${info.connectionId}`)

    // Remove invitation id from cache
    const index = invitations.indexOf(info.connectionId)
    invitations.splice(index, 1)
  }

  return {
    addInvitation,
    handleNewConnection,
    handleProofPaused,
    handleProofDone
  }
}