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

resource "pingdirectory_id_token_validator" "myPingOneIdTokenValidator" {
  type                   = "ping-one"
  name                   = "MyPingOneIdTokenValidator"
  issuer_url             = "example.com"
  enabled                = false
  identity_mapper        = "Exact Match"
  evaluation_order_index = 1
}

resource "pingdirectory_token_claim_validation" "myTokenClaimValidation" {
  name                    = "MyTokenClaimValidation"
  id_token_validator_name = pingdirectory_id_token_validator.myPingOneIdTokenValidator.id
  any_required_value      = ["my_example_value"]
  claim_name              = "my_example_claim_name"
}
