terraform {
  required_providers {
    pingdirectory = {
      source = "pingidentity/pingdirectory"
    }
  }
}

provider "pingdirectory" {
  username   = "cn=administrator"
  password   = "2FederateM0re"
  https_host = "https://localhost:1443"
}

resource "pingdirectory_simple_connection_criteria" "mySimpleConnectionCriteria" {
  id             = "MySimpleConnectionCriteria"
  description    = "Simple connection example"
  user_auth_type = ["internal", "sasl"]
}
