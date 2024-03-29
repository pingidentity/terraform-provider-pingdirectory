---
page_title: "pingdirectory_token_claim_validation Resource - terraform-provider-pingdirectory"
subcategory: "Token Claim Validation"
description: |-
  Manages a Token Claim Validation.
---

# pingdirectory_token_claim_validation (Resource)

Manages a Token Claim Validation.

A Token Claim Validation defines a condition about a token claim that must be satisfied for the token to be considered valid.

## Example Usage

```terraform
resource "pingdirectory_id_token_validator" "myPingOneIdTokenValidator" {
  type                   = "ping-one"
  name                   = "MyPingOneIdTokenValidator"
  issuer_url             = "example.com"
  enabled                = false
  identity_mapper        = "Exact Match"
  evaluation_order_index = 1
}

resource "pingdirectory_token_claim_validation" "myTokenClaimValidation" {
  name                    = "MyTokenClaimValidation"
  id_token_validator_name = pingdirectory_id_token_validator.myPingOneIdTokenValidator.id
  any_required_value      = ["my_example_value"]
  claim_name              = "my_example_claim_name"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `claim_name` (String) The name of the claim to be validated.
- `id_token_validator_name` (String) Name of the parent ID Token Validator
- `name` (String) Name of this config object.
- `type` (String) The type of Token Claim Validation resource. Options are ['string-array', 'boolean', 'string']

### Optional

- `all_required_value` (Set of String) The set of all values that the claim must have to be considered valid.
- `any_required_value` (Set of String) The set of values that the claim may have to be considered valid.
- `description` (String) A description for this Token Claim Validation
- `required_value` (String) Specifies the boolean claim's required value.

### Read-Only

- `id` (String) The ID of this resource.
- `notifications` (Set of String) Notifications returned by the PingDirectory Configuration API.
- `required_actions` (Set of Object) Required actions returned by the PingDirectory Configuration API. (see [below for nested schema](#nestedatt--required_actions))

<a id="nestedatt--required_actions"></a>
### Nested Schema for `required_actions`

Read-Only:

- `property` (String)
- `synopsis` (String)
- `type` (String)

## Import

Import is supported using the following syntax:

```shell
# Importing a Token Claim Validation requires providing the name of all parent resources in the following format
terraform import pingdirectory_token_claim_validation.myTokenClaimValidation "[id-token-validator-name]/[token-claim-validation-name]"
```

