resource "pingdirectory_topology_admin_user" "myTopologyAdminUser" {
  name                            = "MyTopologyAdminUser"
  inherit_default_root_privileges = true
  search_result_entry_limit       = 100
  time_limit_seconds              = 60
  look_through_entry_limit        = 20
  idle_time_limit_seconds         = 120
  password_policy                 = "Default Password Policy"
  require_secure_authentication   = true
  require_secure_connections      = false
}
