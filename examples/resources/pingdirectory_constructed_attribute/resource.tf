resource "pingdirectory_constructed_attribute" "myConstructedAttribute" {
  name           = "MyConstructedAttribute"
  attribute_type = "cn"
  value_pattern  = ["{givenName} {sn}"]
}
