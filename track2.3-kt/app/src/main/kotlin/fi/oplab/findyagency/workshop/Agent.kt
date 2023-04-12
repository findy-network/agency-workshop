package fi.oplab.findyagency.workshop

import org.findy_network.findy_common_kt.*

class Agent {
  public val connection: Connection = Connection(
    authOrigin = System.getenv("FCLI_ORIGIN"),
    authUrl = System.getenv("FCLI_URL"),
    certFolderPath = System.getenv("FCLI_TLS_PATH"),
    key = System.getenv("FCLI_KEY"),
    port = Integer.parseInt(System.getenv("AGENCY_API_SERVER_PORT")),
    seed = "",
    server = System.getenv("AGENCY_API_SERVER"),
    userName = System.getenv("FCLI_USER"),
  )
}