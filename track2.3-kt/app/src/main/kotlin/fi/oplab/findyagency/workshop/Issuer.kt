package fi.oplab.findyagency.workshop

import com.sendgrid.*;
import com.sendgrid.helpers.mail.*;
import com.sendgrid.helpers.mail.objects.*;
import org.findy_network.findy_common_kt.*

class Pairwise(id: String) {
  val id: String = id
  var email: String? = null
  var verified: Boolean = false
}

class Issuer(
  connection: Connection,
  credDefId: String
) : Listener {
  val connection = connection
  val pwConnections: MutableMap<String, Pairwise> =
    java.util.Collections.synchronizedMap(mutableMapOf<String, Pairwise>())
  val credDefId: String = credDefId
  val sgClient: SendGrid

  init {
    sgClient = SendGrid(System.getenv("SENDGRID_API_KEY"));
  }

  private suspend fun askForEmail(connectionId: String) {
    // Ask for user email via basic message
    connection.protocolClient.sendMessage(
      connectionId,
      "Please enter your email to get started."
    )
  }

  // Send email via SendGrid API
  private suspend fun sendEmail(email: String, content: String) {
    val mail = Mail(
      Email(System.getenv("SENDGRID_SENDER")),
      "Email verification",
      Email(email),
      Content("text/plain", content)
    )

    println("Sending email ${content} to ${email}")

    val request = Request()
    request.setMethod(Method.POST)
    request.setEndpoint("mail/send")
    request.setBody(mail.build())
    sgClient.api(request)
  }

  fun addInvitation(id: String) {
    pwConnections.put(id, Pairwise(id = id))
  }

  suspend fun setEmailVerified(connectionId: String): Boolean {
    // Skip handling if the verification is already done
    val pwConnection = pwConnections.get(connectionId)
    if (pwConnection == null || pwConnection.email == null || pwConnection.verified) {
      return false
    }

    pwConnection.verified = true

    val attrs = mapOf("email" to pwConnection.email!!)

    println("Offer credential, conn id: ${connectionId}, credDefID: ${credDefId}, attrs: ${attrs}")

    // Send credential offer to the other agent
    connection.protocolClient.sendCredentialOffer(
      pwConnection.id,
      attrs,
      credDefId
    )
    return true
  }

  override suspend fun handleNewConnection(
    notification: Notification,
    status: ProtocolStatus.DIDExchangeStatus
  ) {
    val pwConnection = pwConnections.get(notification.connectionID)
    if (pwConnection == null) {
      // Connection was not for issuing, skip
      return
    }

    if (pwConnection.email == null) {
      askForEmail(notification.connectionID)
    }
  }

  override suspend fun handleBasicMessageDone(
    notification: Notification,
    status: ProtocolStatus.BasicMessageStatus
  ) {
    // Skip handling if message was sent by us or
    // the verification is already done
    val pwConnection = pwConnections.get(notification.connectionID)
    if (status.getSentByMe() || pwConnection == null || pwConnection.email != null) {
      return
    }

    println("Basic message ${status.getContent()} received with id ${notification.protocolID}")

    // Some sanity checking
    val email = status.getContent()
    val emailValid = email.split(' ').size == 1 && email.indexOf('@') >= 0

    if (emailValid) {
      // Valid email, do verification
      pwConnection.email = email

      val content = "Please verify your email by clicking the following link:\n http://localhost:3001/email/${pwConnection.id}"
      // Send verification mail
      sendEmail(pwConnection.email!!, content)

      // Send confirmation via basic message
      connection.protocolClient.sendMessage(
        pwConnection.id,
        "Email is on it's way! Please check your mailbox 📫."
      )
    } else {
      // Email invalid, ask again
      askForEmail(pwConnection.id)
    }
  }

  override suspend fun handleIssueCredentialDone(
    notification: Notification,
    status: ProtocolStatus.IssueCredentialStatus
  ) {
    println("Credential issued, conn id: ${notification.connectionID} with id ${notification.protocolID}")

    pwConnections.remove(notification.connectionID)
  }
}
