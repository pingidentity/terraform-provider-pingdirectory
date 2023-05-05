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

// extension_class must exist  Make sure the class is neither abstract, nor an interface, and defines a public, nullary constructor.
// see https://docs.ping.directory/PingDirectory/7.0.1.4/config-guide/third-party-monitor-provider.html
resource "pingdirectory_third_party_monitor_provider" "myThirdPartyMonitorProvider" {
  id              = "3rdPartyNew"
  description     = "My new third party monitor provider"
  extension_class = "com.unboundid.directory.sdk.common.api.MonitorProvider"
  enabled         = false
}
