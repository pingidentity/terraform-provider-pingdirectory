terraform {
  required_version = ">=1.1"
  required_providers {
    pingdirectory = {
      version = "~> 1.0.0"
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
  product_version        = "10.2.0.0"
}

variable "pingfederate_hostname" {
  type     = string
  nullable = false
  default  = "pingfederate"
}

variable "pingfederate_https_port" {
  type     = number
  nullable = false
  default  = 9031
}

variable "root_user_dn" {
  type     = string
  nullable = false
  default  = "cn=administrator"
}

variable "root_user_password" {
  type      = string
  nullable  = false
  default   = "2FederateM0re"
  sensitive = true
}

variable "user_base_dn" {
  type     = string
  nullable = false
  default  = "dc=example,dc=com"
}
