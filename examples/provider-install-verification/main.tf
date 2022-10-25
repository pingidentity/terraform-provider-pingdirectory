terraform {
  required_providers {
    pingdirectory = {
      source = "pingidentity.com/terraform/pingdirectory"
    }
  }
}

provider "pingdirectory" {
  username = "cn=Directory Manager"
  password = "2FederateM0re"
  host     = "ldap://localhost:1389"
}

resource "pingdirectory_user" "myuser" {
  dn = "uid=myuser,ou=people,dc=example,dc=com"
  description = "myterraformuser"
}

output "myuser_user" {
  value = pingdirectory_user.myuser
}
