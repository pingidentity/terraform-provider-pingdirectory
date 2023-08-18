resource "pingdirectory_id_token_validator" "myIdTokenValidator" {
  name                   = "MyIdTokenValidator"
  type                   = "ping-one"
  issuer_url             = "example.com"
  enabled                = false
  identity_mapper        = "Exact Match"
  evaluation_order_index = 1
}
