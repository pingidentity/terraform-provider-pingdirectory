---
page_title: "pingdirectory_default_http_servlet_extension Resource - terraform-provider-pingdirectory"
subcategory: "Http Servlet Extension"
description: |-
  Manages a Http Servlet Extension.
---

# pingdirectory_default_http_servlet_extension (Resource)

Manages a Http Servlet Extension.

HTTP Servlet Extensions may be used to define classes and initialization parameters that should be used for a servlet invoked by an HTTP connection handler.

Since this is a 'default' resource, the managed object must already exist in the PingDirectory configuration.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of this config object.

### Optional

- `access_token_scope` (String) When the `type` attribute is set to:
  - `delegated-admin`: The name of a scope that must be present in an access token accepted by the Delegated Admin HTTP Servlet Extension.
  - `directory-rest-api`: The name of a scope that must be present in an access token accepted by the Directory REST API HTTP Servlet Extension.
- `access_token_validator` (Set of String) When the `type` attribute is set to:
  - `delegated-admin`: If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this Delegated Admin HTTP Servlet Extension.
  - `consent`: If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this Consent HTTP Servlet Extension.
  - `file-server`: The access token validators that may be used to verify the authenticity of an OAuth 2.0 bearer token.
  - `scim2`: If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this SCIM2 HTTP Servlet Extension.
  - `directory-rest-api`: If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this Directory REST API HTTP Servlet Extension.
- `additional_response_contents` (String) A JSON-formatted string containing additional fields to be returned in the response body. For example, an additional-response-contents value of '{ "key": "value" }' would add the key and value to the root of the JSON response body.
- `allow_context_override` (Boolean) Indicates whether context providers may override existing context objects with new values.
- `allowed_authentication_type` (Set of String) The types of authentication that may be used to authenticate to the file servlet.
- `allowed_control` (Set of String) Specifies the names of any request controls that should be allowed by the Directory REST API. Any request that contains a critical control not in this list will be rejected. Any non-critical request control which is not supported by the Directory REST API will be removed from the request.
- `always_include_monitor_entry_name_label` (Boolean) Indicates whether generated metrics should always include a "monitor_entry" label whose value is the name of the monitor entry from which the metric was obtained.
- `always_use_permissive_modify` (Boolean) Indicates whether to always use permissive modify behavior for PATCH requests, even if the request did not include the permissive modify request control.
- `audience` (String) When the `type` attribute is set to:
  - `delegated-admin`: A string or URI that identifies the Delegated Admin HTTP Servlet Extension in the context of OAuth2 authorization.
  - `directory-rest-api`: A string or URI that identifies the Directory REST API HTTP Servlet Extension in the context of OAuth2 authorization.
- `available_status_code` (Number) Specifies the HTTP status code that the servlet should return if the server considers itself to be available.
- `base_context_path` (String) When the `type` attribute is set to:
  - One of [`availability-state`, `prometheus-monitoring`]: Specifies the base context path that HTTP clients should use to access this servlet. The value must start with a forward slash and must represent a valid HTTP context path.
  - `velocity`: The context path to use to access all template-based and static content. The value must start with a forward slash and must represent a valid HTTP context path.
  - `ldap-mapped-scim`: The context path to use to access the SCIM interface. The value must start with a forward slash and must represent a valid HTTP context path.
  - `file-server`: Specifies the base context path that should be used by HTTP clients to reference content. The value must start with a forward slash and must represent a valid HTTP context path.
  - `scim2`: The context path to use to access the SCIM 2.0 interface. The value must start with a forward slash and must represent a valid HTTP context path.
- `basic_auth_enabled` (Boolean) When the `type` attribute is set to:
  - One of [`delegated-admin`, `consent`, `directory-rest-api`]: Enables HTTP Basic authentication, using a username and password. The Identity Mapper specified by the identity-mapper property will be used to map the username to a DN.  NOTE: Basic authentication is considered less secure than OAuth2 bearer token authentication.
  - `ldap-mapped-scim`: Enables HTTP Basic authentication, using a username and password.  NOTE: Basic authentication is considered less secure than OAuth2 bearer token authentication.
- `bearer_token_auth_enabled` (Boolean) Enables HTTP bearer token authentication.
- `bulk_max_concurrent_requests` (Number) The maximum number of bulk requests that may be processed concurrently by the server. Any bulk request that would cause this limit to be exceeded is rejected with HTTP status code 503.
- `bulk_max_operations` (Number) The maximum number of operations that are permitted in a bulk request.
- `bulk_max_payload_size` (String) The maximum payload size in bytes of a bulk request.
- `character_encoding` (String) Specifies the value that will be used for all responses' Content-Type headers' charset parameter that indicates the character encoding of the document.
- `correlation_id_response_header` (String) Specifies the name of the HTTP response header that will contain a correlation ID value. Example values are "Correlation-Id", "X-Amzn-Trace-Id", and "X-Request-Id".
- `cross_origin_policy` (String) The cross-origin request policy to use for the HTTP Servlet Extension.
- `debug_enabled` (Boolean) When the `type` attribute is set to:
  - `ldap-mapped-scim`: Enables debug logging of the SCIM SDK. Debug messages will be forwarded to the Directory Server debug logger with the scope of com.unboundid.directory.server.extensions.scim.SCIMHTTPServletExtension.
  - `scim2`: Enables debug logging of the SCIM 2.0 SDK. Debug messages will be forwarded to the Directory Server debug logger with the scope of com.unboundid.directory.broker.http.scim2.extension.SCIM2HTTPServletExtension.
