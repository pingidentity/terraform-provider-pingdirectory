resource "pingdirectory_root_dn_user" "myRootDnUser" {
  name                            = "MyRootDnUser"
  inherit_default_root_privileges = true
  search_result_entry_limit       = 0
  time_limit_seconds              = 0
  look_through_entry_limit        = 0
  idle_time_limit_seconds         = 0
  password_policy                 = "Root Password Policy"
  require_secure_authentication   = false
  require_secure_connections      = false
}
