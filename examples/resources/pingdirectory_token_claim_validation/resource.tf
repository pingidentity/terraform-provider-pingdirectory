resource "pingdirectory_id_token_validator" "myPingOneIdTokenValidator" {
  type                   = "ping-one"
  name                   = "MyPingOneIdTokenValidator"
  issuer_url             = "example.com"
  enabled                = false
  identity_mapper        = "Exact Match"
  evaluation_order_index = 1
}

resource "pingdirectory_token_claim_validation" "myTokenClaimValidation" {
  name                    = "MyTokenClaimValidation"
  id_token_validator_name = pingdirectory_id_token_validator.myPingOneIdTokenValidator.id
  any_required_value      = ["my_example_value"]
  claim_name              = "my_example_claim_name"
}
