package monitorprovider_test

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

const testResource = "General Monitor Entry"

// Attributes to test with. Add optional properties to test here if desired.
type generalMonitorProviderTestModel struct {
	id          string
	description string
	enabled     bool
}

func TestAccGeneralPartyMonitorProvider(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := generalMonitorProviderTestModel{
		id:          testResource,
		description: "Initial general monitor provider",
		enabled:     false,
	}
	// default is enabled=true, set this at end of test
	updatedResourceModel := generalMonitorProviderTestModel{
		id:          testResource,
		description: "Updated general monitor provider",
		enabled:     true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccGeneralMonitorProviderResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedGeneralMonitorProviderAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccGeneralMonitorProviderResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedGeneralMonitorProviderAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccGeneralMonitorProviderResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_default_general_monitor_provider." + resourceName,
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

func testAccGeneralMonitorProviderResource(resourceName string, resourceModel generalMonitorProviderTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_general_monitor_provider" "%[1]s" {
  id          = "%[2]s"
  description = "%[3]s"
  enabled     = %[4]t
}`, resourceName,
		resourceModel.id,
		resourceModel.description,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedGeneralMonitorProviderAttributes(config generalMonitorProviderTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.MonitorProviderApi.GetMonitorProvider(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "general"
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "description",
			config.description, response.GeneralMonitorProviderResponse.Description)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.GeneralMonitorProviderResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}
