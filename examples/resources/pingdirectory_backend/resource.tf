resource "pingdirectory_backend" "myBackend" {
  type                  = "local-db"
  base_dn               = ["dc=example1,dc=com"]
  backend_id            = "myId"
  writability_mode      = "enabled"
  db_directory          = "db"
  import_temp_directory = "tmp"
  enabled               = true
}
