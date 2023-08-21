data "pingdirectory_local_db_index" "myLocalDbIndex" {
  backend_name = "MyBackend"
  attribute    = "myattr"
}
