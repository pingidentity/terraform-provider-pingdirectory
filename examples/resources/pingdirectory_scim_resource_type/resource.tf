resource "pingdirectory_scim_resource_type" "myScimResourceType" {
  name     = "MyScimResourceType"
  type     = "ldap-pass-through"
  enabled  = false
  endpoint = "myendpoint"
}
