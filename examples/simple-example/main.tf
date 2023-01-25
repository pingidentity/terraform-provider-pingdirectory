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
  id = "Drangleic"
  description = "Seek the king"
}

resource "pingdirectory_global_configuration" "global" {
  location = "Docker"
  encrypt_data = false
  sensitive_attribute = ["Delivered One-Time Password", "TOTP Shared Secret"]
  tracked_application = ["Requests by Root Users"]
  result_code_map = "Sun DS Compatible Behavior"
  disabled_privilege = ["jmx-write", "jmx-read"]
}

resource "pingdirectory_blind_trust_manager_provider" "blindtest" {
  id = "Blind Test"
  enabled = true
  include_jvm_default_issuers = true
}

resource "pingdirectory_file_based_trust_manager_provider" "filetest" {
  id = "FileTest"
  enabled = true
  trust_store_file = "config/keystore"
  trust_store_type = "pkcs12"
  include_jvm_default_issuers = true
}

resource "pingdirectory_jvm_default_trust_manager_provider" "jvmtest" {
  id = "jvmtest"
  enabled = false
}

resource "pingdirectory_third_party_trust_manager_provider" "tptest" {
  id = "tptest"
  enabled = false
  extension_class = "com.unboundid.directory.sdk.common.api.TrustManagerProvider"
  extension_argument = ["val1=one", "val2=two"]
}
