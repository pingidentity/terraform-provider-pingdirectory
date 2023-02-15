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

resource "pingdirectory_indicator_gauge" "myIndicatorGauge" {
  id                = "MyIndicatorGauge"
  gauge_data_source = "Replication Connection Status"
  enabled           = false
}
