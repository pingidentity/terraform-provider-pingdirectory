resource "pingdirectory_default_root_dn" "myRootDn" {
  default_root_privilege_name = ["bypass-acl", "config-read", "config-write", "modify-acl", "privilege-change", "use-admin-session"]
}
