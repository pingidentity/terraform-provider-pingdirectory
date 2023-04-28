package gaugedatasource_test

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

const testIdIndicatorGaugeDataSource = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type indicatorGaugeDataSourceTestModel struct {
	id                 string
	monitorObjectclass string
	monitorAttribute   string
}

func TestAccIndicatorGaugeDataSource(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := indicatorGaugeDataSourceTestModel{
		id:                 testIdIndicatorGaugeDataSource,
		monitorObjectclass: "ds-host-system-disk-monitor-entry",
		monitorAttribute:   "pct-busy",
	}
	updatedResourceModel := indicatorGaugeDataSourceTestModel{
		id:                 testIdIndicatorGaugeDataSource,
		monitorObjectclass: "ds-host-system-cpu-memory-monitor-entry",
		monitorAttribute:   "recent-cpu-used",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckIndicatorGaugeDataSourceDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccIndicatorGaugeDataSourceResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedIndicatorGaugeDataSourceAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccIndicatorGaugeDataSourceResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedIndicatorGaugeDataSourceAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccIndicatorGaugeDataSourceResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_indicator_gauge_data_source." + resourceName,
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

func testAccIndicatorGaugeDataSourceResource(resourceName string, resourceModel indicatorGaugeDataSourceTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_indicator_gauge_data_source" "%[1]s" {
  id                  = "%[2]s"
  monitor_objectclass = "%[3]s"
  monitor_attribute   = "%[4]s"
}`, resourceName,
		resourceModel.id,
		resourceModel.monitorObjectclass,
		resourceModel.monitorAttribute)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedIndicatorGaugeDataSourceAttributes(config indicatorGaugeDataSourceTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.GaugeDataSourceApi.GetGaugeDataSource(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Indicator Gauge Data Source"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "monitor-objectclass",
			config.monitorObjectclass, response.IndicatorGaugeDataSourceResponse.MonitorObjectclass)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "monitor-attribute",
			config.monitorAttribute, response.IndicatorGaugeDataSourceResponse.MonitorAttribute)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckIndicatorGaugeDataSourceDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.GaugeDataSourceApi.GetGaugeDataSource(ctx, testIdIndicatorGaugeDataSource).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Indicator Gauge Data Source", testIdIndicatorGaugeDataSource)
	}
	return nil
}
