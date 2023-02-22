---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_profiler_plugin Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Profiler Plugin.
---

# pingdirectory_profiler_plugin (Resource)

Manages a Profiler Plugin.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `description` (String) A description for this Plugin
- `enable_profiling_on_startup` (Boolean) Indicates whether the profiler plug-in is to start collecting data automatically when the Directory Server is started.
- `enabled` (Boolean) Indicates whether the plug-in is enabled for use.
- `profile_action` (String) Specifies the action that should be taken by the profiler.
- `profile_directory` (String) Specifies the path to the directory where profile information is to be written. This path may be either an absolute path or a path that is relative to the root of the Directory Server instance.
- `profile_sample_interval` (String) Specifies the sample interval in milliseconds to be used when capturing profiling information in the server.

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

