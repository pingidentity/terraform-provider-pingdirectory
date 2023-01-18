terraform {
  required_providers {
    pingdirectory = {
      source = "pingidentity/pingdirectory"
    }
  }
}

provider "pingdirectory" {
  username = "cn=administrator"
  password = "2FederateM0re"
  https_host = "https://localhost:1443"
}

resource "pingdirectory_location" "drangleic" {
  name = "Drangleic"
  description = "Seek the king"
}
