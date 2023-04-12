package fi.oplab.findyagency.workshop

import org.findy_network.findy_common_kt.*

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
}