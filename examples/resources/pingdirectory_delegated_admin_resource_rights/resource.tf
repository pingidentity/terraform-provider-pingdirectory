terraform {
  required_version = ">=1.1"
  required_providers {
    pingdirectory = {
      version = "~> 0.3.0"
      source  = "pingidentity/pingdirectory"
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
  product_version        = "9.3.0.0"
}

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
