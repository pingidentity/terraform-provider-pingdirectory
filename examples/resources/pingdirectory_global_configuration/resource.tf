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

resource "pingdirectory_global_configuration" "global" {
  location            = "Docker"
  encrypt_data        = true
  sensitive_attribute = ["Delivered One-Time Password", "TOTP Shared Secret"]
  tracked_application = ["Requests by Root Users"]
  result_code_map     = "Sun DS Compatible Behavior"
  disabled_privilege  = ["jmx-write", "jmx-read"]
}
