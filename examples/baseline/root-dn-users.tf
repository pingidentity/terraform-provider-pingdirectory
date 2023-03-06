resource "pingdirectory_default_root_dn_user" "defaultRootDnUser" {
  id                = "Directory Manager"
  alternate_bind_dn = ["cn=administrator"]
}

resource "pingdirectory_root_dn_user" "syncRootDnUser" {
  id                              = "pingdatasync"
  alternate_bind_dn               = ["cn=sync", "cn=datasync"]
  inherit_default_root_privileges = false
  privilege                       = ["bypass-acl", "bypass-pw-policy", "config-read", "password-reset", "unindexed-search"]
  password                        = "2FederateM0re"
}

resource "pingdirectory_root_dn_user" "federateRootDnUser" {
  id                              = "pingfederate"
  alternate_bind_dn               = ["cn=fed", "cn=pf", "cn=pingfederate"]
  inherit_default_root_privileges = false
  privilege                       = ["password-reset", "permit-get-password-policy-state-issues", "unindexed-search", "config-read", "proxied-auth"]
  is_proxyable                    = "prohibited"
  password                        = "2FederateM0re"
}

resource "pingdirectory_root_dn_user" "proxyRootDnUser" {
  id                = "pingdirectoryproxy"
  alternate_bind_dn = ["cn=pingdirectoryproxy", "cn=proxy"]
  privilege         = ["proxied-auth"]
  password          = "2FederateM0re"
}

resource "pingdirectory_root_dn_user" "authorizeRootDnUser" {
  id                              = "pingauthorize"
  alternate_bind_dn               = ["cn=pingauthorize", "cn=pingdatagovernance", "cn=datagov"]
  inherit_default_root_privileges = false
  privilege                       = ["password-reset", "proxied-auth", "unindexed-search"]
  search_result_entry_limit       = 100000
  password                        = "2FederateM0re"
}
