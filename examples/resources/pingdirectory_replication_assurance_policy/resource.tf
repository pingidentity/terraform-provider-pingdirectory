resource "pingdirectory_replication_assurance_policy" "myReplicationAssurancePolicy" {
  name                   = "MyReplicationAssurancePolicy"
  description            = "My replication assurance policy"
  evaluation_order_index = 3
  timeout                = "3 s"
}
