resource "pingdirectory_velocity_context_provider" "myVelocityContextProvider" {
  name                        = "MyVelocityContextProvider"
  http_servlet_extension_name = "Velocity"
  type                        = "velocity-tools"
  included_view               = ["path/to/view1", "path/to/view2"]
}
