# Use "pingdirectory_default_delegated_admin_rights" if you are adopting existing configuration from the PingDirectory server into Terraform
resource "pingdirectory_delegated_admin_rights" "myDelegatedAdminRights" {
  name          = "MyDelegatedAdminRights"
  enabled       = true
  admin_user_dn = "cn=admin-users,dc=test,dc=com"
}

# Use "pingdirectory_default_rest_resource_type" if you are adopting existing configuration from the PingDirectory server into Terraform
resource "pingdirectory_rest_resource_type" "myUserRestResourceType" {
  type                        = "user"
  name                        = "MyUserRestResourceType"
  enabled                     = true
  resource_endpoint           = "userRestResource"
  structural_ldap_objectclass = "inetOrgPerson"
  search_base_dn              = "cn=users,dc=test,dc=com"
}

# Use "pingdirectory_default_delegated_admin_resource_rights" if you are adopting existing configuration from the PingDirectory server into Terraform
resource "pingdirectory_delegated_admin_resource_rights" "myDelegatedAdminResourceRights" {
  delegated_admin_rights_name = pingdirectory_delegated_admin_rights.myDelegatedAdminRights.id
  enabled                     = true
  admin_permission            = ["create", "read"]
  rest_resource_type          = pingdirectory_rest_resource_type.myUserRestResourceType.id
}
