---
page_title: "pingdirectory_default_trusted_certificate Resource - terraform-provider-pingdirectory"
subcategory: "Trusted Certificate"
description: |-
  Manages a Trusted Certificate.
---

# pingdirectory_default_trusted_certificate (Resource)

Manages a Trusted Certificate.

The Trusted Certificate represents a trusted public key that may be used to verify credentials for digital signatures and public-key encryption. The public key is represented as an X.509v3 certificate. For example, when configured on an access token validator, it may be used to validate the signature of an incoming JWT access token before the product REST APIs consume the access token for Bearer token authentication.

Since this is a 'default' resource, the managed object must already exist in the PingDirectory configuration.



## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_ds_use_locally_config_trusted_cert)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of this config object.

### Optional

- `certificate` (String) The PEM-encoded X.509v3 certificate.

### Read-Only

- `id` (String) The ID of this resource.
- `notifications` (Set of String) Notifications returned by the PingDirectory Configuration API.
- `required_actions` (Set of Object) Required actions returned by the PingDirectory Configuration API. (see [below for nested schema](#nestedatt--required_actions))
- `type` (String) The type of Trusted Certificate resource. Options are ['trusted-certificate']

<a id="nestedatt--required_actions"></a>
### Nested Schema for `required_actions`

Read-Only:

- `property` (String)
- `synopsis` (String)
- `type` (String)



