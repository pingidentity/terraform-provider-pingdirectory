---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_default_subject_equals_dn_certificate_mapper Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Subject Equals Dn Certificate Mapper.
---

# pingdirectory_default_subject_equals_dn_certificate_mapper (Resource)

Manages a Subject Equals Dn Certificate Mapper.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `description` (String) A description for this Certificate Mapper
- `enabled` (Boolean) Indicates whether the Certificate Mapper is enabled.

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

