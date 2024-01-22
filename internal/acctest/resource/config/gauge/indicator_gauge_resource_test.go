package gauge_test

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

const testIdIndicatorGauge = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type indicatorGaugeTestModel struct {
	id              string
	gaugeDataSource string
	enabled         bool
}

func TestAccIndicatorGauge(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := indicatorGaugeTestModel{
		id:              testIdIndicatorGauge,
		gaugeDataSource: "Strong Encryption Not Available",
		enabled:         true,
	}
	updatedResourceModel := indicatorGaugeTestModel{
		id:              testIdIndicatorGauge,
		gaugeDataSource: "Replication Connection Status",
		enabled:         false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckIndicatorGaugeDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccIndicatorGaugeResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedIndicatorGaugeAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccIndicatorGaugeResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedIndicatorGaugeAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccIndicatorGaugeResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_gauge." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIndicatorGaugeResource(resourceName string, resourceModel indicatorGaugeTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_gauge" "%[1]s" {
  type              = "indicator"
  name              = "%[2]s"
  gauge_data_source = "%[3]s"
  enabled           = %[4]t
}`, resourceName, resourceModel.id,
		resourceModel.gaugeDataSource,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedIndicatorGaugeAttributes(config indicatorGaugeTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.GaugeAPI.GetGauge(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Indicator Gauge"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "gauge-data-source",
			config.gaugeDataSource, response.IndicatorGaugeResponse.GaugeDataSource)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.IndicatorGaugeResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckIndicatorGaugeDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.GaugeAPI.GetGauge(ctx, testIdIndicatorGauge).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Indicator Gauge", testIdIndicatorGauge)
	}
	return nil
}
