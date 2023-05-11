package fi.oplab.findyagency.workshop

import org.findy_network.findy_common_kt.*

class Greeter(connection: Connection) : Listener {
  val connection = connection

  override suspend fun handleNewConnection(
    notification: Notification,
    status: ProtocolStatus.DIDExchangeStatus
  ) {
    println("New connection ${status.theirLabel} with id ${notification.connectionID}")
  }
}