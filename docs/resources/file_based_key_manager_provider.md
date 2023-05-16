---
page_title: "pingdirectory_file_based_key_manager_provider Resource - terraform-provider-pingdirectory"
subcategory: "Key Manager Provider"
description: |-
  Manages a File Based Key Manager Provider.
---

# pingdirectory_file_based_key_manager_provider (Resource)

Manages a File Based Key Manager Provider.

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

resource "pingdirectory_file_based_key_manager_provider" "myFileBasedKeyManagerProvider" {
  id             = "MyFileBasedKeyManagerProvider"
  description    = "My file based key manager provider"
  enabled        = false
  key_store_file = "/tmp/key-store-file"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `enabled` (Boolean) Indicates whether the Key Manager Provider is enabled for use.
- `id` (String) Name of this object.
- `key_store_file` (String) Specifies the path to the file that contains the private key information. This may be an absolute path, or a path that is relative to the Directory Server instance root.

### Optional

- `description` (String) A description for this Key Manager Provider
- `key_store_pin` (String, Sensitive) Specifies the PIN needed to access the File Based Key Manager Provider.
- `key_store_pin_file` (String) Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the File Based Key Manager Provider.
- `key_store_pin_passphrase_provider` (String) The passphrase provider to use to obtain the clear-text PIN needed to access the File Based Key Manager Provider.
- `key_store_type` (String) Specifies the format for the data in the key store file.
- `private_key_pin` (String, Sensitive) Specifies the clear-text PIN needed to access the File Based Key Manager Provider private key. If no private key PIN is specified the PIN defaults to the key store PIN.
- `private_key_pin_file` (String) Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the File Based Key Manager Provider private key. If no private key PIN is specified the PIN defaults to the key store PIN.
- `private_key_pin_passphrase_provider` (String) The passphrase provider to use to obtain the clear-text PIN needed to access the File Based Key Manager Provider private key. If no private key PIN is specified the PIN defaults to the key store PIN.

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
# "fileBasedKeyManagerProviderId" should be the id of the File Based Key Manager Provider to be imported
terraform import pingdirectory_file_based_key_manager_provider.myFileBasedKeyManagerProvider fileBasedKeyManagerProviderId
```
