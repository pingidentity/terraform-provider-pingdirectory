resource "pingdirectory_connection_handler" "myConnectionHandler" {
  name        = "MyConnectionHandler"
  type        = "jmx"
  listen_port = 1234
  enabled     = false
}
