resource "pingdirectory_plugin" "myPlugin" {
  name    = "MyPlugin"
  resource_type    = "third-party"
  enabled = false
  extension_class = "com.Example"
}
