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

resource "pingdirectory_topology_admin_user" "myTopologyAdminUser" {
  name                            = "MyTopologyAdminUser"
  inherit_default_root_privileges = true
  search_result_entry_limit       = 100
  time_limit_seconds              = 60
  look_through_entry_limit        = 20
  idle_time_limit_seconds         = 120
  password_policy                 = "Default Password Policy"
  require_secure_authentication   = true
  require_secure_connections      = false
}
