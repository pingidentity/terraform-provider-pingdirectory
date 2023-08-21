data "pingdirectory_delegated_admin_attribute" "myDelegatedAdminAttribute" {
  rest_resource_type_name = "MyRestResourceType"
  attribute_type          = "myattr"
}
