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

resource "pingdirectory_debug_target" "myDebugTarget" {
  id                 = "com.example.MyClass"
  log_publisher_name = "File-Based Debug Logger"
  debug_scope        = "com.example.MyClass"
  debug_level        = "all"
}
