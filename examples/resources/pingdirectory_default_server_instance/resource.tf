resource "pingdirectory_default_directory_server_instance" "myServerInstance" {
  name                 = "MyServerInstance"
  server_instance_name = "MyDirectoryServerInstance"
  server_version       = "10.2.0.0"
}
