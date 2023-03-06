resource "pingdirectory_local_db_index" "pfConnectedIdentityIndex" {
  backend_name = "userRoot"
  attribute    = "pf-connected-identity"
  index_type   = ["equality"]
}

resource "pingdirectory_local_db_index" "pfOauthClientIdIndex" {
  backend_name = "userRoot"
  attribute    = "pf-oauth-client-id"
  index_type   = ["equality", "ordering", "substring"]
}

resource "pingdirectory_local_db_index" "pfOauthClientNameIndex" {
  backend_name = "userRoot"
  attribute    = "pf-oauth-client-name"
  index_type   = ["equality", "ordering", "substring"]
}

resource "pingdirectory_local_db_index" "pfOauthClientLastModifiedIndex" {
  backend_name = "userRoot"
  attribute    = "pf-oauth-client-last-modified"
  index_type   = ["ordering"]
}

resource "pingdirectory_local_db_index" "accessGrantGuidIndex" {
  backend_name = "userRoot"
  attribute    = "accessGrantGuid"
  index_type   = ["equality"]
}

resource "pingdirectory_local_db_index" "accessGrantUniqueUserIndentifierIndex" {
  backend_name = "userRoot"
  attribute    = "accessGrantUniqueUserIdentifier"
  index_type   = ["equality"]
}

resource "pingdirectory_local_db_index" "accessGrantHashedRefreshTokenValueIndex" {
  backend_name = "userRoot"
  attribute    = "accessGrantHashedRefreshTokenValue"
  index_type   = ["equality"]
}

resource "pingdirectory_local_db_index" "accessGrantClientIdIndex" {
  backend_name = "userRoot"
  attribute    = "accessGrantClientId"
  index_type   = ["equality"]
}

resource "pingdirectory_local_db_index" "accessGrantExpiresIndex" {
  backend_name = "userRoot"
  attribute    = "accessGrantExpires"
  index_type   = ["ordering"]
}
