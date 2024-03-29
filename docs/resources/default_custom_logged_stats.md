---
page_title: "pingdirectory_default_custom_logged_stats Resource - terraform-provider-pingdirectory"
subcategory: "Custom Logged Stats"
description: |-
  Manages a Custom Logged Stats.
---

# pingdirectory_default_custom_logged_stats (Resource)

Manages a Custom Logged Stats.

A custom Custom Logged Stats object enables additional statistics to be included in the output of a Periodic Stats Logger.

Since this is a 'default' resource, the managed object must already exist in the PingDirectory configuration.



## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_ds_config_custom_logged_stat_dsconfig_int)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of this config object.
- `plugin_name` (String) Name of the parent Plugin

### Optional

- `attribute_to_log` (Set of String) Specifies the attributes on the monitor entries that should be included in the output.
- `column_name` (Set of String) Optionally, specifies an explicit name for each column header instead of having these names automatically generated from the monitored attribute name.
- `decimal_format` (String) This provides a way to format the monitored attribute value in the output to control the precision for instance.
- `description` (String) A description for this Custom Logged Stats
- `divide_value_by` (String) An optional floating point value that can be used to scale the resulting value.
- `divide_value_by_attribute` (String) An optional property that can scale the resulting value by another attribute in the monitored entry.
- `enabled` (Boolean) Indicates whether the Custom Logged Stats object is enabled.
- `header_prefix` (String) An optional prefix that is included in the header before the column name.
- `header_prefix_attribute` (String) An optional attribute from the monitor entry that is included as a prefix before the column name in the column header.
- `include_filter` (String) An optional LDAP filter that can be used restrict which monitor entries are used to produce the output.
- `monitor_objectclass` (String) The objectclass name of the monitor entries to examine for generating these statistics.
- `non_zero_implies_not_idle` (Boolean) If this property is set to true, then the value of any of the monitored attributes here can contribute to whether an interval is considered "idle" by the Periodic Stats Logger.
- `regex_pattern` (String) An optional regular expression pattern, that when used in conjunction with regex-replacement, can alter the value of the attribute being monitored.
- `regex_replacement` (String) An optional regular expression replacement value, that when used in conjunction with regex-pattern, can alter the value of the attribute being monitored.
- `statistic_type` (Set of String) Specifies the type of statistic to include in the output for each monitored attribute.

### Read-Only

- `id` (String) The ID of this resource.
- `notifications` (Set of String) Notifications returned by the PingDirectory Configuration API.
- `required_actions` (Set of Object) Required actions returned by the PingDirectory Configuration API. (see [below for nested schema](#nestedatt--required_actions))
- `type` (String) The type of Custom Logged Stats resource. Options are ['custom-logged-stats']

<a id="nestedatt--required_actions"></a>
### Nested Schema for `required_actions`

Read-Only:

- `property` (String)
- `synopsis` (String)
- `type` (String)



