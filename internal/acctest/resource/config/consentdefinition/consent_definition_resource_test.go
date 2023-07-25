package consentdefinition_test

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

const testIdConsentDefinition = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type consentDefinitionTestModel struct {
	uniqueId    string
	displayName string
}

func TestAccConsentDefinition(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := consentDefinitionTestModel{
		uniqueId:    testIdConsentDefinition,
		displayName: "DisplayName",
	}
	updatedResourceModel := consentDefinitionTestModel{
		uniqueId:    testIdConsentDefinition,
		displayName: "DisplayName1",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckConsentDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccConsentDefinitionResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedConsentDefinitionAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccConsentDefinitionResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedConsentDefinitionAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccConsentDefinitionResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_consent_definition." + resourceName,
				ImportStateId:           updatedResourceModel.uniqueId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccConsentDefinitionResource(resourceName string, resourceModel consentDefinitionTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_consent_definition" "%[1]s" {
  unique_id    = "%[2]s"
  display_name = "%[3]s"
}`, resourceName,
		resourceModel.uniqueId,
		resourceModel.displayName)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedConsentDefinitionAttributes(config consentDefinitionTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ConsentDefinitionApi.GetConsentDefinition(ctx, config.uniqueId).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Consent Definition"
		err = acctest.TestAttributesMatchString(resourceType, &config.uniqueId, "unique-id",
			config.uniqueId, response.UniqueID)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.uniqueId, "display-name",
			config.displayName, response.DisplayName)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckConsentDefinitionDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ConsentDefinitionApi.GetConsentDefinition(ctx, testIdConsentDefinition).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Consent Definition", testIdConsentDefinition)
	}
	return nil
}
