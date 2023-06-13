resource "pingdirectory_plugin" "expiredSessionAutoPurgePlugin" {
  id                 = "ExpiredSessionAutoPurge"
  resource_type               = "purge-expired-data"
  enabled            = true
  datetime_attribute = "pf-authn-session-group-expiry-time"
  expiration_offset  = "1 h"
  purge_behavior     = "subtree-delete-entries"
  base_dn            = ["ou=sessions,${var.user_base_dn}"]
  filter             = ["(objectClass=pf-authn-session-groups)"]
  polling_interval   = "20 m"
}

resource "pingdirectory_plugin" "idleSessionAutoPurgePlugin" {
  resource_type               = "purge-expired-data"
  id                 = "IdleSessionAutoPurge"
  enabled            = true
  datetime_attribute = "pf-authn-session-group-last-activity-time"
  expiration_offset  = "1 w"
  purge_behavior     = "subtree-delete-entries"
  base_dn            = ["ou=sessions,${var.user_base_dn}"]
  filter             = ["(objectClass=pf-authn-session-groups)"]
  polling_interval   = "1 d"
}
