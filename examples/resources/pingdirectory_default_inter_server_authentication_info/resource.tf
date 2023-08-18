resource "pingdirectory_default_inter_server_authentication_info" "myInterServerAuthenticationInfo" {
  name                          = "certificate-auth-mirrored-config"
  server_instance_listener_name = "ldap-listener-mirrored-config"
  server_instance_name          = "instance-name"
  purpose                       = ["mirrored-config"]
}
