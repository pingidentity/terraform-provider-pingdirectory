data "pingdirectory_local_db_composite_index" "myLocalDbCompositeIndex" {
  name         = "MyLocalDbCompositeIndex"
  backend_name = "MyBackend"
}
