---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_default_smtp_external_server Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Smtp External Server.
---

# pingdirectory_default_smtp_external_server (Resource)

Manages a Smtp External Server.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `description` (String) A description for this External Server
- `passphrase_provider` (String) The passphrase provider to use to obtain the login password for the specified user.
- `password` (String, Sensitive) The login password for the specified user name. Both username and password must be supplied if this attribute is set.
- `server_host_name` (String) The host name of the smtp server.
- `server_port` (Number) The port number where the smtp server listens for requests.
- `smtp_connection_properties` (Set of String) Specifies the connection properties for the smtp server.
- `smtp_security` (String) This property specifies type of connection security to use when connecting to the outgoing mail server.
- `smtp_timeout` (String) Specifies the maximum length of time that a connection or attempted connection to a SMTP server may take.
- `user_name` (String) The name of the login account to use when connecting to the smtp server. Both username and password must be supplied if this attribute is set.

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


