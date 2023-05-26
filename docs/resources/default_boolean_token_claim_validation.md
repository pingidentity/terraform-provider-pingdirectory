---
page_title: "pingdirectory_default_boolean_token_claim_validation Resource - terraform-provider-pingdirectory"
subcategory: "Token Claim Validation"
description: |-
  Manages a Boolean Token Claim Validation.
---

# pingdirectory_default_boolean_token_claim_validation (Resource)

Manages a Boolean Token Claim Validation.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.
- `id_token_validator_name` (String) Name of the parent ID Token Validator

### Optional

- `claim_name` (String) The name of the claim to be validated.
- `description` (String) A description for this Token Claim Validation
- `required_value` (String) Specifies the boolean claim's required value.

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


