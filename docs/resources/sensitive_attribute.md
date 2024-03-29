---
page_title: "pingdirectory_sensitive_attribute Resource - terraform-provider-pingdirectory"
subcategory: "Sensitive Attribute"
description: |-
  Manages a Sensitive Attribute.
---

# pingdirectory_sensitive_attribute (Resource)

Manages a Sensitive Attribute.

Sensitive Attributes provide a means of declaring one or more attributes to contain sensitive data so that the server can enforce additional protection for operations attempting to interact with them.

## Example Usage

```terraform
resource "pingdirectory_sensitive_attribute" "mySensitiveAttribute" {
  name           = "MySensitiveAttribute"
  attribute_type = ["userPassword", "pwdHistory"]
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_ds_config_sensitive_attriutes)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `attribute_type` (Set of String) The name(s) or OID(s) of the attribute types for attributes whose values may be considered sensitive.
- `name` (String) Name of this config object.

### Optional

- `allow_in_add` (String) Indicates whether clients will be allowed to include sensitive attributes in add requests.
- `allow_in_compare` (String) Indicates whether clients will be allowed to target sensitive attributes with compare requests.
- `allow_in_filter` (String) Indicates whether clients will be allowed to include sensitive attributes in search filters. This also includes filters that may be used in other forms, including assertion and LDAP join request controls.
- `allow_in_modify` (String) Indicates whether clients will be allowed to target sensitive attributes with modify requests.
- `allow_in_returned_entries` (String) Indicates whether sensitive attributes should be included in entries returned to the client. This includes not only search result entries, but also other forms including in the values of controls like the pre-read, post-read, get authorization entry, and LDAP join response controls.
- `description` (String) A description for this Sensitive Attribute
- `include_default_sensitive_operational_attributes` (Boolean) Indicates whether to automatically include any server-generated operational attributes that may contain sensitive data.
- `type` (String) The type of Sensitive Attribute resource. Options are ['sensitive-attribute']

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
# "sensitiveAttributeId" should be the name of the Sensitive Attribute to be imported
terraform import pingdirectory_sensitive_attribute.mySensitiveAttribute sensitiveAttributeId
```

