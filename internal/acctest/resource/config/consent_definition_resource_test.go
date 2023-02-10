package config_test

import (
	"fmt"
	"testing"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testIdConsentDefinition = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type consentDefinitionTestModel struct {
	id          string
	unique_id   string
	displayName string
}

func TestAccConsentDefinition(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := consentDefinitionTestModel{
		id:          testIdConsentDefinition,
		unique_id:   testIdConsentDefinition,
		displayName: "DisplayName",
	}
	updatedResourceModel := consentDefinitionTestModel{
		id:          testIdConsentDefinition,
		unique_id:   testIdConsentDefinition,
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
				ImportStateId:           updatedResourceModel.id,
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
	 id = "%[2]s"
	 unique_id = "%[3]s"
	 display_name = "%[4]s"
}`, resourceName, resourceModel.id, resourceModel.unique_id, resourceModel.displayName)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedConsentDefinitionAttributes(config consentDefinitionTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ConsentDefinitionApi.GetConsentDefinition(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Consent Definition"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "unique-id",
			testIdConsentDefinition, response.UniqueID)
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
