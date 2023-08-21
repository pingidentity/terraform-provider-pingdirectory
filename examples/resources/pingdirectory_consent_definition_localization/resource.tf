resource "pingdirectory_consent_definition" "myConsentDefinition" {
  unique_id    = "myConsentDefinition"
  display_name = "example display name"
}

resource "pingdirectory_consent_definition_localization" "myConsentDefinitionLocalization" {
  consent_definition_name = pingdirectory_consent_definition.myConsentDefinition.unique_id
  locale                  = "en-US"
  version                 = "1.1"
  data_text               = "example data text"
  purpose_text            = "example purpose text"
}
