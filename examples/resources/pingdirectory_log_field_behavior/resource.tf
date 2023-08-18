resource "pingdirectory_log_field_behavior" "myLogFieldBehavior" {
  name        = "MyLogFieldBehavior"
  type        = "text-access"
  description = "My text access log field behavior"
}
