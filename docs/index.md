---
page_title: "Provider: PingDirectory"
description: |-
  The PingDirectory provider is used to manage the configuration of a PingDirectory server through the Configuration API.
---

# PingDirectory Provider

The PingDirectory provider manages the configuration of a PingDirectory server through the Configuration API. The provider only manages configuration, similar to the `dsconfig` command-line tool. The provider does not manage other aspects of the PingDirectory server, such as schema and user data.

The Configuration API requires credentials for basic auth, which must be passed to the provider.

## PingDirectory Version Support

The PingDirectory provider supports versions `9.3` through `10.2` of PingDirectory.

## Documentation

Detailed documentation on PingDirectory configuration can be found in the [online docs](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_ds_landing_page).

## Example usage and relation to the `dsconfig` tool

The following example configures a `Location` object and updates the PingDirectory `Global Configuration`.

Applying this Terraform configuration file will create a `Location` on the PingDirectory server managed by Terraform, and will update the PingDirectory `Global Configuration` to match the defined resource. Terraform can then manage this configuration, rather than using dsconfig commands such as:

```
dsconfig create-location --location-name MyLocation --set "description:My description"
dsconfig set-location-prop --location-name MyLocation --set "description:My changed description"
dsconfig delete-location --location-name MyLocation
dsconfig set-global-configuration-prop --set encrypt-data:true --set location:MyLocation
```

```terraform
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
  insecure_trust_all_tls = true
  product_version        = "10.2.0.0"
}

# Create a sample location
resource "pingdirectory_location" "myLocation" {
  name        = "MyLocation"
  description = "My description"
}

# Update the default global configuration to use the created location, and to enable encryption
resource "pingdirectory_default_global_configuration" "global" {
  location     = pingdirectory_location.myLocation.id
  encrypt_data = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `ca_certificate_pem_files` (Set of String) Paths to files containing PEM-encoded certificates to be trusted as root CAs when connecting to the PingDirectory server over HTTPS. If not set, the host's root CA set will be used. Default value can be set with the `PINGDIRECTORY_PROVIDER_CA_CERTIFICATE_PEM_FILES` environment variable, using commas to delimit multiple PEM files if necessary.
- `https_host` (String) URI for PingDirectory HTTPS port. Default value can be set with the `PINGDIRECTORY_PROVIDER_HTTPS_HOST` environment variable.
- `insecure_trust_all_tls` (Boolean) Set to true to trust any certificate when connecting to the PingDirectory server. This is insecure and should not be enabled outside of testing. Default value can be set with the `PINGDIRECTORY_PROVIDER_INSECURE_TRUST_ALL_TLS` environment variable.
- `password` (String, Sensitive) Password for PingDirectory admin user. Default value can be set with the `PINGDIRECTORY_PROVIDER_PASSWORD` environment variable.
- `product_version` (String) Version of the PingDirectory server being configured. Default value can be set with the `PINGDIRECTORY_PROVIDER_PRODUCT_VERSION` environment variable.
- `username` (String) Username for PingDirectory admin user. Default value can be set with the `PINGDIRECTORY_PROVIDER_USERNAME` environment variable.

## Server profile examples

Examples duplicating behavior of existing [server profiles](https://github.com/pingidentity/pingidentity-server-profiles) can be found in the `examples/server-profiles` directory.

## Default resources

PingDirectory comes with many default values configured out of the box. To make adopting these pre-configured values easier, resources in the provider support a "default_" prefix which allows adopting default resources from the PingDirectory server without having to specify every attribute to match the default configuration that PingDirectory provides. When using "default_" resources, any optional values that are not specified in your HCL code will be computed based on the values that PingDirectory provides. Required values must still be specified when using "default_".

When using the "default_" prefix, the configuration object being managed must already exist on the PingDirectory server prior to being managed by the provider. The provider will not create the configuration object.

When destroying default resources, the provider will stop managing the resource and remove it from state, but it will not destroy the resource in the PingDirectory configuration.

Here is an example showing the differences between default and non-default resources.

```terraform
# Disable the default failed operations access logger
resource "pingdirectory_default_log_publisher" "defaultFileBasedAccessLogPublisher" {
  name    = "Failed Operations Access Logger"
  enabled = false
}

# Create a new custom file based access logger
resource "pingdirectory_log_publisher" "myNewFileBasedAccessLogPublisher" {
  type                 = "file-based-access"
  name                 = "MyNewFileBasedAccessLogPublisher"
  log_file             = "logs/example.log"
  log_file_permissions = "600"
  rotation_policy      = ["Size Limit Rotation Policy"]
  retention_policy     = ["File Count Retention Policy"]
  asynchronous         = true
  enabled              = false
}

# Enable the default JMX connection handler
resource "pingdirectory_default_connection_handler" "defaultJmxConnHandler" {
  name    = "JMX Connection Handler"
  enabled = true
}

# Create a new custom JMX connection handler
resource "pingdirectory_connection_handler" "myJmxConnHandler" {
  type        = "jmx"
  name        = "MyJmxConnHandler"
  enabled     = false
  listen_port = 8888
}
```

## The `type` attribute

All resources in this provider include a `type` attribute that serves as a discriminator between config objects that support multiple types. For example, the `pingdirectory_log_publisher` resource supports `type` values such as `file-based-error`, `json-access`, `console-json-audit`, etc. It is required for any non-default resources that support multiple types.

For example, the following `dsconfig` command and HCL code are equivalent:

```
dsconfig create-log-publisher \
  --publisher-name MyLogPublisher \
  --type file-based-debug \
  --set enabled:false \
  --set log-file:logs/example.log \
  --set "rotation-policy:24 Hours Time Limit Rotation Policy" \
  --set "retention-policy:File Count Growth Limit Policy"

resource "pingdirectory_log_publisher" "myLogPublisher" {
  name             = "MyLogPublisher"
  type             = "file-based-debug"
  enabled          = false
  log_file         = "logs/example.log"
  rotation_policy  = "24 Hours Time Limit Rotation Policy"
  retention_policy = "File Count Growth Limit Policy"
}
```

The `pingdirectory_plugin`, `pingdirectory_scim_attribute`, and `pingdirectory_scim_subattribute` resources are exceptions, as they already have a `type` attribute used by PingDirectory. The `resource_type` attribute serves the purpose that would normally be filled by the `type` attribute for these resources.
