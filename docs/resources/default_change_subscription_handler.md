---
page_title: "pingdirectory_default_change_subscription_handler Resource - terraform-provider-pingdirectory"
subcategory: "Change Subscription Handler"
description: |-
  Manages a Change Subscription Handler.
---

# pingdirectory_default_change_subscription_handler (Resource)

Manages a Change Subscription Handler.

Change Subscription Handlers may be used to provide notification of changes that match one or more change subscriptions.

Since this is a 'default' resource, the managed object must already exist in the PingDirectory configuration.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of this config object.

### Optional

- `change_subscription` (Set of String) The set of change subscriptions for which this change subscription handler should be notified. If no values are provided then it will be notified for all change subscriptions defined in the server.
- `description` (String) A description for this Change Subscription Handler
- `enabled` (Boolean) Indicates whether this change subscription handler is enabled within the server.
- `extension_argument` (Set of String) The set of arguments used to customize the behavior for the Third Party Change Subscription Handler. Each configuration property should be given in the form 'name=value'.
- `extension_class` (String) The fully-qualified name of the Java class providing the logic for the Third Party Change Subscription Handler.
- `log_file` (String) Specifies the log file in which the change notification messages will be written.
- `script_argument` (Set of String) The set of arguments used to customize the behavior for the Scripted Change Subscription Handler. Each configuration property should be given in the form 'name=value'.
- `script_class` (String) The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Change Subscription Handler.

### Read-Only

- `id` (String) The ID of this resource.
- `notifications` (Set of String) Notifications returned by the PingDirectory Configuration API.
- `required_actions` (Set of Object) Required actions returned by the PingDirectory Configuration API. (see [below for nested schema](#nestedatt--required_actions))
- `type` (String) The type of Change Subscription Handler resource. Options are ['groovy-scripted', 'logging', 'third-party']

<a id="nestedatt--required_actions"></a>
### Nested Schema for `required_actions`

Read-Only:

- `property` (String)
- `synopsis` (String)
- `type` (String)



