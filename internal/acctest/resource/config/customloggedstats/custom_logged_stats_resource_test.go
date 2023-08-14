package customloggedstats_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdCustomLoggedStats = "MyId"
const testPluginName = "JSON Stats Logger"

// Attributes to test with. Add optional properties to test here if desired.
type customLoggedStatsTestModel struct {
	id                 string
	pluginName         string
	monitorObjectclass string
	attributeToLog     []string
	statisticType      []string
}

func TestAccCustomLoggedStats(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := customLoggedStatsTestModel{
		id:                 testIdCustomLoggedStats,
		pluginName:         testPluginName,
		monitorObjectclass: "ds-memory-usage-monitor-entry",
		attributeToLog:     []string{"total-bytes-used-by-memory-consumers"},
		statisticType:      []string{"raw"},
	}
	updatedResourceModel := customLoggedStatsTestModel{
		id:                 testIdCustomLoggedStats,
		pluginName:         testPluginName,
		monitorObjectclass: "ds-memory-usage-monitor-entry",
		attributeToLog:     []string{"total-bytes-used-by-memory-consumers", "non-heap-memory-bytes-used"},
		statisticType:      []string{"average"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckCustomLoggedStatsDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccCustomLoggedStatsResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedCustomLoggedStatsAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_custom_logged_stats.%s", resourceName), "monitor_objectclass", initialResourceModel.monitorObjectclass),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_custom_logged_stats.%s", resourceName), "attribute_to_log.*", initialResourceModel.attributeToLog[0]),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_custom_logged_stats.%s", resourceName), "statistic_type.*", initialResourceModel.statisticType[0]),
					resource.TestCheckResourceAttrSet("data.pingdirectory_custom_logged_stats_list.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccCustomLoggedStatsResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedCustomLoggedStatsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccCustomLoggedStatsResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_custom_logged_stats." + resourceName,
				ImportStateId:     updatedResourceModel.pluginName + "/" + updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.CustomLoggedStatsApi.DeleteCustomLoggedStats(ctx, updatedResourceModel.id, updatedResourceModel.pluginName).Execute()
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

func testAccCustomLoggedStatsResource(resourceName string, resourceModel customLoggedStatsTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_custom_logged_stats" "%[1]s" {
  name                = "%[2]s"
  plugin_name         = "%[3]s"
  monitor_objectclass = "%[4]s"
  attribute_to_log    = %[5]s
  statistic_type      = %[6]s
}

data "pingdirectory_custom_logged_stats" "%[1]s" {
  name        = "%[2]s"
  plugin_name = "%[3]s"
  depends_on = [
    pingdirectory_custom_logged_stats.%[1]s
  ]
}

data "pingdirectory_custom_logged_stats_list" "list" {
  plugin_name = "%[3]s"
  depends_on = [
    pingdirectory_custom_logged_stats.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.pluginName,
		resourceModel.monitorObjectclass,
		acctest.StringSliceToTerraformString(resourceModel.attributeToLog),
		acctest.StringSliceToTerraformString(resourceModel.statisticType))
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedCustomLoggedStatsAttributes(config customLoggedStatsTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.CustomLoggedStatsApi.GetCustomLoggedStats(ctx, config.id, config.pluginName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Custom Logged Stats"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "monitor-objectclass",
			config.monitorObjectclass, response.MonitorObjectclass)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "attribute-to-log",
			config.attributeToLog, response.AttributeToLog)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "statistic-type",
			config.statisticType, client.StringSliceEnumcustomLoggedStatsStatisticTypeProp(response.StatisticType))
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckCustomLoggedStatsDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.CustomLoggedStatsApi.GetCustomLoggedStats(ctx, testIdCustomLoggedStats, testPluginName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Custom Logged Stats", testIdCustomLoggedStats)
	}
	return nil
}
