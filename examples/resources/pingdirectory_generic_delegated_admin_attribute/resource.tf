terraform {
  required_version = ">=1.1"
  required_providers {
    pingdirectory = {
      version = "~> 0.3.0"
      source = "pingidentity/pingdirectory"
    }
  }
}

provider "pingdirectory" {
  username   = "cn=administrator"
  password   = "2FederateM0re"
  https_host = "https://localhost:1443"
  # Warning: The insecure_trust_all_tls attribute configures the provider to trust any certificate presented by the PingDirectory server.
  # It should not be used in production. If you need to specify trusted CA certificates, use the
  # ca_certificate_pem_files attribute to point to any number of trusted CA certificate files
  # in PEM format. If you do not specify certificates, the host's default root CA set will be used.
  # Example:
  # ca_certificate_pem_files = ["/example/path/to/cacert1.pem", "/example/path/to/cacert2.pem"]
  insecure_trust_all_tls = true
  product_version        = "9.2.0.0"
}

# Use "pingdirectory_default_generic_delegated_admin_attribute" if you are adopting existing configuration from the PingDirectory server into Terraform
resource "pingdirectory_generic_delegated_admin_attribute" "myDelegatedAdminAttributeDevice" {
  rest_resource_type_name = pingdirectory_generic_rest_resource_type.myRestResourceTypeDevice.id
  attribute_type          = "cn"
  display_name            = "Device Name"
  display_order_index     = 1
}

# Use "pingdirectory_default_generic_delegated_admin_attribute" if you are adopting existing configuration from the PingDirectory server into Terraform
resource "pingdirectory_generic_delegated_admin_attribute" "myDelegatedAdminAttributeSerialNumber" {
  rest_resource_type_name = pingdirectory_generic_rest_resource_type.myRestResourceTypeDevice.id
  attribute_type          = "serialNumber"
  display_name            = "Serial Number"
  display_order_index     = 2
}

# Use "pingdirectory_default_generic_rest_resource_type" if you are adopting existing configuration from the PingDirectory server into Terraform
resource "pingdirectory_generic_rest_resource_type" "myRestResourceTypeDevice" {
  id                             = "device"
  enabled                        = true
  resource_endpoint              = "device"
  display_name                   = "Device"
  structural_ldap_objectclass    = "device"
  search_base_dn                 = "dc=example,dc=com"
  parent_dn                      = "dc=example,dc=com"
  search_filter_pattern          = "(cn=*%%*)"
  primary_display_attribute_type = "cn"
}
