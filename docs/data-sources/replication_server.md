---
page_title: "pingdirectory_replication_server Data Source - terraform-provider-pingdirectory"
subcategory: "Replication Server"
description: |-
  Describes a Replication Server.
---

# pingdirectory_replication_server (Data Source)

Describes a Replication Server.

Replication Servers publish updates to Directory Server instances within a Replication Domain.

## Example Usage

```terraform
data "pingdirectory_replication_server" "myReplicationServer" {
  synchronization_provider_name = "MySynchronizationProvider"
  replication_server_id         = 1234
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `replication_server_id` (Number) Specifies a unique identifier for the Replication Server.
- `synchronization_provider_name` (String) Name of the parent Synchronization Provider

### Read-Only

- `compression_criteria` (String) Specifies when the replication traffic should be compressed.
- `gateway_priority` (Number) Specifies the gateway priority of the Replication Server in the current location.
- `heartbeat_interval` (String) Specifies the heartbeat interval that the Directory Server will use when communicating with Replication Servers.
- `id` (String) The ID of this resource.
- `include_all_remote_servers_state_in_monitor_message` (Boolean) Supported in PingDirectory product version 10.0.0.0+. Indicates monitor messages should include information about remote servers.
- `je_property` (Set of String) Specifies the database and environment properties for the Berkeley DB Java Edition database for the replication changelog.
- `listen_on_all_addresses` (Boolean) Indicates whether the Replication Server should listen on all addresses for this host. If set to FALSE, then the Replication Server will listen only to the address resolved from the hostname provided.
- `missing_changes_alert_threshold_percent` (Number) Specifies the missing changes alert threshold as a percentage of the total pending changes. For instance, a value of 80 indicates that the replica is 80% of the way to losing changes.
- `missing_changes_policy` (String) Supported in PingDirectory product version 10.0.0.0+. Determines how the server responds when replication detects that some changes might have been missed. Each missing changes policy is a set of missing changes actions to take for a set of missing changes types. The value configured here acts as a default for all replication domains on this replication server.
- `remote_monitor_update_interval` (String) Specifies the duration that topology monitor data will be cached before it is requested again from a remote server.
- `replication_db_directory` (String) The path where the Replication Server stores all persistent information.
- `replication_port` (Number) The port on which this Replication Server waits for connections from other Replication Servers or Directory Server instances.
- `replication_purge_delay` (String) Changes are guaranteed to be maintained in the changelog database for at least this duration. Setting target-database-size can allow additional changes to be maintained up to the configured size on disk.
- `restricted_domain` (Set of String) Specifies the base DN of domains that are only replicated between server instances that belong to the same replication set.
- `target_database_size` (String) The replication changelog database is allowed to grow up to this size even if changes are older than the configured replication-purge-delay.
- `type` (String) The type of Replication Server resource. Options are ['replication-server']

