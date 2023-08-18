resource "pingdirectory_client_connection_policy" "myClientConnectionPolicy" {
  policy_id              = "default"
  enabled                = false
  evaluation_order_index = 1
}
