package pluginroot_test

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

// Attributes to test with. Add optional properties to test here if desired.
type pluginRootTestModel struct {
	pluginOrderPreParseSearch string
}

func TestAccPluginRoot(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := pluginRootTestModel{
		pluginOrderPreParseSearch: "7-Bit Clean",
	}
	updatedResourceModel := pluginRootTestModel{
		pluginOrderPreParseSearch: "Entry UUID",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccPluginRootResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedPluginRootAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_plugin_root.%s", resourceName), "plugin_order_pre_parse_search", initialResourceModel.pluginOrderPreParseSearch),
				),
			},
			{
				// Test updating some fields
				Config: testAccPluginRootResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedPluginRootAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccPluginRootResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_default_plugin_root." + resourceName,
				ImportStateId:     initialResourceModel.pluginOrderPreParseSearch,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPluginRootResource(resourceName string, resourceModel pluginRootTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_plugin_root" "%[1]s" {
  plugin_order_pre_parse_search = "%[2]s"
}

data "pingdirectory_plugin_root" "%[1]s" {
  depends_on = [
    pingdirectory_default_plugin_root.%[1]s
  ]
}`, resourceName,
		resourceModel.pluginOrderPreParseSearch)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedPluginRootAttributes(config pluginRootTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "Plugin Root"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.PluginRootAPI.GetPluginRoot(ctx).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, nil, "plugin-order-pre-parse-search",
			config.pluginOrderPreParseSearch, *response.PluginOrderPreParseSearch)
		if err != nil {
			return err
		}
		return nil
	}
}
