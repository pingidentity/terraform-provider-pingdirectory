terraform {
  required_providers {
    pingdirectory = {
      source = "pingidentity.com/terraform/pingdirectory"
    }
  }
}

provider "pingdirectory" {
  username = "cn=administrator"
  password = "2FederateM0re"
  ldap_host = "ldap://localhost:1389"
  https_host = "https://localhost:1443"
  default_user_password = "2FederateM0re"
}

resource "pingdirectory_location" "drangleic" {
  name = "Drangleic"
  description = "Seek the king"
}
