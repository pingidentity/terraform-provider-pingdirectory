package synchronizationprovider_test

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

const testIdCustomSynchronizationProvider = "Changelog Ordering"

// Attributes to test with. Add optional properties to test here if desired.
type customSynchronizationProviderTestModel struct {
	id      string
	enabled bool
}

func TestAccCustomSynchronizationProvider(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := customSynchronizationProviderTestModel{
		id:      testIdCustomSynchronizationProvider,
		enabled: true,
	}
	updatedResourceModel := customSynchronizationProviderTestModel{
		id:      testIdCustomSynchronizationProvider,
		enabled: false,
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
				Config: testAccCustomSynchronizationProviderResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedCustomSynchronizationProviderAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_synchronization_provider.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
				),
			},
			{
				// Test updating some fields
				Config: testAccCustomSynchronizationProviderResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedCustomSynchronizationProviderAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccCustomSynchronizationProviderResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_default_synchronization_provider." + resourceName,
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

func testAccCustomSynchronizationProviderResource(resourceName string, resourceModel customSynchronizationProviderTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_synchronization_provider" "%[1]s" {
  type    = "custom"
  id      = "%[2]s"
  enabled = %[3]t
}

data "pingdirectory_synchronization_provider" "%[1]s" {
  id = "%[2]s"
  depends_on = [
    pingdirectory_default_synchronization_provider.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedCustomSynchronizationProviderAttributes(config customSynchronizationProviderTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.SynchronizationProviderApi.GetSynchronizationProvider(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Custom Synchronization Provider"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.CustomSynchronizationProviderResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}
