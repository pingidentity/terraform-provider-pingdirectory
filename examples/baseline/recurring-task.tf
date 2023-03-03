resource "pingdirectory_default_ldif_export_recurring_task" "defaultExportAllBackendsRecurringTask" {
    id = "Export All Non-Administrative Backends"
    description = "Exports the contents of all non-administrative backends to LDIF.  The LDIF exports will be gzip-compressed, and they will be encrypted if the global configuration is set to encrypt LDIF exports by default.  The export will be rate-limited to ten megabytes per second of compressed data to help avoid any noticeable impact on server performance.  The exports will be written into the \"/opt/backup\" directory of the container and will be kept for seven days."
    ldif_directory = "/opt/backup"
    encrypt = true
}
