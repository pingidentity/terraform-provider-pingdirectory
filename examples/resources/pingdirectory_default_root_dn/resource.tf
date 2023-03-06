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

// This set is approximately the minimum set required for you to be able to run
// 'dsconfig get-root-dn-prop' successfully.  If you remove any of these permissions, 
// you risk loss of access to the RootDN permission object.
resource "pingdirectory_default_root_dn" "myrootdn" {
  default_root_privilege_name = ["bypass-acl", "config-read", "config-write", "modify-acl", "privilege-change", "use-admin-session"]
}
