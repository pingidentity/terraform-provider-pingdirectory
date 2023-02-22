---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_snmp_subagent_plugin Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Snmp Subagent Plugin.
---

# pingdirectory_snmp_subagent_plugin (Resource)

Manages a Snmp Subagent Plugin.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `agentx_address` (String) The hostname or IP address of the SNMP master agent.
- `agentx_port` (Number) The port number on which the SNMP master agent will be contacted.
- `enabled` (Boolean) Indicates whether the plug-in is enabled for use.
- `id` (String) Name of this object.

### Optional

- `connect_retry_max_wait` (String) The maximum amount of time to wait between attempts to establish a connection to the master agent.
- `context_name` (String) The SNMP context name for this sub-agent. The context name must not be longer than 30 ASCII characters. Each server in a topology must have a unique SNMP context name.
- `description` (String) A description for this Plugin
- `invoke_for_internal_operations` (Boolean) Indicates whether the plug-in should be invoked for internal operations.
- `num_worker_threads` (Number) The number of worker threads to use to handle SNMP requests.
- `ping_interval` (String) The amount of time between consecutive pings sent by the sub-agent on its connection to the master agent. A value of zero disables the sending of pings by the sub-agent.
- `session_timeout` (String) Specifies the maximum amount of time to wait for a session to the master agent to be established.

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

