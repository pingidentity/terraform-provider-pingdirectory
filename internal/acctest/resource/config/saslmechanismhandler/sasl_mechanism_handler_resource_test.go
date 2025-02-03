// Copyright © 2025 Ping Identity Corporation

package saslmechanismhandler_test

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

const testIdUnboundidMsChapV2SaslMechanismHandler = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type unboundidMsChapV2SaslMechanismHandlerTestModel struct {
	id             string
	identityMapper string
	enabled        bool
}

func TestAccUnboundidMsChapV2SaslMechanismHandler(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := unboundidMsChapV2SaslMechanismHandlerTestModel{
		id:             testIdUnboundidMsChapV2SaslMechanismHandler,
		identityMapper: "Exact Match",
		enabled:        true,
	}
	updatedResourceModel := unboundidMsChapV2SaslMechanismHandlerTestModel{
		id:             testIdUnboundidMsChapV2SaslMechanismHandler,
		identityMapper: "All Admin Users",
		enabled:        false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckUnboundidMsChapV2SaslMechanismHandlerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccUnboundidMsChapV2SaslMechanismHandlerResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedUnboundidMsChapV2SaslMechanismHandlerAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_sasl_mechanism_handler.%s", resourceName), "identity_mapper", initialResourceModel.identityMapper),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_sasl_mechanism_handler.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_sasl_mechanism_handlers.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccUnboundidMsChapV2SaslMechanismHandlerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedUnboundidMsChapV2SaslMechanismHandlerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccUnboundidMsChapV2SaslMechanismHandlerResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_sasl_mechanism_handler." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.SaslMechanismHandlerAPI.DeleteSaslMechanismHandler(ctx, updatedResourceModel.id).Execute()
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

func testAccUnboundidMsChapV2SaslMechanismHandlerResource(resourceName string, resourceModel unboundidMsChapV2SaslMechanismHandlerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_sasl_mechanism_handler" "%[1]s" {
  type            = "unboundid-ms-chap-v2"
  name            = "%[2]s"
  identity_mapper = "%[3]s"
  enabled         = %[4]t
}

data "pingdirectory_sasl_mechanism_handler" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_sasl_mechanism_handler.%[1]s
  ]
}

data "pingdirectory_sasl_mechanism_handlers" "list" {
  depends_on = [
    pingdirectory_sasl_mechanism_handler.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.identityMapper,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedUnboundidMsChapV2SaslMechanismHandlerAttributes(config unboundidMsChapV2SaslMechanismHandlerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.SaslMechanismHandlerAPI.GetSaslMechanismHandler(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Unboundid Ms Chap V2 Sasl Mechanism Handler"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "identity-mapper",
			config.identityMapper, response.UnboundidMsChapV2SaslMechanismHandlerResponse.IdentityMapper)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.UnboundidMsChapV2SaslMechanismHandlerResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckUnboundidMsChapV2SaslMechanismHandlerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.SaslMechanismHandlerAPI.GetSaslMechanismHandler(ctx, testIdUnboundidMsChapV2SaslMechanismHandler).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Unboundid Ms Chap V2 Sasl Mechanism Handler", testIdUnboundidMsChapV2SaslMechanismHandler)
	}
	return nil
}
