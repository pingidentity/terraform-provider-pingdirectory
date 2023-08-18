data "pingdirectory_token_claim_validation" "myTokenClaimValidation" {
  name                    = "MyTokenClaimValidation"
  id_token_validator_name = "MyIdTokenValidator"
}
