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

resource "pingdirectory_user" "mahomes" {
  uid = "pm"
  sn = "Mahomes"
  given_name = "Patrick"
  mail = "pm@kcchiefs.com"
}

resource "pingdirectory_user" "knight" {
  uid = "hk"
  description = "the knight"
  sn = "Knight"
  given_name = "Hollow"
  mail = "hk@hallownest.com"
}

resource "pingdirectory_location" "drangleic" {
  name = "Drangleic"
  description = "Seek the king"
}

output "mahomes_user" {
  value = pingdirectory_user.mahomes.cn
}

output "knight_user" {
  value = pingdirectory_user.knight.cn
}

output "drangleic_location" {
  value = pingdirectory_location.drangleic
}
