resource "pingdirectory_scim_attribute" "myScimAttribute" {
  scim_schema_name = "urn:com:example"
  name             = "cn"
}
