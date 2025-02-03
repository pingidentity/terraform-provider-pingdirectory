// Copyright © 2025 Ping Identity Corporation

package consentdefinitionlocalization_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdConsentDefinitionLocalization = "en-US"
const testConsentDefinitionName = "myConsentDefinition"

// Attributes to test with. Add optional properties to test here if desired.
type consentDefinitionLocalizationTestModel struct {
	consentDefinitionName string
	locale                string
	version               string
	dataText              string
	purposeText           string
}

func TestAccConsentDefinitionLocalization(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := consentDefinitionLocalizationTestModel{
		consentDefinitionName: testConsentDefinitionName,
		locale:                testIdConsentDefinitionLocalization,
		version:               "1.1",
		dataText:              "example data text",
		purposeText:           "example purpose text",
	}
	updatedResourceModel := consentDefinitionLocalizationTestModel{
		consentDefinitionName: testConsentDefinitionName,
		locale:                testIdConsentDefinitionLocalization,
		version:               "1.2",
		dataText:              "example updated data text",
		purposeText:           "example updated purpose text",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckConsentDefinitionLocalizationDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccConsentDefinitionLocalizationResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedConsentDefinitionLocalizationAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_consent_definition_localization.%s", resourceName), "locale", initialResourceModel.locale),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_consent_definition_localization.%s", resourceName), "version", initialResourceModel.version),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_consent_definition_localization.%s", resourceName), "data_text", initialResourceModel.dataText),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_consent_definition_localization.%s", resourceName), "purpose_text", initialResourceModel.purposeText),
					resource.TestCheckResourceAttrSet("data.pingdirectory_consent_definition_localizations.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccConsentDefinitionLocalizationResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedConsentDefinitionLocalizationAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccConsentDefinitionLocalizationResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_consent_definition_localization." + resourceName,
				ImportStateId:     updatedResourceModel.consentDefinitionName + "/" + updatedResourceModel.locale,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.ConsentDefinitionLocalizationAPI.DeleteConsentDefinitionLocalization(ctx, updatedResourceModel.locale, updatedResourceModel.consentDefinitionName).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccConsentDefinitionLocalizationResource(resourceName string, resourceModel consentDefinitionLocalizationTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_consent_definition" "%[2]s" {
  unique_id    = "%[2]s"
  display_name = "example display name"
}
resource "pingdirectory_consent_definition_localization" "%[1]s" {
  consent_definition_name = pingdirectory_consent_definition.%[2]s.unique_id
  locale                  = "%[3]s"
  version                 = "%[4]s"
  data_text               = "%[5]s"
  purpose_text            = "%[6]s"
}

data "pingdirectory_consent_definition_localization" "%[1]s" {
  consent_definition_name = "%[2]s"
  locale                  = "%[3]s"
  depends_on = [
    pingdirectory_consent_definition_localization.%[1]s
  ]
}

data "pingdirectory_consent_definition_localizations" "list" {
  consent_definition_name = "%[2]s"
  depends_on = [
    pingdirectory_consent_definition_localization.%[1]s
  ]
}`, resourceName,
		resourceModel.consentDefinitionName,
		resourceModel.locale,
		resourceModel.version,
		resourceModel.dataText,
		resourceModel.purposeText)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedConsentDefinitionLocalizationAttributes(config consentDefinitionLocalizationTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ConsentDefinitionLocalizationAPI.GetConsentDefinitionLocalization(ctx, config.locale, config.consentDefinitionName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Consent Definition Localization"
		err = acctest.TestAttributesMatchString(resourceType, &config.locale, "locale",
			config.locale, response.Locale)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.locale, "version",
			config.version, response.Version)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.locale, "data-text",
			config.dataText, response.DataText)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.locale, "purpose-text",
			config.purposeText, response.PurposeText)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckConsentDefinitionLocalizationDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ConsentDefinitionLocalizationAPI.GetConsentDefinitionLocalization(ctx, testIdConsentDefinitionLocalization, testConsentDefinitionName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Consent Definition Localization", testIdConsentDefinitionLocalization)
	}
	return nil
}
