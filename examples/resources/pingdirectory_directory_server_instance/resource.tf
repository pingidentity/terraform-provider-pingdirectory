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

resource "pingdirectory_directory_server_instance" "mine" {
  //NOTE This id needs to match the instance name generated for the running instance
  id                            = "instanceName"
  jmx_port                      = 1112
  start_tls_enabled             = true
  load_balancing_algorithm_name = []
}