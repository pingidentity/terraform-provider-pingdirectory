data "pingdirectory_replication_domain" "myReplicationDomain" {
  name                          = "MyReplicationDomain"
  synchronization_provider_name = "MySynchronizationProvider"
}
