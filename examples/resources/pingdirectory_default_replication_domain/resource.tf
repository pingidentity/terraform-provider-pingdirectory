resource "pingdirectory_default_replication_domain" "myReplicationDomain" {
  name                          = "MyReplicationDomain"
  synchronization_provider_name = "Multimaster Synchronization"
  server_id                     = 1234
  base_dn                       = "dc=example,dc=com"
}
