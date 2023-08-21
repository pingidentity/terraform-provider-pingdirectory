resource "pingdirectory_log_field_syntax" "myLogFieldSyntax" {
  name             = "MyLogFieldSyntax"
  type             = "json"
  default_behavior = "omit"
}
