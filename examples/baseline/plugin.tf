resource "pingdirectory_purge_expired_data_plugin" "expiredSessionAutoPurgePlugin" {
  id                 = "ExpiredSessionAutoPurge"
  enabled            = true
  datetime_attribute = "pf-authn-session-group-expiry-time"
  expiration_offset  = "1 h"
  purge_behavior     = "subtree-delete-entries"
  base_dn            = "ou=sessions,${var.user_base_dn}"
  filter             = "(objectClass=pf-authn-session-groups)"
  polling_interval   = "20 m"
}

resource "pingdirectory_purge_expired_data_plugin" "idleSessionAutoPurgePlugin" {
  id                 = "IdleSessionAutoPurge"
  enabled            = true
  datetime_attribute = "pf-authn-session-group-last-activity-time"
  expiration_offset  = "1 w"
  purge_behavior     = "subtree-delete-entries"
  base_dn            = "ou=sessions,${var.user_base_dn}"
  filter             = "(objectClass=pf-authn-session-groups)"
  polling_interval   = "1 d"
}
