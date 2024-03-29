---
page_title: "pingdirectory_consent_definition Resource - terraform-provider-pingdirectory"
subcategory: "Consent Definition"
description: |-
  Manages a Consent Definition.
---

# pingdirectory_consent_definition (Resource)

Manages a Consent Definition.

A Consent Definition represents a type of consent to share data.

## Example Usage

```terraform
resource "pingdirectory_consent_definition" "myConsentDefinition" {
  unique_id    = "myConsentDefinition"
  display_name = "example display name"
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_cs_create_consent_def_localization)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `unique_id` (String) A version-independent unique identifier for this Consent Definition.

### Optional

- `description` (String) A description for this Consent Definition
- `display_name` (String) A human-readable display name for this Consent Definition.
- `parameter` (Set of String) Optional parameters for this Consent Definition.
- `type` (String) The type of Consent Definition resource. Options are ['consent-definition']

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
# "consentDefinitionId" should be the unique_id of the Consent Definition to be imported
terraform import pingdirectory_consent_definition.myConsentDefinition consentDefinitionId
```

