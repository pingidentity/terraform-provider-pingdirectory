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
  product_version        = "9.2.0.0"
}

# Use "pingdirectory_default_consent_definition" if you are adopting existing configuration from the PingDirectory server into Terraform
resource "pingdirectory_consent_definition" "myConsentDefinition" {
  unique_id    = "myConsentDefinition"
  display_name = "example display name"
}

# Use "pingdirectory_default_consent_definition_localization" if you are adopting existing configuration from the PingDirectory server into Terraform
resource "pingdirectory_consent_definition_localization" "myConsentDefinitionLocalization" {
  consent_definition_name = pingdirectory_consent_definition.myConsentDefinition.unique_id
  locale                  = "en-US"
  version                 = "1.1"
  data_text               = "example data text"
  purpose_text            = "example purpose text"
}
