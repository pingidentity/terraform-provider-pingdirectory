---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_default_periodic_gc_plugin Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Periodic Gc Plugin.
---

# pingdirectory_default_periodic_gc_plugin (Resource)

Manages a Periodic Gc Plugin.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `delay_after_alert` (String) Specifies the length of time that the Directory Server should wait after sending the "force-gc-starting" administrative alert before actually invoking the garbage collection processing.
- `delay_post_gc` (String) Specifies the length of time that the Directory Server should wait after successfully completing the garbage collection processing, before removing the "force-gc-starting" administrative alert, which marks the server as unavailable.
- `description` (String) A description for this Plugin
- `enabled` (Boolean) Indicates whether the plug-in is enabled for use.
- `invoke_for_internal_operations` (Boolean) Indicates whether the plug-in should be invoked for internal operations.
- `invoke_gc_day_of_week` (Set of String) Specifies the days of the week which the Periodic GC Plugin should run. If no values are provided, then the plugin will run every day at the specified time.
- `invoke_gc_time_utc` (Set of String) Specifies the times of the day at which garbage collection may be explicitly invoked. The times should be specified in "HH:MM" format, with "HH" as a two-digit numeric value between 00 and 23 representing the hour of the day, and MM as a two-digit numeric value between 00 and 59 representing the minute of the hour. All times will be interpreted in the UTC time zone.
- `plugin_type` (Set of String) Specifies the set of plug-in types for the plug-in, which specifies the times at which the plug-in is invoked.

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


