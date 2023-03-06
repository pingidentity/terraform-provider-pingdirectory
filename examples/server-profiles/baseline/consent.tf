resource "pingdirectory_consent_definition" "emailConsentDefinition" {
  unique_id    = "email"
  display_name = "Email Address"
  description  = "Share your email address"
}

resource "pingdirectory_consent_definition_localization" "emailConsentDefinitionLocalization" {
  locale                  = "en-US"
  consent_definition_name = pingdirectory_consent_definition.emailConsentDefinition.unique_id
  version                 = "1.0"
  title_text              = "Share your email address"
  data_text               = "Your email address"
  purpose_text            = "Join Mailing List"
}

resource "pingdirectory_directory_rest_api_http_servlet_extension" "defaultDirectoryRestApiExtension" {
  id                 = "Directory REST API"
  access_token_scope = "ds"
}

resource "pingdirectory_exact_match_identity_mapper" "userIdIdentityMapper" {
  id              = "user-id-identity-mapper"
  enabled         = true
  match_attribute = ["cn", "entryUUID", "uid"]
  match_base_dn   = ["cn=config", "ou=people,${var.user_base_dn}"]
}

resource "pingdirectory_mock_access_token_validator" "mockAccessTokenValidate" {
  id                     = "mock-access-token-validator"
  identity_mapper        = pingdirectory_exact_match_identity_mapper.userIdIdentityMapper.id
  enabled                = true
  evaluation_order_index = 1
}

resource "pingdirectory_topology_admin_user" "consentInternalServiceAccount" {
  id                              = "Consent API internal service account"
  alternate_bind_dn               = ["cn=consent service account"]
  first_name                      = ["Consent"]
  inherit_default_root_privileges = false
  last_name                       = ["Internal Service Account"]
  password                        = "rootpassword"
  privilege                       = ["bypass-acl", "config-read"]
}

resource "pingdirectory_consent_service" "defaultConsentService" {
  enabled                        = true
  base_dn                        = "ou=Consents,${var.user_base_dn}"
  bind_dn                        = "cn=consent service account"
  consent_record_identity_mapper = [pingdirectory_exact_match_identity_mapper.userIdIdentityMapper.id]
  # The above attribute must be changed to allow destroying the user-id-identity-mapper object.
  # consent_record_identity_mapper = []
  service_account_dn             = ["uid=Consent Admin,ou=people,${var.user_base_dn}"]
  unprivileged_consent_scope     = "consent"
  privileged_consent_scope       = "consent_admin"
}

resource "pingdirectory_consent_http_servlet_extension" "defaultConsentServletExtension" {
  id              = "Consent"
  identity_mapper = pingdirectory_exact_match_identity_mapper.userIdIdentityMapper.id
  # The above attribute must be changed to allow destroying the user-id-identity-mapper object.
  # Exact Match is the default identity mapper.
  # identity_mapper = "Exact Match"
  depends_on = [
    pingdirectory_consent_service.defaultConsentService
  ]
}
