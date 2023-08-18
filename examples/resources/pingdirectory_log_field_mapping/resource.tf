resource "pingdirectory_log_field_mapping" "myLogFieldMapping" {
  name            = "MyLogFieldMapping"
  type            = "access"
  description     = "My access log field mapping"
  log_field_scope = "search_scope"
}
