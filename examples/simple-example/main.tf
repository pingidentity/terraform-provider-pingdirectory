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
  host     = "ldap://localhost:1389"
}

resource "pingdirectory_user" "myuser" {
  uid = "myuid"
  description = "myterraformuser"
}

output "myuser_user" {
  value = pingdirectory_user.myuser
}
