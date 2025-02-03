// Copyright © 2025 Ping Identity Corporation

package extendedoperationhandler_test

import (
	"fmt"
	"strconv"
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
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckValidateTotpPasswordExtendedOperationHandlerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccValidateTotpPasswordExtendedOperationHandlerResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedValidateTotpPasswordExtendedOperationHandlerAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_extended_operation_handler.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_extended_operation_handlers.list", "objects.0.id"),
				),
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
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.ExtendedOperationHandlerAPI.DeleteExtendedOperationHandler(ctx, updatedResourceModel.id).Execute()
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

func testAccValidateTotpPasswordExtendedOperationHandlerResource(resourceName string, resourceModel validateTotpPasswordExtendedOperationHandlerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_extended_operation_handler" "%[1]s" {
  type    = "validate-totp-password"
  name    = "%[2]s"
  enabled = %[3]t
}

data "pingdirectory_extended_operation_handler" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_extended_operation_handler.%[1]s
  ]
}

data "pingdirectory_extended_operation_handlers" "list" {
  depends_on = [
    pingdirectory_extended_operation_handler.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedValidateTotpPasswordExtendedOperationHandlerAttributes(config validateTotpPasswordExtendedOperationHandlerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ExtendedOperationHandlerAPI.GetExtendedOperationHandler(ctx, config.id).Execute()
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
	_, _, err := testClient.ExtendedOperationHandlerAPI.GetExtendedOperationHandler(ctx, testIdValidateTotpPasswordExtendedOperationHandler).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Validate Totp Password Extended Operation Handler", testIdValidateTotpPasswordExtendedOperationHandler)
	}
	return nil
}
