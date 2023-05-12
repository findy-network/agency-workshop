package fi.oplab.findyagency.workshop

import org.findy_network.findy_common_kt.*

class Pairwise(id: String) {
  val id: String = id
}

class Issuer(
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
      // Connection was not for issuing, skip
      return
    }

    val attrs = mapOf("foo" to "bar")

    println("Offer credential, conn id: ${notification.connectionID}, credDefID: ${credDefId}, attrs: ${attrs}")

    // Send credential offer to the other agent
    connection.protocolClient.sendCredentialOffer(
      notification.connectionID,
      attrs,
      credDefId
    )
  }

  override suspend fun handleIssueCredentialDone(
    notification: Notification,
    status: ProtocolStatus.IssueCredentialStatus
  ) {
    println("Credential issued, conn id: ${notification.connectionID} with id ${notification.protocolID}")

    pwConnections.remove(notification.connectionID)
  }
}
