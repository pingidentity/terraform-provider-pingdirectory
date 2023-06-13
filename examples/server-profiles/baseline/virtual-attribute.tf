resource "pingdirectory_default_virtual_attribute" "defaultPwpStateJsonVirtualAttribute" {
  type    = "password-policy-state-json"
  id      = "Password Policy State JSON"
  enabled = true
}
