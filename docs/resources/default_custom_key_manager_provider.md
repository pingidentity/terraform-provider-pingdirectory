---
page_title: "pingdirectory_default_custom_key_manager_provider Resource - terraform-provider-pingdirectory"
subcategory: "Key Manager Provider"
description: |-
  Manages a Custom Key Manager Provider.
---

# pingdirectory_default_custom_key_manager_provider (Resource)

Manages a Custom Key Manager Provider.

## Example Usage

```terraform
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

// For default, resource must exist.
// OOTB defaults are JKS / Null / PKCS11 / PKCS12 (Null used in example)
// enabled is required
resource "pingdirectory_default_custom_key_manager_provider" "myCustomKeyManagerProvider" {
  id      = "Null"
  enabled = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `description` (String) A description for this Key Manager Provider
- `enabled` (Boolean) Indicates whether the Key Manager Provider is enabled for use.

### Read-Only

- `last_updated` (String) Timestamp of the last Terraform update of this resource.
- `notifications` (Set of String) Notifications returned by the PingDirectory Configuration API.
- `required_actions` (Set of Object) Required actions returned by the PingDirectory Configuration API. (see [below for nested schema](#nestedatt--required_actions))

<a id="nestedatt--required_actions"></a>
### Nested Schema for `required_actions`

Read-Only:

- `property` (String)
- `synopsis` (String)
- `type` (String)

## Import

Import is supported using the following syntax:

```shell
# "customKeyManagerProviderId" should be the id of the Custom Key Manager Provider to be imported
terraform import pingdirectory_default_custom_key_manager_provider.myCustomKeyManagerProvider customKeyManagerProviderId
```
