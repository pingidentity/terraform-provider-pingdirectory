package logfieldbehavior_test

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

const testIdTextAccessLogFieldBehavior = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type textAccessLogFieldBehaviorTestModel struct {
	id          string
	description string
}

func TestAccTextAccessLogFieldBehavior(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := textAccessLogFieldBehaviorTestModel{
		id:          testIdTextAccessLogFieldBehavior,
		description: "initial description",
	}
	updatedResourceModel := textAccessLogFieldBehaviorTestModel{
		id:          testIdTextAccessLogFieldBehavior,
		description: "updated description",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckTextAccessLogFieldBehaviorDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccTextAccessLogFieldBehaviorResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedTextAccessLogFieldBehaviorAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_log_field_behavior.%s", resourceName), "description", initialResourceModel.description),
					resource.TestCheckResourceAttrSet("data.pingdirectory_log_field_behaviors.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccTextAccessLogFieldBehaviorResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedTextAccessLogFieldBehaviorAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccTextAccessLogFieldBehaviorResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_log_field_behavior." + resourceName,
				ImportStateId:     initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccTextAccessLogFieldBehaviorResource(resourceName string, resourceModel textAccessLogFieldBehaviorTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_log_field_behavior" "%[1]s" {
  type        = "text-access"
  name        = "%[2]s"
  description = "%[3]s"
}

data "pingdirectory_log_field_behavior" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_log_field_behavior.%[1]s
  ]
}

data "pingdirectory_log_field_behaviors" "list" {
  depends_on = [
    pingdirectory_log_field_behavior.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedTextAccessLogFieldBehaviorAttributes(config textAccessLogFieldBehaviorTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LogFieldBehaviorApi.GetLogFieldBehavior(ctx, config.id).Execute()
		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		resourceType := "Text Access Log Field Behavior"
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "description",
			config.description, response.TextAccessLogFieldBehaviorResponse.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckTextAccessLogFieldBehaviorDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.LogFieldBehaviorApi.GetLogFieldBehavior(ctx, testIdTextAccessLogFieldBehavior).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Text Access Log Field Behavior", testIdTextAccessLogFieldBehavior)
	}
	return nil
}
