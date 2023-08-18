resource "pingdirectory_velocity_template_loader" "myVelocityTemplateLoader" {
  name                        = "MyVelocityTemplateLoader"
  evaluation_order_index      = 10100
  http_servlet_extension_name = "Velocity"
  mime_type_matcher           = "text/html"
  template_suffix             = ".vm"
}
