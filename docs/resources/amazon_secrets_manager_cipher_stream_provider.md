---
page_title: "pingdirectory_amazon_secrets_manager_cipher_stream_provider Resource - terraform-provider-pingdirectory"
subcategory: "Cipher Stream Provider"
description: |-
  Manages a Amazon Secrets Manager Cipher Stream Provider.
---

# pingdirectory_amazon_secrets_manager_cipher_stream_provider (Resource)

Manages a Amazon Secrets Manager Cipher Stream Provider.

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

resource "pingdirectory_amazon_secrets_manager_cipher_stream_provider" "myAmazonSecretsManagerCipherStreamProvider" {
  id                  = "MyAmazonSecretsManagerCipherStreamProvider"
  aws_external_server = "my_example_aws_external_server"
  secret_id           = "my_example_secret_id"
  secret_field_name   = "my_example_secret_field_name"
  enabled             = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `aws_external_server` (String) The external server with information to use when interacting with the AWS Secrets Manager.
- `enabled` (Boolean) Indicates whether this Cipher Stream Provider is enabled for use in the Directory Server.
- `id` (String) Name of this object.
- `secret_field_name` (String) The name of the JSON field whose value is the passphrase that will be used to generate the encryption key for protecting the contents of the encryption settings database.
- `secret_id` (String) The Amazon Resource Name (ARN) or the user-friendly name of the secret to be retrieved.

### Optional

- `description` (String) A description for this Cipher Stream Provider
- `encryption_metadata_file` (String) The path to a file that will hold metadata about the encryption performed by this Amazon Secrets Manager Cipher Stream Provider.
- `secret_version_id` (String) The unique identifier for the version of the secret to be retrieved.
- `secret_version_stage` (String) The staging label for the version of the secret to be retrieved.

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
# "amazonSecretsManagerCipherStreamProviderId" should be the id of the Amazon Secrets Manager Cipher Stream Provider to be imported
terraform import pingdirectory_amazon_secrets_manager_cipher_stream_provider.myAmazonSecretsManagerCipherStreamProvider amazonSecretsManagerCipherStreamProviderId
```
