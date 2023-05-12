package fi.oplab.findyagency.workshop

import org.findy_network.findy_common_kt.*

class Verifier(
  connection: Connection,
  credDefId: String
) : Listener {
  val connection = connection
  val pwConnections: MutableMap<String, Pairwise> =
    java.util.Collections.synchronizedMap(mutableMapOf<String, Pairwise>())
  val credDefId: String = credDefId

  fun addInvitation(id: String) {
    pwConnections.put(id, Pairwise(id = id))
  }

  override suspend fun handleNewConnection(
    notification: Notification,
    status: ProtocolStatus.DIDExchangeStatus
  ) {
    if (!pwConnections.contains(notification.connectionID)) {
      // Connection was not for verifying, skip
      return
    }

    val attrs = listOf(ProofRequestAttribute("foo", credDefId))

    println("Request proof, conn id: ${notification.connectionID}, credDefID: ${credDefId}, attrs: ${attrs}")

    // Send credential offer to the other agent
    connection.protocolClient.sendProofRequest(
        notification.connectionID,
        attrs,
    )
  }

  // This function is called after proof is verified cryptographically.
  // The application can execute its business logic and reject the proof
  // if the attribute values are not valid.
  override suspend fun handlePresentProofPaused(
    notification: Notification,
    status: ProtocolStatus.PresentProofStatus
  ) {
    // we have no special logic here - accept all received values
    connection.protocolClient.resumeProofRequest(notification.protocolID, true)
    println("Proof continued with id ${notification.protocolID}")
  }

  override suspend fun handlePresentProofDone(
    notification: Notification,
    status: ProtocolStatus.PresentProofStatus
  ) {
    if (!pwConnections.contains(notification.connectionID)) {
      // Connection was not for verifying, skip
      return
    }
    println("Proof verified from ${notification.connectionID} with id ${notification.protocolID}")

    pwConnections.remove(notification.connectionID)
  }
}
