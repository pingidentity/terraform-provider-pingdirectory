---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_default_regular_expression_identity_mapper Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Regular Expression Identity Mapper.
---

# pingdirectory_default_regular_expression_identity_mapper (Resource)

Manages a Regular Expression Identity Mapper.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `description` (String) A description for this Identity Mapper
- `enabled` (Boolean) Indicates whether the Identity Mapper is enabled for use.
- `match_attribute` (Set of String) Specifies the name or OID of the attribute whose value should match the provided identifier string after it has been processed by the associated regular expression.
- `match_base_dn` (Set of String) Specifies the base DN(s) that should be used when performing searches to map the provided ID string to a user entry. If multiple values are given, searches are performed below all the specified base DNs.
- `match_filter` (String) An optional filter that mapped users must match.
- `match_pattern` (String) Specifies the regular expression pattern that is used to identify portions of the ID string that will be replaced.
- `replace_pattern` (String) Specifies the replacement pattern that should be used for substrings in the ID string that match the provided regular expression pattern.

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

