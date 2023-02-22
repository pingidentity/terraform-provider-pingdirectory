---
page_title: "Provider: PingDirectory"
description: |-
  The PingDirectory provider is used to manage the configuration of a PingDirectory server through the Configuration API.
---

# PingDirectory Provider

The PingDirectory provider manages the configuration of a PingDirectory server through the Configuration API. The Configuration API requires credentials for basic auth, which must be passed to the provider.

## Providing credentials

The server host, username, and password can either be provided in the Terraform configuration file, or they can be provided via environment variables:

```
PINGDIRECTORY_PROVIDER_HTTPS_HOST
PINGDIRECTORY_PROVIDER_USERNAME
PINGDIRECTORY_PROVIDER_PASSWORD
```

## An example managing several config objects

```terraform
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

resource "pingdirectory_location" "drangleic" {
  id          = "Drangleic"
  description = "Seek the king"
}

resource "pingdirectory_global_configuration" "global" {
  location            = "Docker"
  encrypt_data        = false
  sensitive_attribute = ["Delivered One-Time Password", "TOTP Shared Secret"]
  tracked_application = ["Requests by Root Users"]
  result_code_map     = "Sun DS Compatible Behavior"
  disabled_privilege  = ["jmx-write", "jmx-read"]
}

resource "pingdirectory_blind_trust_manager_provider" "blindtest" {
  id                          = "Blind Test"
  enabled                     = true
  include_jvm_default_issuers = true
}

resource "pingdirectory_file_based_trust_manager_provider" "filetest" {
  id                          = "FileTest"
  enabled                     = true
  trust_store_file            = "config/keystore"
  trust_store_type            = "pkcs12"
  include_jvm_default_issuers = true
}

resource "pingdirectory_jvm_default_trust_manager_provider" "jvmtest" {
  id      = "jvmtest"
  enabled = false
}

resource "pingdirectory_third_party_trust_manager_provider" "tptest" {
  id                 = "tptest"
  enabled            = false
  extension_class    = "com.unboundid.directory.sdk.common.api.TrustManagerProvider"
  extension_argument = ["val1=one", "val2=two"]
}
```
