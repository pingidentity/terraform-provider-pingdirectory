resource "pingdirectory_scim_attribute_mapping" "myScimAttributeMapping" {
  name                         = "MyScimAttributeMapping"
  scim_resource_type_name      = "MyScimResourceType"
  scim_resource_type_attribute = "name"
  ldap_attribute               = "name"
}
