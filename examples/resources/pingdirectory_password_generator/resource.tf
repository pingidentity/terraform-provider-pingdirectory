resource "pingdirectory_password_generator" "myPasswordGenerator" {
  name                   = "MyPasswordGenerator"
  type                   = "random"
  password_character_set = ["set:abcdefghijklmnopqrstuvwxyz"]
  password_format        = "set:15"
  enabled                = true
}
