resource "pingdirectory_delegated_admin_correlated_rest_resource" "myDelegatedAdminCorrelatedRestResource" {
  name                                          = "MyDelegatedAdminCorrelatedRestResource"
  rest_resource_type_name                       = "MyRestResourceType"
  display_name                                  = "MyDisplayName"
  correlated_rest_resource                      = "MyRestResourceType"
  primary_rest_resource_correlation_attribute   = "cn"
  secondary_rest_resource_correlation_attribute = "sn"
}