- `debug_level` (String) The minimum debug level that should be used for messages to be logged.
- `debug_type` (Set of String) The types of debug messages that should be logged.
- `default_mime_type` (String) When the `type` attribute is set to:
  - `velocity`: Specifies the default value that will be used in the response's Content-Type header that indicates the type of content to return.
  - `file-server`: Specifies the default MIME type to use for the Content-Type header when a mapping cannot be found.
- `default_operational_attribute` (Set of String) A set of operational attributes that will be returned with entries by default.
- `degraded_status_code` (Number) Specifies the HTTP status code that the servlet should return if the server considers itself to be degraded.
- `description` (String) A description for this HTTP Servlet Extension
- `document_root_directory` (String) Specifies the path to the directory on the local filesystem containing the files to be served by this File Server HTTP Servlet Extension. The path must exist, and it must be a directory.
- `enable_directory_indexing` (Boolean) Indicates whether to generate a default HTML page with a listing of available files if the requested path refers to a directory rather than a file, and that directory does not contain an index file.
- `entity_tag_ldap_attribute` (String) Specifies the LDAP attribute whose value should be used as the entity tag value to enable SCIM resource versioning support.
- `exclude_ldap_base_dn` (Set of String) Specifies the base DNs for the branches of the DIT that should not be exposed via the Identity Access API.
- `exclude_ldap_objectclass` (Set of String) Specifies the LDAP object classes that should be not be exposed directly as SCIM resources.
- `expose_request_attributes` (Boolean) Specifies whether the HTTP request will be exposed to templates.
- `expose_server_context` (Boolean) Specifies whether a server context will be exposed under context key 'ubid_server' for all template contexts.
- `expose_session_attributes` (Boolean) Specifies whether the HTTP session will be exposed to templates.
- `extension_argument` (Set of String) The set of arguments used to customize the behavior for the Third Party HTTP Servlet Extension. Each configuration property should be given in the form 'name=value'.
- `extension_class` (String) The fully-qualified name of the Java class providing the logic for the Third Party HTTP Servlet Extension.
- `id_token_validator` (Set of String) The ID token validators that may be used to verify the authenticity of an of an OpenID Connect ID token.
- `identity_mapper` (String) When the `type` attribute is set to:
  - `delegated-admin`: Specifies the Identity Mapper that is to be used for associating user entries with basic authentication user names.
  - `velocity`: Specifies the name of the identity mapper that is to be used for associating basic authentication credentials with user entries.
  - `consent`: Specifies the Identity Mapper that is to be used for associating basic authentication usernames with DNs.
  - `ldap-mapped-scim`: Specifies the name of the identity mapper that is to be used to match the username included in the HTTP Basic authentication header to the corresponding user in the directory.
  - `file-server`: The identity mapper that will be used to identify the entry with which a username is associated.
  - `config`: Specifies the name of the identity mapper that is to be used for associating user entries with basic authentication user names.
  - `directory-rest-api`: Specifies the Identity Mapper that is to be used for associating user entries with basic authentication usernames.
- `include_instance_name_label` (Boolean) Indicates whether generated metrics should include an "instance" label whose value is the instance name for this Directory Server instance.
- `include_ldap_base_dn` (Set of String) Specifies the base DNs for the branches of the DIT that should be exposed via the Identity Access API.
- `include_ldap_objectclass` (Set of String) Specifies the LDAP object classes that should be exposed directly as SCIM resources.
- `include_location_name_label` (Boolean) Indicates whether generated metrics should include a "location" label whose value is the location name for this Directory Server instance.
- `include_monitor_attribute_name_label` (Boolean) Indicates whether generated metrics should include a "monitor_attribute" label whose value is the name of the monitor attribute from which the metric was obtained.
- `include_monitor_object_class_name_label` (Boolean) Indicates whether generated metrics should include a "monitor_object_class" label whose value is the name of the object class for the monitor entry from which the metric was obtained.
- `include_product_name_label` (Boolean) Indicates whether generated metrics should include a "product" label whose value is the product name for this Directory Server instance.
- `include_response_body` (Boolean) Indicates whether the response should include a body that is a JSON object.
- `include_stack_trace` (Boolean) Indicates whether a stack trace of the thread which called the debug method should be included in debug log messages.
- `index_file` (Set of String) Specifies the name of a file whose contents may be returned to the client if the requested path refers to a directory rather than a file.
- `label_name_value_pair` (Set of String) A set of name-value pairs for labels that should be included in all metrics exposed by this Directory Server instance.
- `map_access_tokens_to_local_users` (String) Indicates whether the SCIM2 servlet should attempt to map the presented access token to a local user.
- `max_page_size` (Number) The maximum number of entries to be returned in one page of search results.
- `max_results` (Number) The maximum number of resources that are returned in a response.
- `mime_types_file` (String) When the `type` attribute is set to:
  - `velocity`: Specifies the path to a file that contains MIME type mappings that will be used to determine the appropriate value to return for the Content-Type header based on the extension of the requested static content file.
  - `file-server`: Specifies the path to a file that contains MIME type mappings that will be used to determine the appropriate value to return for the Content-Type header based on the extension of the requested file.
