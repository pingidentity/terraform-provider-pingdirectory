resource "pingdirectory_password_validator" "myPasswordValidator" {
  name                = "MyPasswordValidator"
  type                = "length-based"
  min_password_length = 8
  max_password_length = 100
  enabled             = true
}
