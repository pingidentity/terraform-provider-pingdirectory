---
page_title: "pingdirectory_default_file_based_key_manager_provider Resource - terraform-provider-pingdirectory"
subcategory: "Key Manager Provider"
description: |-
  Manages a File Based Key Manager Provider.
---

# pingdirectory_default_file_based_key_manager_provider (Resource)

Manages a File Based Key Manager Provider.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `description` (String) A description for this Key Manager Provider
- `enabled` (Boolean) Indicates whether the Key Manager Provider is enabled for use.
- `key_store_file` (String) Specifies the path to the file that contains the private key information. This may be an absolute path, or a path that is relative to the Directory Server instance root.
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


