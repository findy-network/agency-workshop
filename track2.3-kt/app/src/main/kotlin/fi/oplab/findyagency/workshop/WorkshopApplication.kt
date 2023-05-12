package fi.oplab.findyagency.workshop

import kotlinx.serialization.decodeFromString
import org.findy_network.findy_common_kt.*
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RestController

@SpringBootApplication class WorkshopApplication

fun main(args: Array<String>) {
  runApplication<WorkshopApplication>(*args)
}

@kotlinx.serialization.Serializable data class InvitationData(
  @kotlinx.serialization.SerialName("@id") val id: String
)

@RestController
class AppController {
  val agent = Agent()
  val greeter = Greeter(agent.connection)
  // Create issuer instance
  val issuer = Issuer(agent.connection, agent.credDefId)

  init {
    val listeners = ArrayList<Listener>()
    listeners.add(greeter)
    // Add issuer to the listener array
    listeners.add(issuer)
    // Start listening to agent notifications 
    agent.listen(listeners)
  }

  fun createInvitationPage(header: String): Pair<String, String> {

    // Create invitation to our agent
    val invitation = createInvitation()

    // Parse invitation id
    val data = kotlinx.serialization.json.Json { ignoreUnknownKeys = true }
      .decodeFromString<InvitationData>(invitation.getJSON())
    println("Created invitation with id ${data.id}")

    val imgSrc = "data:image/png;base64," + createQRCode(invitation.url)
    val html = """<html>
    <h1>${header}</h1>
    <p>Read the QR code with the wallet application:</p>
    <img src="${imgSrc}"/>
    <p>or copy-paste the invitation:</p>
    <textarea onclick="this.focus();this.select()" readonly="readonly" rows="10" cols="60">${invitation.url}</textarea></html>"""

    // return both html page and invitation id
    return Pair(html, data.id)
  }

  fun createInvitation(): Invitation = kotlinx.coroutines.runBlocking {
    // Use as label whichever name we want to expose from ourselves to the other end
    agent.connection.agentClient.createInvitation(label = System.getenv("FCLI_USER"))
  }

  // Utility for converting string to QR code
  fun createQRCode(value: String): String {
    val writer = com.google.zxing.qrcode.QRCodeWriter()
    val bitMatrix = writer.encode(value, com.google.zxing.BarcodeFormat.QR_CODE, 512, 512)
    val width = bitMatrix.width
    val height = bitMatrix.height
    val bitmap = java.awt.image.BufferedImage(
      width,
      height,
      java.awt.image.BufferedImage.TYPE_USHORT_565_RGB
    )
    for (x in 0 until width) {
      for (y in 0 until height) {
        bitmap.setRGB(
          x,
          y,
          if (bitMatrix.get(x, y)) java.awt.Color.BLACK.getRGB() else java.awt.Color.WHITE.getRGB()
        )
      }
    }
    val out = java.io.ByteArrayOutputStream()
    javax.imageio.ImageIO.write(bitmap, "PNG", out)
    return java.util.Base64.getEncoder().encodeToString(out.toByteArray())
  }


  @GetMapping("/") fun index(): String = "Kotlin example"

  @GetMapping("/greet") fun greet(): String {
    val (html) = createInvitationPage("Greet")
    return html
  }

  // Show pairwise invitation. Once connection is established, issue credential.
  @GetMapping("/issue") fun issue(): String {
    val (html, id) = createInvitationPage("Issue")
    issuer.addInvitation(id)
    return html
  }

  @GetMapping("/verify") fun verify(): String = "IMPLEMENT ME"
}
