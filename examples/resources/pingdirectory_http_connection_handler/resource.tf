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

resource "pingdirectory_http_connection_handler" "http" {
  id                     = "example"
  description            = "Description of http connection handler"
  listen_port            = 2443
  enabled                = true
  http_servlet_extension = ["Available or Degraded State", "Available State"]
}
