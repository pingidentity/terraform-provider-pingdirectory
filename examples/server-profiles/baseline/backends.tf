resource "pingdirectory_default_backend" "changelogBackend" {
  backend_id            = "changelog"
  enabled               = true
  changelog_maximum_age = "2 h"
}

resource "pingdirectory_default_backend" "defaultUserRootBackend" {
  backend_id               = "userRoot"
  compact_common_parent_dn = ["ou=people,${var.user_base_dn}"]
}
