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

resource "pingdirectory_exact_match_identity_mapper" "myExactMatchIdentityMapper" {
  id              = "MyExactMatchIdentityMapper"
  match_attribute = ["uid"]
  enabled         = true
}
