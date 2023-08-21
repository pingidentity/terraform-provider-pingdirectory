data "pingdirectory_replication_server" "myReplicationServer" {
  synchronization_provider_name = "MySynchronizationProvider"
  replication_server_id         = 1234
}
