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

variable "root_user_dn" {
  type     = string
  nullable = false
  default  = "cn=administrator"
}

resource "pingdirectory_default_connection_handler" "defaultHttpsConnectionHandler" {
  type                      = "http"
  name                      = "HTTPS Connection Handler"
  web_application_extension = []
}

resource "pingdirectory_default_gauge" "defaultCpuUsageGauge" {
  type    = "numeric"
  name    = "CPU Usage (Percent)"
  enabled = true
}

resource "pingdirectory_default_gauge" "defaultLicenseExpirationGauge" {
  type    = "numeric"
  name    = "License Expiration (Days)"
  enabled = false
}

resource "pingdirectory_default_gauge" "defaultAvailableFileDescriptorsGauge" {
  type    = "numeric"
  name    = "Available File Descriptors"
  enabled = false
}

resource "pingdirectory_default_log_publisher" "defaultDataRecoveryLog" {
  type    = "file-based-audit"
  name    = "Data Recovery Log"
  enabled = false
}

resource "pingdirectory_default_log_publisher" "defaultFileBasedAuditLogger" {
  type    = "file-based-audit"
  name    = "File-Based Audit Logger"
  enabled = true
}

resource "pingdirectory_default_root_dn_user" "defaultRootDnUser" {
  name              = "Directory Manager"
  alternate_bind_dn = [var.root_user_dn]
}
