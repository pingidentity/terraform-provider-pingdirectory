resource "pingdirectory_default_file_based_debug_log_publisher" "defaultFileBasedDebugLogger" {
  id      = "File-Based Debug Logger"
  enabled = true
}

resource "pingdirectory_debug_target" "debugTargetReplicationCli" {
    log_publisher_name = "File-Based Debug Logger"
  debug_scope = "com.unboundid.guitools.replicationcli"
  debug_level = "verbose"
}

resource "pingdirectory_debug_target" "debugTargetFooBar" {
    log_publisher_name = "File-Based Debug Logger"
  debug_scope = "foo.bar"
  debug_level = "verbose"
}
