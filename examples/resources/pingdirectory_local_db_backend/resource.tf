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

# Use "pingdirectory_default_local_db_backend" if you are adopting existing configuration from the PingDirectory server into Terraform
resource "pingdirectory_backend" "myLocalDbBackend" {
  type = "local-db"
  backend_id            = "MyLocalDbBackendasdfdfdd"
  base_dn               = ["dc=exampleanotherdfdd,dc=com"]
  writability_mode      = "enabled"
  db_directory          = "db"
  import_temp_directory = "tmp"
  enabled               = true
}
