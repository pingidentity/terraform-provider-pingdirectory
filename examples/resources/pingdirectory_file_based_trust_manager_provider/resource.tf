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

resource "pingdirectory_file_based_trust_manager_provider" "filetest" {
  id                          = "FileTest"
  enabled                     = true
  trust_store_file            = "config/keystore"
  trust_store_type            = "pkcs12"
  include_jvm_default_issuers = true
}
