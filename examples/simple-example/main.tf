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

resource "pingdirectory_user" "mahomes" {
  uid = "pm"
  sn = "Mahomes"
  given_name = "Patrick"
  mail = "pmbro@kcchiefs.com"
}

resource "pingdirectory_user" "knight" {
  uid = "hk"
  description = "the knight"
  sn = "Knight"
  given_name = "Hollow"
  mail = "hk@hallownest.com"
}

output "mahomes_user" {
  value = pingdirectory_user.mahomes.cn
}

output "knight_user" {
  value = pingdirectory_user.knight.cn
}
