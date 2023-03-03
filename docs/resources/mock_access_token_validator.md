---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_mock_access_token_validator Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Mock Access Token Validator.
---

# pingdirectory_mock_access_token_validator (Resource)

Manages a Mock Access Token Validator.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `enabled` (Boolean) Indicates whether this Access Token Validator is enabled for use in Directory Server.
- `id` (String) Name of this object.

### Optional

- `client_id_claim_name` (String) The name of the token claim that contains the OAuth2 client ID.
- `description` (String) A description for this Access Token Validator
- `evaluation_order_index` (Number) When multiple Mock Access Token Validators are defined for a single Directory Server, this property determines the evaluation order for determining the correct validator class for an access token received by the Directory Server. Values of this property must be unique among all Mock Access Token Validators defined within Directory Server but not necessarily contiguous. Mock Access Token Validators with a smaller value will be evaluated first to determine if they are able to validate the access token.
- `identity_mapper` (String) Specifies the name of the Identity Mapper that should be used for associating user entries with Bearer token subject names. The claim name from which to obtain the subject (i.e. the currently logged-in user) may be configured using the subject-claim-name property.
- `scope_claim_name` (String) The name of the token claim that contains the scopes granted by the token.
- `subject_claim_name` (String) The name of the token claim that contains the subject, i.e. the logged-in user in an access token. This property goes hand-in-hand with the identity-mapper property and tells the Identity Mapper which field to use to look up the user entry on the server.

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

