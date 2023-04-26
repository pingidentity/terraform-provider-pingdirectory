---
page_title: "pingdirectory_default_unboundid_ms_chap_v2_sasl_mechanism_handler Resource - terraform-provider-pingdirectory"
subcategory: "Sasl Mechanism Handler"
description: |-
  Manages a Unboundid Ms Chap V2 Sasl Mechanism Handler.
---

# pingdirectory_default_unboundid_ms_chap_v2_sasl_mechanism_handler (Resource)

Manages a Unboundid Ms Chap V2 Sasl Mechanism Handler.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `description` (String) A description for this SASL Mechanism Handler
- `enabled` (Boolean) Indicates whether the SASL mechanism handler is enabled for use.
- `identity_mapper` (String) The identity mapper that should be used to identify the entry associated with the username provided in the bind request.

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


