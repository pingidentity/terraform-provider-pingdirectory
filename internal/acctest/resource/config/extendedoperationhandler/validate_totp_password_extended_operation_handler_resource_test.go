package extendedoperationhandler_test

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

const testIdValidateTotpPasswordExtendedOperationHandler = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type validateTotpPasswordExtendedOperationHandlerTestModel struct {
	id      string
	enabled bool
}

func TestAccValidateTotpPasswordExtendedOperationHandler(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := validateTotpPasswordExtendedOperationHandlerTestModel{
		id:      testIdValidateTotpPasswordExtendedOperationHandler,
		enabled: true,
	}
	updatedResourceModel := validateTotpPasswordExtendedOperationHandlerTestModel{
		id:      testIdValidateTotpPasswordExtendedOperationHandler,
		enabled: false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckValidateTotpPasswordExtendedOperationHandlerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccValidateTotpPasswordExtendedOperationHandlerResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedValidateTotpPasswordExtendedOperationHandlerAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccValidateTotpPasswordExtendedOperationHandlerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedValidateTotpPasswordExtendedOperationHandlerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccValidateTotpPasswordExtendedOperationHandlerResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_extended_operation_handler." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccValidateTotpPasswordExtendedOperationHandlerResource(resourceName string, resourceModel validateTotpPasswordExtendedOperationHandlerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_extended_operation_handler" "%[1]s" {
  type    = "validate-totp-password"
  id      = "%[2]s"
  enabled = %[3]t
}`, resourceName,
		resourceModel.id,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedValidateTotpPasswordExtendedOperationHandlerAttributes(config validateTotpPasswordExtendedOperationHandlerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Validate Totp Password Extended Operation Handler"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.ValidateTotpPasswordExtendedOperationHandlerResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckValidateTotpPasswordExtendedOperationHandlerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(ctx, testIdValidateTotpPasswordExtendedOperationHandler).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Validate Totp Password Extended Operation Handler", testIdValidateTotpPasswordExtendedOperationHandler)
	}
	return nil
}
