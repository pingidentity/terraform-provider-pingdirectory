resource "pingdirectory_local_db_vlv_index" "myLocalDbVlvIndex" {
  backend_name = "userRoot"
  base_dn      = ["dc=example,dc=com"]
  scope        = "base-object"
  filter       = "uid=user.1"
  sort_order   = "givenName"
  name         = "my_example"
}
