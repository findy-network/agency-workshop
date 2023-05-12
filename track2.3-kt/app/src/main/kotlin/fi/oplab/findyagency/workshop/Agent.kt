package fi.oplab.findyagency.workshop

import org.findy_network.findy_common_kt.*
import kotlinx.coroutines.launch

interface Listener {
  suspend fun handleNewConnection(
    notification: Notification,
    status: ProtocolStatus.DIDExchangeStatus
  ) {}

  // Send notification to listener when basic message protocol is completed
  suspend fun handleBasicMessageDone(
    notification: Notification,
    status: ProtocolStatus.BasicMessageStatus
  ) {}
}

class Agent {
  public val connection: Connection = Connection(
    authOrigin = System.getenv("FCLI_ORIGIN"),
    authUrl = System.getenv("FCLI_URL"),
    // NOTE: we currently assume that we do not need certs for cloud installation
    // as the cert is issued by a trusted issuer
    certFolderPath = if (System.getenv("AGENCY_API_SERVER") == "localhost") System.getenv("FCLI_TLS_PATH") else null,
    key = System.getenv("FCLI_KEY"),
    port = Integer.parseInt(System.getenv("AGENCY_API_SERVER_PORT")),
    seed = "",
    server = System.getenv("AGENCY_API_SERVER"),
    userName = System.getenv("FCLI_USER"),
  )
  public val credDefId = createCredentialDefinition()

  fun listen(listeners: List<Listener>) {
    kotlinx.coroutines.GlobalScope.launch {
      connection.agentClient.listen().collect {
        println("Received from Agency:\n$it")
        val status = it.notification
        when (status.typeID) {
          Notification.Type.STATUS_UPDATE -> {
            // info contains the protocol related information
            val info = connection.protocolClient.status(status.protocolID)
            val getType =
                fun(): Protocol.Type =
                    if (info.state.state == ProtocolState.State.OK) status.protocolType
                    else Protocol.Type.NONE

            when (getType()) {
              // New connection established
              Protocol.Type.DIDEXCHANGE -> {
                listeners.map{ it.handleNewConnection(status, info.didExchange) }
              }
              // Notify basic message protocol events
              Protocol.Type.BASIC_MESSAGE -> {
                listeners.map{ it.handleBasicMessageDone(status, info.basicMessage) }
              }
              else -> println("no handler for protocol type: ${status.protocolType}")
            }
          }
          else -> println("no handler for notification type: ${status.typeID}")
        }
      }
    }
  }


  private fun createCredentialDefinition(): String = kotlinx.coroutines.runBlocking {
    var credDefId = ""
    try {
      credDefId = java.io.File("CRED_DEF_ID").readLines()[0]
    } catch (e: Exception) {}

    // Create cred def only if it does not exist already
    if (credDefId == "") {
      // Note: if schema creation fails, you have probably created
      // the same schema already with the same name and version number.
      // If this happens, change the version number and retry the creation
      val schemaRes =
          connection.agentClient.createSchema(
              name = "foobar",
              attributes = listOf("foo"),
              version = "1.0"
          )
      do {
        var schemaCreated = false
        try {
          val schema = connection.agentClient.getSchema(id = schemaRes.id)
          println("Created schema ${schema.id}")
          schemaCreated = true
        } catch (e: Exception) {
          println("Waiting for the schema to be created...")
          Thread.sleep(1_000 * 1)
        }
      } while (!schemaCreated)

      println("Starting to create credential definition (may take a while)...")
      val credDef = connection.agentClient.createCredDef(
        schemaId = schemaRes.id,
        tag = System.getenv("FCLI_USER")
      )
      credDefId = credDef.id
      java.io.File("CRED_DEF_ID").writeText(credDefId)
      println("Cred def ${credDefId} created successfully!")
    }
    credDefId
  }
}