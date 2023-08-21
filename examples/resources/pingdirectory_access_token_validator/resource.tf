resource "pingdirectory_access_token_validator" "myAccessTokenValidator" {
  name                 = "MyAccessTokenValidator"
  type                 = "ping-federate"
  client_id            = "my-client-id"
  enabled              = false
  client_secret        = "my-client-secrets"
  authorization_server = "PingOne Auth Service"
}
