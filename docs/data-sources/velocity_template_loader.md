---
page_title: "pingdirectory_velocity_template_loader Data Source - terraform-provider-pingdirectory"
subcategory: "Velocity Template Loader"
description: |-
  Describes a Velocity Template Loader.
---

# pingdirectory_velocity_template_loader (Data Source)

Describes a Velocity Template Loader.

Velocity Template Loaders load templates from the filesystem.

## Example Usage

```terraform
data "pingdirectory_velocity_template_loader" "myVelocityTemplateLoader" {
  name                        = "MyVelocityTemplateLoader"
  http_servlet_extension_name = "MyHttpServletExtension"
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_met_support_mult_content_types)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `http_servlet_extension_name` (String) Name of the parent HTTP Servlet Extension
- `name` (String) Name of this config object.

### Read-Only

- `enabled` (Boolean) Indicates whether this Velocity Template Loader is enabled.
- `evaluation_order_index` (Number) This property determines the evaluation order for determining the correct Velocity Template Loader to load a template for generating content for a particular request.
- `id` (String) The ID of this resource.
- `mime_type` (String) Specifies a the value that will be used in the response's Content-Type header that indicates the type of content to return.
- `mime_type_matcher` (String) Specifies a media type for matching Accept request-header values.
- `template_directory` (String) Specifies the directory in which to search for the template files.
- `template_suffix` (String) Specifies the suffix to append to the requested resource name when searching for the template file with which to form a response.
- `type` (String) The type of Velocity Template Loader resource. Options are ['velocity-template-loader']

