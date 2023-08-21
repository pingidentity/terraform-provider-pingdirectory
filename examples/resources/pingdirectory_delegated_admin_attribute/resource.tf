resource "pingdirectory_delegated_admin_attribute" "myDelegatedAdminAttribute" {
  rest_resource_type_name = "MyRestResourceType"
  type                    = "certificate"
  attribute_type          = "myattr"
  display_name            = "MyAttribute"
}
