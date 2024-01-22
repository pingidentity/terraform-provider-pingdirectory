package debugtarget_test

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

const testIdDebugTarget = "com.example.MyClass"
const testLogPublisherName = "File-Based Debug Logger"

// Attributes to test with. Add optional properties to test here if desired.
type debugTargetTestModel struct {
	logPublisherName string
	debugScope       string
	debugLevel       string
}

func TestAccDebugTarget(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := debugTargetTestModel{
		logPublisherName: testLogPublisherName,
		debugScope:       testIdDebugTarget,
		debugLevel:       "all",
	}
	updatedResourceModel := debugTargetTestModel{
		logPublisherName: testLogPublisherName,
		debugScope:       testIdDebugTarget,
		debugLevel:       "info",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckDebugTargetDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDebugTargetResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedDebugTargetAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_debug_target.%s", resourceName), "debug_scope", initialResourceModel.debugScope),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_debug_target.%s", resourceName), "debug_level", initialResourceModel.debugLevel),
					resource.TestCheckResourceAttrSet("data.pingdirectory_debug_targets.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccDebugTargetResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedDebugTargetAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccDebugTargetResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_debug_target." + resourceName,
				ImportStateId:     updatedResourceModel.logPublisherName + "/" + updatedResourceModel.debugScope,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.DebugTargetAPI.DeleteDebugTarget(ctx, updatedResourceModel.debugScope, updatedResourceModel.logPublisherName).Execute()
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

func testAccDebugTargetResource(resourceName string, resourceModel debugTargetTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_debug_target" "%[1]s" {
  log_publisher_name = "%[2]s"
  debug_scope        = "%[3]s"
  debug_level        = "%[4]s"
}

data "pingdirectory_debug_target" "%[1]s" {
  log_publisher_name = "%[2]s"
  debug_scope        = "%[3]s"
  depends_on = [
    pingdirectory_debug_target.%[1]s
  ]
}

data "pingdirectory_debug_targets" "list" {
  log_publisher_name = "%[2]s"
  depends_on = [
    pingdirectory_debug_target.%[1]s
  ]
}`, resourceName,
		resourceModel.logPublisherName,
		resourceModel.debugScope,
		resourceModel.debugLevel)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDebugTargetAttributes(config debugTargetTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.DebugTargetAPI.GetDebugTarget(ctx, config.debugScope, config.logPublisherName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Debug Target"
		err = acctest.TestAttributesMatchString(resourceType, &config.debugScope, "debug-scope",
			config.debugScope, response.DebugScope)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.debugScope, "debug-level",
			config.debugLevel, response.DebugLevel.String())
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckDebugTargetDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.DebugTargetAPI.GetDebugTarget(ctx, testIdDebugTarget, testLogPublisherName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Debug Target", testIdDebugTarget)
	}
	return nil
}
