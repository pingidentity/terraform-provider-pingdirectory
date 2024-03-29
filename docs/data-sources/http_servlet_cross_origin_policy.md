---
page_title: "pingdirectory_http_servlet_cross_origin_policy Data Source - terraform-provider-pingdirectory"
subcategory: "Http Servlet Cross Origin Policy"
description: |-
  Describes a Http Servlet Cross Origin Policy.
---

# pingdirectory_http_servlet_cross_origin_policy (Data Source)

Describes a Http Servlet Cross Origin Policy.

This object describes a configuration for handling Cross-Origin HTTP requests using the Cross Origin Resource Sharing (CORS) protocol, as defined at http://www.w3.org/TR/cors. An instance of HTTP Servlet Cross Origin Policy can be associated with zero or more HTTP Servlet Extensions to set the Cross-Origin policy for those servlets.

## Example Usage

```terraform
data "pingdirectory_http_servlet_cross_origin_policy" "myHttpServletCrossOriginPolicy" {
  name = "MyHttpServletCrossOriginPolicy"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of this config object.

### Read-Only

- `cors_allow_credentials` (Boolean) Indicates whether the servlet extension allows CORS requests with username/password credentials.
- `cors_allowed_headers` (Set of String) A list of HTTP headers that are supported by the resource and can be specified in a cross-origin request.
- `cors_allowed_methods` (Set of String) A list of HTTP methods allowed for cross-origin access to resources. i.e. one or more of GET, POST, PUT, DELETE, etc.
- `cors_allowed_origins` (Set of String) A list of origins that are allowed to execute cross-origin requests.
- `cors_exposed_headers` (Set of String) A list of HTTP headers other than the simple response headers that browsers are allowed to access.
- `cors_preflight_max_age` (String) The maximum amount of time that a preflight request can be cached by a client.
- `description` (String) A description for this HTTP Servlet Cross Origin Policy
- `id` (String) The ID of this resource.
- `type` (String) The type of HTTP Servlet Cross Origin Policy resource. Options are ['http-servlet-cross-origin-policy']

