data "pingdirectory_correlated_ldap_data_view" "myCorrelatedLdapDataView" {
  name                    = "MyCorrelatedLdapDataView"
  scim_resource_type_name = "MyScimResourceType"
}
