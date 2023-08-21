resource "pingdirectory_virtual_attribute" "myVirtualAttribute" {
  name             = "MyVirtualAttribute"
  type             = "mirror"
  source_attribute = "mail"
  enabled          = true
  attribute_type   = "name"
}
