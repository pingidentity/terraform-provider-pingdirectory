---
page_title: "pingdirectory_soft_delete_policies Data Source - terraform-provider-pingdirectory"
subcategory: "Soft Delete Policy"
description: |-
  Lists Soft Delete Policy objects in the server configuration.
---

# pingdirectory_soft_delete_policies (Data Source)

Lists Soft Delete Policy objects in the server configuration.

General policy settings for soft delete operations.

## Example Usage

```terraform
data "pingdirectory_soft_delete_policies" "list" {
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_ds_config_soft_deletes_on_server)

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (String) SCIM filter used when searching the configuration.

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (Set of String) Soft Delete Policy IDs found in the configuration

