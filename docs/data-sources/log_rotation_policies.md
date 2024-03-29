---
page_title: "pingdirectory_log_rotation_policies Data Source - terraform-provider-pingdirectory"
subcategory: "Log Rotation Policy"
description: |-
  Lists Log Rotation Policy objects in the server configuration.
---

# pingdirectory_log_rotation_policies (Data Source)

Lists Log Rotation Policy objects in the server configuration.

Log Rotation Policies are used to specify when log files should be rotated.

## Example Usage

```terraform
data "pingdirectory_log_rotation_policies" "list" {
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_ds_config_log_rotation)

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (String) SCIM filter used when searching the configuration.

### Read-Only

- `id` (String) The ID of this resource.
- `objects` (Set of Object) Log Rotation Policy objects found in the configuration (see [below for nested schema](#nestedatt--objects))

<a id="nestedatt--objects"></a>
### Nested Schema for `objects`

Read-Only:

- `id` (String)
- `type` (String)

