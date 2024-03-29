---
page_title: "pingdirectory_custom_logged_stats_list Data Source - terraform-provider-pingdirectory"
subcategory: "Custom Logged Stats"
description: |-
  Lists Custom Logged Stats objects in the server configuration.
---

# pingdirectory_custom_logged_stats_list (Data Source)

Lists Custom Logged Stats objects in the server configuration.

A custom Custom Logged Stats object enables additional statistics to be included in the output of a Periodic Stats Logger.

## Example Usage

```terraform
data "pingdirectory_custom_logged_stats_list" "list" {
  plugin_name = "MyPlugin"
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_ds_config_custom_logged_stat_dsconfig_int)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `plugin_name` (String) Name of the parent Plugin

### Optional

- `filter` (String) SCIM filter used when searching the configuration.

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (Set of String) Custom Logged Stats IDs found in the configuration

