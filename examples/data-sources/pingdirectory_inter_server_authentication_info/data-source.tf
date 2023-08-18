data "pingdirectory_inter_server_authentication_info" "myInterServerAuthenticationInfo" {
  name                          = "MyInterServerAuthenticationInfo"
  server_instance_listener_name = "MyServerInstanceListener"
  server_instance_name          = "MyServerInstance"
}
