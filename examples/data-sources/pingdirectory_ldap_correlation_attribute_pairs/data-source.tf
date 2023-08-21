data "pingdirectory_ldap_correlation_attribute_pairs" "list" {
  correlated_ldap_data_view_name = "MyCorrelatedLdapDataView"
  scim_resource_type_name        = "MyScimResourceType"
}
