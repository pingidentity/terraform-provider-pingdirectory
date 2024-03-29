---
page_title: "pingdirectory_sensitive_attributes Data Source - terraform-provider-pingdirectory"
subcategory: "Sensitive Attribute"
description: |-
  Lists Sensitive Attribute objects in the server configuration.
---

# pingdirectory_sensitive_attributes (Data Source)

Lists Sensitive Attribute objects in the server configuration.

Sensitive Attributes provide a means of declaring one or more attributes to contain sensitive data so that the server can enforce additional protection for operations attempting to interact with them.

## Example Usage

```terraform
data "pingdirectory_sensitive_attributes" "list" {
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_ds_config_sensitive_attriutes)

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (String) SCIM filter used when searching the configuration.

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (Set of String) Sensitive Attribute IDs found in the configuration

