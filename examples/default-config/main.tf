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

# Disable the default failed operations access logger
resource "pingdirectory_default_file_based_access_log_publisher" "defaultFileBasedAccessLogPublisher" {
  id      = "Failed Operations Access Logger"
  enabled = false
}

# Create a new custom file based access logger
resource "pingdirectory_file_based_access_log_publisher" "myNewFileBasedAccessLogPublisher" {
  id                   = "MyNewFileBasedAccessLogPublisher"
  log_file             = "logs/example.log"
  log_file_permissions = "600"
  rotation_policy      = ["Size Limit Rotation Policy"]
  retention_policy     = ["File Count Retention Policy"]
  asynchronous         = true
  enabled              = false
}

# Enable the default JMX connection handler
resource "pingdirectory_default_jmx_connection_handler" "defaultJmxConnHandler" {
  id      = "JMX Connection Handler"
  enabled = true
}

# Create a new custom JMX connection handler
resource "pingdirectory_jmx_connection_handler" "myJmxConnHandler" {
  id          = "MyJmxConnHandler"
  enabled     = false
  listen_port = 8888
}
