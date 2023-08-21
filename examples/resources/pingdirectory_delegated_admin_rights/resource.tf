resource "pingdirectory_delegated_admin_rights" "myDelegatedAdminRights" {
  name          = "MyDelegatedAdminRights"
  enabled       = true
  admin_user_dn = "cn=admin-users,dc=test,dc=com"
}
