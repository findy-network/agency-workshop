package fi.oplab.findyagency.workshop

import org.findy_network.findy_common_kt.*

class Greeter(connection: Connection) : Listener {
  val connection = connection

  override suspend fun handleNewConnection(
    notification: Notification,
    status: ProtocolStatus.DIDExchangeStatus
  ) {
    println("New connection ${status.theirLabel} with id ${notification.connectionID}")

    // Greet each new connection with basic message
    connection.protocolClient.sendMessage(
      notification.connectionID,
      "Hi there 👋!"
    )
  }

  override suspend fun handleBasicMessageDone(
    notification: Notification,
    status: ProtocolStatus.BasicMessageStatus
  ) {

    if (!status.sentByMe) {
      println("Received basic message ${status.content} from ${notification.connectionID}")
    }
  }
}