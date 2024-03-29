---
page_title: "pingdirectory_log_rotation_policy Data Source - terraform-provider-pingdirectory"
subcategory: "Log Rotation Policy"
description: |-
  Describes a Log Rotation Policy.
---

# pingdirectory_log_rotation_policy (Data Source)

Describes a Log Rotation Policy.

Log Rotation Policies are used to specify when log files should be rotated.

## Example Usage

```terraform
data "pingdirectory_log_rotation_policy" "myLogRotationPolicy" {
  name = "MyLogRotationPolicy"
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_ds_config_log_rotation)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of this config object.

### Read-Only

- `description` (String) A description for this Log Rotation Policy
- `file_size_limit` (String) Specifies the maximum size that a log file can reach before it is rotated.
- `id` (String) The ID of this resource.
- `rotation_interval` (String) Specifies the time interval between rotations.
- `time_of_day` (Set of String) Specifies the time of day at which log rotation should occur.
- `type` (String) The type of Log Rotation Policy resource. Options are ['time-limit', 'fixed-time', 'never-rotate', 'size-limit']

