resource "pingdirectory_synchronization_provider" "mySynchronizationProvider" {
  name    = "MySynchronizationProvider"
  type    = "replication"
  enabled = false
}
