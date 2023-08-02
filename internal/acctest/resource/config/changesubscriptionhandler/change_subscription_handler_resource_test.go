package changesubscriptionhandler_test

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

const testIdLoggingChangeSubscriptionHandler = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type loggingChangeSubscriptionHandlerTestModel struct {
	id      string
	enabled bool
}

func TestAccLoggingChangeSubscriptionHandler(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := loggingChangeSubscriptionHandlerTestModel{
		id:      testIdLoggingChangeSubscriptionHandler,
		enabled: true,
	}
	updatedResourceModel := loggingChangeSubscriptionHandlerTestModel{
		id:      testIdLoggingChangeSubscriptionHandler,
		enabled: false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckLoggingChangeSubscriptionHandlerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccLoggingChangeSubscriptionHandlerResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedLoggingChangeSubscriptionHandlerAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_change_subscription_handler.%s", resourceName), "id", initialResourceModel.id),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_change_subscription_handler.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_change_subscription_handlers.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccLoggingChangeSubscriptionHandlerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedLoggingChangeSubscriptionHandlerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccLoggingChangeSubscriptionHandlerResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_change_subscription_handler." + resourceName,
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

func testAccLoggingChangeSubscriptionHandlerResource(resourceName string, resourceModel loggingChangeSubscriptionHandlerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_change_subscription_handler" "%[1]s" {
  type    = "logging"
  name    = "%[2]s"
  enabled = %[3]t
}

data "pingdirectory_change_subscription_handler" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_change_subscription_handler.%[1]s
  ]
}

data "pingdirectory_change_subscription_handlers" "list" {
  depends_on = [
    pingdirectory_change_subscription_handler.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedLoggingChangeSubscriptionHandlerAttributes(config loggingChangeSubscriptionHandlerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ChangeSubscriptionHandlerApi.GetChangeSubscriptionHandler(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Logging Change Subscription Handler"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.LoggingChangeSubscriptionHandlerResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckLoggingChangeSubscriptionHandlerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ChangeSubscriptionHandlerApi.GetChangeSubscriptionHandler(ctx, testIdLoggingChangeSubscriptionHandler).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Logging Change Subscription Handler", testIdLoggingChangeSubscriptionHandler)
	}
	return nil
}
