data "pingdirectory_inter_server_authentication_info_list" "list" {
  server_instance_listener_name = "MyServerInstanceListener"
  server_instance_name          = "MyServerInstance"
}
