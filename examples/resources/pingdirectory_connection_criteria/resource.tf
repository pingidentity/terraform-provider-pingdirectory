resource "pingdirectory_connection_criteria" "myConnectionCriteria" {
  name           = "MyConnectionCriteria"
  type           = "simple"
  description    = "Simple connection example"
  user_auth_type = ["internal", "sasl"]
}
