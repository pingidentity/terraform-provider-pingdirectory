resource "pingdirectory_external_server" "myLdapExternalServer" {
  type                  = "ldap"
  name                  = "MyLdapExternalServer"
  server_host_name      = "example.com"
  authentication_method = "none"
}

resource "pingdirectory_pass_through_authentication_handler" "myPassThroughAuthenticationHandler" {
  name   = "MyPassThroughAuthenticationHandler"
  type   = "ldap"
  server = [pingdirectory_ldap_external_server.myLdapExternalServer.id]
}
