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

resource "pingdirectory_local_db_backend" "myLocalDbBackend" {
  id                    = "MyLocalDbBackend"
  backend_id            = "MyLocalDbBackend"
  base_dn               = ["dc=example1,dc=com"]
  writability_mode      = "enabled"
  db_directory          = "db"
  import_temp_directory = "tmp"
  enabled               = true
}
