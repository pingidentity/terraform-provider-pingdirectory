---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_default_http_external_server Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Http External Server.
---

# pingdirectory_default_http_external_server (Resource)

Manages a Http External Server.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `base_url` (String) The base URL of the external server, optionally including port number, for example "https://externalService:9031".
- `connect_timeout` (String) Specifies the maximum length of time to wait for a connection to be established before aborting a request to the server.
- `description` (String) A description for this External Server
- `hostname_verification_method` (String) The mechanism for checking if the hostname of the HTTP External Server matches the name(s) stored inside the server's X.509 certificate. This is only applicable if SSL is being used for connection security.
- `key_manager_provider` (String) The key manager provider to use if SSL (HTTPS) is to be used for connection-level security. When specifying a value for this property (except when using the Null key manager provider) you must ensure that the external server trusts this server's public certificate by adding this server's public certificate to the external server's trust store.
- `response_timeout` (String) Specifies the maximum length of time to wait for response data to be read from an established connection before aborting a request to the server.
- `ssl_cert_nickname` (String) The certificate alias within the keystore to use if SSL (HTTPS) is to be used for connection-level security. When specifying a value for this property you must ensure that the external server trusts this server's public certificate by adding this server's public certificate to the external server's trust store.
- `trust_manager_provider` (String) The trust manager provider to use if SSL (HTTPS) is to be used for connection-level security.

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


