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

const testIdDebugTarget = "MyId"
const testLogPublisherName = "File-Based Debug Logger"

// Attributes to test with. Add optional properties to test here if desired.
type debugTargetTestModel struct {
	id               string
	logPublisherName string
	debugScope       string
	debugLevel       string
}

func TestAccDebugTarget(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := debugTargetTestModel{
		id:               testIdDebugTarget,
		logPublisherName: testLogPublisherName,
		debugScope:       testIdDebugTarget,
		debugLevel:       "all",
	}
	updatedResourceModel := debugTargetTestModel{
		id:               testIdDebugTarget,
		logPublisherName: testLogPublisherName,
		debugScope:       testIdDebugTarget,
		debugLevel:       "info",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckDebugTargetDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDebugTargetResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedDebugTargetAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccDebugTargetResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedDebugTargetAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccDebugTargetResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_debug_target." + resourceName,
				ImportStateId:           updatedResourceModel.logPublisherName + "/" + updatedResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccDebugTargetResource(resourceName string, resourceModel debugTargetTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_debug_target" "%[1]s" {
	 id = "%[2]s"
	 log_publisher_name = "%[3]s"
	 debug_scope = "%[4]s"
	 debug_level = "%[5]s"
}`, resourceName, resourceModel.id,
		resourceModel.logPublisherName,
		resourceModel.debugScope,
		resourceModel.debugLevel)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDebugTargetAttributes(config debugTargetTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.DebugTargetApi.GetDebugTarget(ctx, config.id, config.logPublisherName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Debug Target"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "debug-scope",
			config.debugScope, response.DebugScope)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "debug-level",
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
	_, _, err := testClient.DebugTargetApi.GetDebugTarget(ctx, testIdDebugTarget, testLogPublisherName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Debug Target", testIdDebugTarget)
	}
	return nil
}
