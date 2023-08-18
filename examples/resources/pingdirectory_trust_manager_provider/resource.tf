resource "pingdirectory_trust_manager_provider" "myTrustManagerProvider" {
  name    = "MyTrustManagerProvider"
  type    = "blind"
  enabled = false
}
