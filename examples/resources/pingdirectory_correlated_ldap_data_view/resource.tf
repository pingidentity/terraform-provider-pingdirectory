//URN is needed for the LDAP mapping SCIM resource type
resource "pingdirectory_scim_schema" "myScimSchema" {
  schema_urn = "urn:com:example2"
}

// LDAP mapping SCIM resource type is needed for the correlated data view resource
resource "pingdirectory_scim_resource_type" "myLdapMappingScimResourceType" {
  name        = "MyLdapMappingScimResourceType2"
  type        = "ldap-mapping"
  core_schema = pingdirectory_scim_schema.myScimSchema.schema_urn
  enabled     = false
  endpoint    = "myendpoint"
}

resource "pingdirectory_correlated_ldap_data_view" "myCorrelatedLdapDataView" {
  name                            = "MyCorrelatedLdapDataView"
  scim_resource_type_name         = pingdirectory_scim_resource_type.myLdapMappingScimResourceType.id
  structural_ldap_objectclass     = "ldapObject"
  include_base_dn                 = "cn=com.example"
  primary_correlation_attribute   = "cn"
  secondary_correlation_attribute = "cn"
}
