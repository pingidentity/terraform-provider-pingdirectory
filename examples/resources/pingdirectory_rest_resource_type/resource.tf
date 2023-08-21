resource "pingdirectory_rest_resource_type" "myRestResourceType" {
  name                        = "MyRestResourceType"
  type                        = "user"
  enabled                     = true
  resource_endpoint           = "userRestResource"
  structural_ldap_objectclass = "inetOrgPerson"
  search_base_dn              = "cn=users,dc=test,dc=com"
}
