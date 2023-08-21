resource "pingdirectory_default_server_instance_listener" "myServerInstanceListener" {
  name                 = "ldap-listener-mirrored-config"
  server_instance_name = "MyServerInstance"
}
