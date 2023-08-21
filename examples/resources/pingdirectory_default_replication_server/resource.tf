resource "pingdirectory_default_replication_server" "myReplicationServer" {
  synchronization_provider_name = "Multimaster Synchronization"
  replication_server_id         = 26554
  replication_port              = 8989
  gateway_priority              = 5
}
