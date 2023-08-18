resource "pingdirectory_identity_mapper" "myIdentityMapper" {
  name            = "MyIdentityMapper"
  type            = "exact-match"
  match_attribute = ["uid"]
  enabled         = true
}
