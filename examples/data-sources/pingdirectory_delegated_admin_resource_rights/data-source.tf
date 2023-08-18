data "pingdirectory_delegated_admin_resource_rights" "myDelegatedAdminResourceRights" {
  delegated_admin_rights_name = "MyDelegatedAdminRights"
  rest_resource_type          = "myRestResourceType"
}
