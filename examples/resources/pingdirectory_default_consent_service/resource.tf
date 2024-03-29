resource "pingdirectory_default_consent_service" "myConsentService" {
  enabled                    = true
  base_dn                    = "ou=consents,dc=example,dc=com"
  bind_dn                    = "cn=consent service account"
  unprivileged_consent_scope = "urn:pingdirectory:consent"
  privileged_consent_scope   = "urn:pingdirectory:consent_admin"
  search_size_limit          = 90
}
