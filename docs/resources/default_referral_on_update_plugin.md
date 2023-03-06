---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_default_referral_on_update_plugin Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Referral On Update Plugin.
---

# pingdirectory_default_referral_on_update_plugin (Resource)

Manages a Referral On Update Plugin.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `base_dn` (Set of String) Specifies a base DN for requests for which to send referrals in response to update operations.
- `description` (String) A description for this Plugin
- `enabled` (Boolean) Indicates whether the plug-in is enabled for use.
- `invoke_for_internal_operations` (Boolean) Indicates whether the plug-in should be invoked for internal operations.
- `plugin_type` (Set of String) Specifies the set of plug-in types for the plug-in, which specifies the times at which the plug-in is invoked.
- `referral_base_url` (Set of String) Specifies the base URL to use for the referrals generated by this plugin. It should include only the scheme, address, and port to use to communicate with the target server (e.g., "ldap://server.example.com:389/").

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

