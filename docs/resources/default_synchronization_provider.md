---
page_title: "pingdirectory_default_synchronization_provider Resource - terraform-provider-pingdirectory"
subcategory: "Synchronization Provider"
description: |-
  Manages a Synchronization Provider.
---

# pingdirectory_default_synchronization_provider (Resource)

Manages a Synchronization Provider.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.
- `type` (String) The type of Synchronization Provider resource. Options are ['replication', 'custom']

### Optional

- `description` (String) A description for this Synchronization Provider
- `enabled` (Boolean) Indicates whether the Synchronization Provider is enabled for use.
- `num_update_replay_threads` (Number) Specifies the number of update replay threads.

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


