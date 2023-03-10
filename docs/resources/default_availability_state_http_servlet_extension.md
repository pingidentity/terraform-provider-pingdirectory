---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_default_availability_state_http_servlet_extension Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Availability State Http Servlet Extension.
---

# pingdirectory_default_availability_state_http_servlet_extension (Resource)

Manages a Availability State Http Servlet Extension.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `additional_response_contents` (String) A JSON-formatted string containing additional fields to be returned in the response body. For example, an additional-response-contents value of '{ "key": "value" }' would add the key and value to the root of the JSON response body.
- `available_status_code` (Number) Specifies the HTTP status code that the servlet should return if the server considers itself to be available.
- `base_context_path` (String) Specifies the base context path that HTTP clients should use to access this servlet. The value must start with a forward slash and must represent a valid HTTP context path.
- `correlation_id_response_header` (String) Specifies the name of the HTTP response header that will contain a correlation ID value. Example values are "Correlation-Id", "X-Amzn-Trace-Id", and "X-Request-Id".
- `cross_origin_policy` (String) The cross-origin request policy to use for the HTTP Servlet Extension.
- `degraded_status_code` (Number) Specifies the HTTP status code that the servlet should return if the server considers itself to be degraded.
- `description` (String) A description for this HTTP Servlet Extension
- `include_response_body` (Boolean) Indicates whether the response should include a body that is a JSON object.
- `override_status_code` (Number) Specifies a HTTP status code that the servlet should always return, regardless of the server's availability. If this value is defined, it will override the availability-based return codes.
- `response_header` (Set of String) Specifies HTTP header fields and values added to response headers for all requests.
- `unavailable_status_code` (Number) Specifies the HTTP status code that the servlet should return if the server considers itself to be unavailable.

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


