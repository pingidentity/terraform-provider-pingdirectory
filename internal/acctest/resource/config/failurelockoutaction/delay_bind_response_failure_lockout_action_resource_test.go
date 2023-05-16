package failurelockoutaction_test

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

const testIdDelayBindResponseFailureLockoutAction = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type delayBindResponseFailureLockoutActionTestModel struct {
	id    string
	delay string
}

func TestAccDelayBindResponseFailureLockoutAction(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := delayBindResponseFailureLockoutActionTestModel{
		id:    testIdDelayBindResponseFailureLockoutAction,
		delay: "1 s",
	}
	updatedResourceModel := delayBindResponseFailureLockoutActionTestModel{
		id:    testIdDelayBindResponseFailureLockoutAction,
		delay: "10 s",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckDelayBindResponseFailureLockoutActionDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDelayBindResponseFailureLockoutActionResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedDelayBindResponseFailureLockoutActionAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccDelayBindResponseFailureLockoutActionResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedDelayBindResponseFailureLockoutActionAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccDelayBindResponseFailureLockoutActionResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_delay_bind_response_failure_lockout_action." + resourceName,
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

func testAccDelayBindResponseFailureLockoutActionResource(resourceName string, resourceModel delayBindResponseFailureLockoutActionTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_delay_bind_response_failure_lockout_action" "%[1]s" {
  id    = "%[2]s"
  delay = "%[3]s"
}`, resourceName,
		resourceModel.id,
		resourceModel.delay)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDelayBindResponseFailureLockoutActionAttributes(config delayBindResponseFailureLockoutActionTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.FailureLockoutActionApi.GetFailureLockoutAction(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Delay Bind Response Failure Lockout Action"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "delay",
			config.delay, response.DelayBindResponseFailureLockoutActionResponse.Delay)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckDelayBindResponseFailureLockoutActionDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.FailureLockoutActionApi.GetFailureLockoutAction(ctx, testIdDelayBindResponseFailureLockoutAction).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Delay Bind Response Failure Lockout Action", testIdDelayBindResponseFailureLockoutAction)
	}
	return nil
}