- `oauth_token_handler` (String) Specifies the OAuth Token Handler implementation that should be used to validate OAuth 2.0 bearer tokens when they are included in a SCIM request.
- `override_status_code` (Number) Specifies a HTTP status code that the servlet should always return, regardless of the server's availability. If this value is defined, it will override the availability-based return codes.
- `reject_expansion_attribute` (Set of String) A set of attributes which the client is not allowed to provide for the expand query parameters. This should be used for attributes that could either have a large number of values or that reference entries that are very large like groups.
- `require_authentication` (Boolean) When the `type` attribute is set to:
  - `velocity`: Require authentication when accessing Velocity templates.
  - `file-server`: Indicates whether the servlet extension should only accept requests from authenticated clients.
- `require_file_servlet_access_privilege` (Boolean) Indicates whether the servlet extension should only accept requests from authenticated clients that have the file-servlet-access privilege.
- `require_group` (Set of String) The DN of a group whose members will be permitted to access to the associated files. If multiple group DNs are configured, then anyone who is a member of at least one of those groups will be granted access.
- `resource_mapping_file` (String) The path to an XML file defining the resources supported by the SCIM interface and the SCIM-to-LDAP attribute mappings to use.
- `response_header` (Set of String) When the `type` attribute is set to:
  - One of [`delegated-admin`, `quickstart`, `availability-state`, `prometheus-monitoring`, `consent`, `ldap-mapped-scim`, `groovy-scripted`, `file-server`, `config`, `scim2`, `directory-rest-api`, `third-party`]: Specifies HTTP header fields and values added to response headers for all requests.
  - `velocity`: Specifies HTTP header fields and values added to response headers for all template page requests.
- `schemas_endpoint_objectclass` (Set of String) The list of object classes which will be returned by the schemas endpoint.
- `script_argument` (Set of String) The set of arguments used to customize the behavior for the Scripted HTTP Servlet Extension. Each configuration property should be given in the form 'name=value'.
- `script_class` (String) The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted HTTP Servlet Extension.
- `server` (String) Specifies the PingFederate server to be configured.
- `static_content_directory` (String) Specifies the base directory in which static, non-template content such as images, CSS, and Javascript files are stored on the filesystem.
- `static_context_path` (String) The path below the base context path by which static, non-template content such as images, CSS, and Javascript files are accessible.
- `static_custom_directory` (String) Specifies the base directory in which custom static, non-template content such as images, CSS, and Javascript files are stored on the filesystem. Files in this directory will override those with the same name in the directory specified by the static-content-directory property.
- `static_response_header` (Set of String) Specifies HTTP header fields and values added to response headers for static content requests such as images and scripts.
- `swagger_enabled` (Boolean) Indicates whether the SCIM2 HTTP Servlet Extension will generate a Swagger specification document.
- `template_directory` (Set of String) Specifies an ordered list of directories in which to search for the template files.
- `temporary_directory` (String) Specifies the location of the directory that is used to create temporary files containing SCIM request data.
- `temporary_directory_permissions` (String) Specifies the permissions that should be applied to the directory that is used to create temporary files.
- `unavailable_status_code` (Number) Specifies the HTTP status code that the servlet should return if the server considers itself to be unavailable.

### Read-Only

- `id` (String) The ID of this resource.
- `notifications` (Set of String) Notifications returned by the PingDirectory Configuration API.
- `required_actions` (Set of Object) Required actions returned by the PingDirectory Configuration API. (see [below for nested schema](#nestedatt--required_actions))
- `type` (String) The type of HTTP Servlet Extension resource. Options are ['delegated-admin', 'quickstart', 'availability-state', 'prometheus-monitoring', 'velocity', 'consent', 'ldap-mapped-scim', 'groovy-scripted', 'file-server', 'config', 'scim2', 'directory-rest-api', 'third-party']

<a id="nestedatt--required_actions"></a>
### Nested Schema for `required_actions`

Read-Only:

- `property` (String)
- `synopsis` (String)
- `type` (String)



