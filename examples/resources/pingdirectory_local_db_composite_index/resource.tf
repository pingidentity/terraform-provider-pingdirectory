resource "pingdirectory_local_db_composite_index" "myLocalDbCompositeIndex" {
  name                 = "MyLocalDbCompositeIndex"
  backend_name         = "userRoot"
  description          = "My local DB composite index"
  index_filter_pattern = "(sn=?)"
}
