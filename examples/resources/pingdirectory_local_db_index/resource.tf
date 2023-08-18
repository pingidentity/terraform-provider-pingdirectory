resource "pingdirectory_local_db_index" "myLocalDbIndex" {
  backend_name = "userRoot"
  attribute    = "dc"
  index_type   = ["equality"]
}
