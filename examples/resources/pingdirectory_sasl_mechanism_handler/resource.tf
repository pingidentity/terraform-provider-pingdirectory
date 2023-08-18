resource "pingdirectory_sasl_mechanism_handler" "mySaslMechanismHandler" {
  name            = "MySaslMechanismHandler"
  type            = "unboundid-ms-chap-v2"
  identity_mapper = "Exact Match"
  enabled         = false
}
