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

# If these are uncommented, the "Sensitive Password Attributes" sensitive_attribute value needs to be removed from the global config below
#resource "pingdirectory_user" "mahomes" {
#  uid = "pm"
#  sn = "Mahomes"
#  given_name = "Patrick"
#  mail = "pm@kcchiefs.com"
#}
#
#resource "pingdirectory_user" "knight" {
#  uid = "hk"
#  description = "the knight"
#  sn = "Knight"
#  given_name = "Hollow"
#  mail = "hk@hallownest.com"
#}

resource "pingdirectory_location" "drangleic" {
  name = "Drangleic"
  description = "Seek the king"
}

resource "pingdirectory_global_configuration" "global" {
  location = "Docker"
  encrypt_data = true
  sensitive_attribute = ["Delivered One-Time Password", "TOTP Shared Secret", "Sensitive Password Attributes"]
  tracked_application = ["Requests by Root Users"]
  result_code_map = "Sun DS Compatible Behavior"
  #result_code_map = ""
}
