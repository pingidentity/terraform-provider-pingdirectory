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
  product_version = "9.2.0.0"
}

//URN is needed for the LDAP mapping SCIM resource type
resource "pingdirectory_scim_schema" "myScimSchema" {
  schema_urn = "urn:com:example2"
}

// LDAP mapping SCIM resource type is needed for the correlated data view resource
resource "pingdirectory_scim_resource_type" "myLdapMappingScimResourceType" {
  id          = "MyLdapMappingScimResourceType2"
  type = "ldap-mapping"
  core_schema = pingdirectory_scim_schema.myScimSchema.schema_urn
  enabled     = false
  endpoint    = "myendpoint"
}

resource "pingdirectory_correlated_ldap_data_view" "myCorrelatedLdapDataView" {
  id                              = "MyCorrelatedLdapDataView"
  scim_resource_type_name         = pingdirectory_scim_resource_type.myLdapMappingScimResourceType.id
  structural_ldap_objectclass     = "ldapObject"
  include_base_dn                 = "cn=com.example"
  primary_correlation_attribute   = "cn"
  secondary_correlation_attribute = "cn"
}
