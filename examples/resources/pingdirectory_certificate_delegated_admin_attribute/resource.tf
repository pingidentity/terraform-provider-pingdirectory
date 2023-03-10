terraform {
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

# Use "pingdirectory_default_certificate_delegated_admin_attribute" if you are adopting existing configuration from the PingDirectory server into Terraform
resource "pingdirectory_certificate_delegated_admin_attribute" "myCertificateDelegatedAdminAttribute" {
  rest_resource_type_name = "MyRestResourceType"
  attribute_type          = "myattr"
  display_name            = "MyAttribute"
}
