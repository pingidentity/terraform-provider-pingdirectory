resource "pingdirectory_monitor_provider" "myMonitorProvider" {
  name            = "MyMonitorProvider"
  type            = "third-party"
  description     = "My third party monitor entry resource provider"
  enabled         = false
  extension_class = "com.Example"
}
