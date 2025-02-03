// Copyright © 2025 Ping Identity Corporation

package plugin_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdInternalSearchRatePlugin = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type internalSearchRatePluginTestModel struct {
	id           string
	description  string
	pluginType   []string
	numThreads   int64
	baseDn       string
	filterPrefix string
	enabled      bool
}

func TestAccInternalSearchRatePlugin(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := internalSearchRatePluginTestModel{
		id:           testIdInternalSearchRatePlugin,
		description:  "Test simple internal search rate plugin",
		pluginType:   []string{"shutdown", "startup"},
		numThreads:   10,
		baseDn:       "dc=example1,dc=com",
		filterPrefix: "myprefix",
		enabled:      true,
	}
	updatedResourceModel := internalSearchRatePluginTestModel{
		id:           testIdInternalSearchRatePlugin,
		description:  "Test simple internal search rate plugin modified",
		pluginType:   []string{"startup"},
		numThreads:   9,
		baseDn:       "dc=example2,dc=com",
		filterPrefix: "myprefix2",
		enabled:      true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckInternalSearchRatePluginDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccInternalSearchRatePluginResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedInternalSearchRatePluginAttributes(initialResourceModel),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_plugin.%s", resourceName), "base_dn.*", initialResourceModel.baseDn),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_plugin.%s", resourceName), "filter_prefix", initialResourceModel.filterPrefix),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_plugin.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_plugins.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccInternalSearchRatePluginResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedInternalSearchRatePluginAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccInternalSearchRatePluginResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_plugin." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.PluginAPI.DeletePlugin(ctx, updatedResourceModel.id).Execute()
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

func testAccInternalSearchRatePluginResource(resourceName string, resourceModel internalSearchRatePluginTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_plugin" "%[1]s" {
  resource_type = "internal-search-rate"
  name          = "%[2]s"
  description   = "%[3]s"
  plugin_type   = %[4]s
  num_threads   = %[5]d
  base_dn       = ["%[6]s"]
  filter_prefix = "%[7]s"
  enabled       = %[8]t
}

data "pingdirectory_plugin" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_plugin.%[1]s
  ]
}

data "pingdirectory_plugins" "list" {
  depends_on = [
    pingdirectory_plugin.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.description,
		acctest.StringSliceToTerraformString(resourceModel.pluginType),
		resourceModel.numThreads,
		resourceModel.baseDn,
		resourceModel.filterPrefix,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedInternalSearchRatePluginAttributes(config internalSearchRatePluginTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.PluginAPI.GetPlugin(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Internal Search Rate Plugin"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "description",
			config.description, *response.InternalSearchRatePluginResponse.Description)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "plugin-type",
			config.pluginType, client.StringSliceEnumpluginPluginTypeProp((response.InternalSearchRatePluginResponse.PluginType)))
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "num-threads",
			config.numThreads, response.InternalSearchRatePluginResponse.NumThreads)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "base-dn",
			config.baseDn, response.InternalSearchRatePluginResponse.BaseDN)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "filter-prefix",
			config.filterPrefix, response.InternalSearchRatePluginResponse.FilterPrefix)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.InternalSearchRatePluginResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckInternalSearchRatePluginDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.PluginAPI.GetPlugin(ctx, testIdInternalSearchRatePlugin).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Internal Search Rate Plugin", testIdInternalSearchRatePlugin)
	}
	return nil
}
