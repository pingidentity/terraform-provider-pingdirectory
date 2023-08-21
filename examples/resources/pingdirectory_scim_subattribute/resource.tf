resource "pingdirectory_scim_schema" "myScimSchema" {
  schema_urn = "urn:com:example"
}

resource "pingdirectory_scim_attribute" "myScimAttribute" {
  scim_schema_name = pingdirectory_scim_schema.myScimSchema.schema_urn
  name             = "cn"
}

resource "pingdirectory_scim_subattribute" "myScimSubattribute" {
  name                = "MyScimSubattribute"
  scim_attribute_name = pingdirectory_scim_attribute.myScimAttribute.name
  scim_schema_name    = pingdirectory_scim_schema.myScimSchema.schema_urn
}
