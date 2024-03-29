---
page_title: "pingdirectory_password_validators Data Source - terraform-provider-pingdirectory"
subcategory: "Password Validator"
description: |-
  Lists Password Validator objects in the server configuration.
---

# pingdirectory_password_validators (Data Source)

Lists Password Validator objects in the server configuration.

Password Validators are responsible for determining whether a proposed password is acceptable for use and could include checks like ensuring it meets minimum length requirements, that it has an appropriate range of characters, or that it is not in the history.

## Example Usage

```terraform
data "pingdirectory_password_validators" "list" {
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_ds_config_password_validators)

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (String) SCIM filter used when searching the configuration.

### Read-Only

- `id` (String) The ID of this resource.
- `objects` (Set of Object) Password Validator objects found in the configuration (see [below for nested schema](#nestedatt--objects))

<a id="nestedatt--objects"></a>
### Nested Schema for `objects`

Read-Only:

- `id` (String)
- `type` (String)

