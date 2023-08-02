resource "pingdirectory_default_root_dn_user" "defaultRootDnUser" {
  name              = "Directory Manager"
  alternate_bind_dn = [var.root_user_dn]
}

resource "pingdirectory_root_dn_user" "syncRootDnUser" {
  name                            = "pingdatasync"
  alternate_bind_dn               = ["cn=sync", "cn=datasync"]
  inherit_default_root_privileges = false
  privilege                       = ["bypass-acl", "bypass-pw-policy", "config-read", "password-reset", "unindexed-search"]
  password                        = var.root_user_password
}

resource "pingdirectory_root_dn_user" "federateRootDnUser" {
  name                            = "pingfederate"
  alternate_bind_dn               = ["cn=fed", "cn=pf", "cn=pingfederate"]
  inherit_default_root_privileges = false
  privilege                       = ["password-reset", "permit-get-password-policy-state-issues", "unindexed-search", "config-read", "proxied-auth"]
  is_proxyable                    = "prohibited"
  password                        = var.root_user_password
}

resource "pingdirectory_root_dn_user" "proxyRootDnUser" {
  name              = "pingdirectoryproxy"
  alternate_bind_dn = ["cn=pingdirectoryproxy", "cn=proxy"]
  privilege         = ["proxied-auth"]
  password          = var.root_user_password
}

resource "pingdirectory_root_dn_user" "authorizeRootDnUser" {
  name                            = "pingauthorize"
  alternate_bind_dn               = ["cn=pingauthorize", "cn=pingdatagovernance", "cn=datagov"]
  inherit_default_root_privileges = false
  privilege                       = ["password-reset", "proxied-auth", "unindexed-search"]
  search_result_entry_limit       = 100000
  password                        = var.root_user_password
}
