terraform {
  required_version = ">=1.1"
  required_providers {
    pingdirectory = {
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
}

# Use "pingdirectory_default_reverse_dn_join_virtual_attribute" if you are adopting existing configuration from the PingDirectory server into Terraform
resource "pingdirectory_reverse_dn_join_virtual_attribute" "myReverseDnJoinVirtualAttribute" {
  id                = "MyReverseDnJoinVirtualAttribute"
  join_dn_attribute = "sn"
  join_base_dn_type = "use-search-base-dn"
  enabled           = false
  attribute_type    = "cn"
}
