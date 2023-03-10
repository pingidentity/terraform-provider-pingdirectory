---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_default_attribute_mapper_plugin Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Attribute Mapper Plugin.
---

# pingdirectory_default_attribute_mapper_plugin (Resource)

Manages a Attribute Mapper Plugin.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `always_map_responses` (Boolean) Indicates whether the target attribute in response messages should always be remapped back to the source attribute. If this is "false", then the mapping will be performed for a response message only if one or more elements of the associated request are mapped. Otherwise, the mapping will be performed for all responses regardless of whether the mapping was applied to the request.
- `description` (String) A description for this Plugin
- `enable_control_mapping` (Boolean) Indicates whether mapping should be applied to attribute types that may be present in specific controls. If enabled, attribute mapping will only be applied for control types which are specifically supported by the attribute mapper plugin.
- `enabled` (Boolean) Indicates whether the plug-in is enabled for use.
- `invoke_for_internal_operations` (Boolean) Indicates whether the plug-in should be invoked for internal operations.
- `plugin_type` (Set of String) Specifies the set of plug-in types for the plug-in, which specifies the times at which the plug-in is invoked.
- `source_attribute` (String) Specifies the source attribute type that may appear in client requests which should be remapped to the target attribute. Note that the source attribute type must be defined in the server schema and must not be equal to the target attribute type.
- `target_attribute` (String) Specifies the target attribute type to which the source attribute type should be mapped. Note that the target attribute type must be defined in the server schema and must not be equal to the source attribute type.

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


