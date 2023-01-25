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

resource "pingdirectory_directory_server_instance" "mine" {
  //NOTE This instance name needs to match the instance name generated for the running instance
  id = "d1db4a163621"
  server_instance_name = "d1db4a163621"
  jmx_port = 1112
  start_tls_enabled = true
  load_balancing_algorithm_name = []
}