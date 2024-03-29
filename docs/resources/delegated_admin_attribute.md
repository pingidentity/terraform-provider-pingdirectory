---
page_title: "pingdirectory_delegated_admin_attribute Resource - terraform-provider-pingdirectory"
subcategory: "Delegated Admin Attribute"
description: |-
  Manages a Delegated Admin Attribute.
---

# pingdirectory_delegated_admin_attribute (Resource)

Manages a Delegated Admin Attribute.

A Delegated Admin Attribute defines an LDAP attribute which can be accessed by a client of the Delegated Admin API.

## Example Usage

```terraform
resource "pingdirectory_delegated_admin_attribute" "myDelegatedAdminAttribute" {
  rest_resource_type_name = "MyRestResourceType"
  type                    = "certificate"
  attribute_type          = "myattr"
  display_name            = "MyAttribute"
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_da_config_delegated_admin)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `attribute_type` (String) Specifies the name or OID of the LDAP attribute type.
- `display_name` (String) A human readable display name for this Delegated Admin Attribute.
- `rest_resource_type_name` (String) Name of the parent REST Resource Type
- `type` (String) The type of Delegated Admin Attribute resource. Options are ['certificate', 'photo', 'generic']

### Optional

- `allowed_mime_type` (Set of String) The list of file types allowed to be uploaded. If no types are specified, then all types will be allowed.
- `attribute_category` (String) Specifies which attribute category this attribute belongs to.
- `attribute_presentation` (String) Indicates how the attribute is presented to the user of the app.
- `date_time_format` (String) Specifies the format string that is used to present a date and/or time value to the user of the app. This property only applies to LDAP attribute types whose LDAP syntax is GeneralizedTime and is ignored if the attribute type has any other syntax.
- `description` (String) A description for this Delegated Admin Attribute
- `display_order_index` (Number) This property determines a display order for attributes within a given attribute category. Attributes are ordered within their category based on this index from least to greatest.
- `include_in_summary` (Boolean) Indicates whether this Delegated Admin Attribute is to be included in the summary display for a resource.
- `multi_valued` (Boolean) Indicates whether this Delegated Admin Attribute may have multiple values.
- `mutability` (String) Specifies the circumstances under which the values of the attribute can be written.
- `reference_resource_type` (String) For LDAP attributes with DN syntax, specifies what kind of resource is referenced.

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
# Importing a Delegated Admin Attribute requires providing the name of all parent resources in the following format
terraform import pingdirectory_delegated_admin_attribute.myDelegatedAdminAttribute "[rest-resource-type-name]/[delegated-admin-attribute-attribute-type]"
```

