resource "pingdirectory_default_attribute_syntax" "myAttributeSyntax" {
  name                    = "MyAttributeSyntax"
  enabled                 = false
  require_binary_transfer = true
}
