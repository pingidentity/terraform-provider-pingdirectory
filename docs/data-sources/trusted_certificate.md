---
page_title: "pingdirectory_trusted_certificate Data Source - terraform-provider-pingdirectory"
subcategory: "Trusted Certificate"
description: |-
  Describes a Trusted Certificate.
---

# pingdirectory_trusted_certificate (Data Source)

Describes a Trusted Certificate.

The Trusted Certificate represents a trusted public key that may be used to verify credentials for digital signatures and public-key encryption. The public key is represented as an X.509v3 certificate. For example, when configured on an access token validator, it may be used to validate the signature of an incoming JWT access token before the product REST APIs consume the access token for Bearer token authentication.

## Example Usage

```terraform
data "pingdirectory_trusted_certificate" "myTrustedCertificate" {
  name = "MyTrustedCertificate"
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_ds_use_locally_config_trusted_cert)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of this config object.

### Read-Only

- `certificate` (String) The PEM-encoded X.509v3 certificate.
- `id` (String) The ID of this resource.
- `type` (String) The type of Trusted Certificate resource. Options are ['trusted-certificate']

