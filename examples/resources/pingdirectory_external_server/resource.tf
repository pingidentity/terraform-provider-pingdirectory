resource "pingdirectory_external_server" "myExternalServer" {
  name             = "MyExternalServer"
  type             = "smtp"
  server_host_name = "mysmtp.mailserver.com"
  server_port      = 25
}
