---
page_title: "pingdirectory_default_failure_lockout_action Resource - terraform-provider-pingdirectory"
subcategory: "Failure Lockout Action"
description: |-
  Manages a Failure Lockout Action.
---

# pingdirectory_default_failure_lockout_action (Resource)

Manages a Failure Lockout Action.

Failure Lockout Actions may be used to specify the behavior that the server should exhibit for accounts that have too many failed authentication attempts.

Since this is a 'default' resource, the managed object must already exist in the PingDirectory configuration.



## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_sec_alt_failure_lockout_actions)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of this config object.

### Optional

- `allow_blocking_delay` (Boolean) Indicates whether to delay the response for authentication attempts even if that delay may block the thread being used to process the attempt.
- `delay` (String) The length of time to delay the bind response for accounts with too many failed authentication attempts.
- `description` (String) A description for this Failure Lockout Action
- `generate_account_status_notification` (Boolean) When the `type` attribute is set to:
  - `delay-bind-response`: Indicates whether to generate an account status notification for cases in which a bind response is delayed because of failure lockout.
  - `no-operation`: Indicates whether to generate an account status notification for cases in which this failure lockout action is invoked for a bind attempt with too many outstanding authentication failures.

### Read-Only

- `id` (String) The ID of this resource.
- `notifications` (Set of String) Notifications returned by the PingDirectory Configuration API.
- `required_actions` (Set of Object) Required actions returned by the PingDirectory Configuration API. (see [below for nested schema](#nestedatt--required_actions))
- `type` (String) The type of Failure Lockout Action resource. Options are ['delay-bind-response', 'no-operation', 'lock-account']

<a id="nestedatt--required_actions"></a>
### Nested Schema for `required_actions`

Read-Only:

- `property` (String)
- `synopsis` (String)
- `type` (String)



