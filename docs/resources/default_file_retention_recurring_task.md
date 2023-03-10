---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_default_file_retention_recurring_task Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a File Retention Recurring Task.
---

# pingdirectory_default_file_retention_recurring_task (Resource)

Manages a File Retention Recurring Task.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `alert_on_failure` (Boolean) Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task fails to complete successfully.
- `alert_on_start` (Boolean) Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task starts running.
- `alert_on_success` (Boolean) Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task completes successfully.
- `cancel_on_task_dependency_failure` (Boolean) Indicates whether an instance of this Recurring Task should be canceled if the task immediately before it in the recurring task chain fails to complete successfully (including if it is canceled by an administrator before it starts or while it is running).
- `description` (String) A description for this Recurring Task
- `email_on_failure` (Set of String) The email addresses to which a message should be sent if an instance of this Recurring Task fails to complete successfully. If this option is used, then at least one smtp-server must be configured in the global configuration.
- `email_on_start` (Set of String) The email addresses to which a message should be sent whenever an instance of this Recurring Task starts running. If this option is used, then at least one smtp-server must be configured in the global configuration.
- `email_on_success` (Set of String) The email addresses to which a message should be sent whenever an instance of this Recurring Task completes successfully. If this option is used, then at least one smtp-server must be configured in the global configuration.
- `filename_pattern` (String) A pattern that specifies the names of the files to examine. The pattern may contain zero or more asterisks as wildcards, where each wildcard matches zero or more characters. It may also contain at most one occurrence of the special string "${timestamp}", which will match a timestamp with the format specified using the timestamp-format property. All other characters in the pattern will be treated literally.
- `retain_aggregate_file_size` (String) The minimum aggregate size of files that will be retained. The size should be specified as an integer followed by a unit that is one of "b" or "bytes", "kb" or "kilobytes", "mb" or "megabytes", "gb" or "gigabytes", or "tb" or "terabytes". For example, a value of "1 gb" indicates that at least one gigabyte of files should be retained.
- `retain_file_age` (String) The minimum age of files matching the pattern that will be retained.
- `retain_file_count` (Number) The minimum number of files matching the pattern that will be retained.
- `target_directory` (String) The path to the directory containing the files to examine. The directory must exist.
- `timestamp_format` (String) The format to use for the timestamp represented by the "${timestamp}" token in the filename pattern.

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


