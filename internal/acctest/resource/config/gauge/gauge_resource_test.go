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

const testIdGauge = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type gaugeTestModel struct {
	id              string
	gaugeDataSource string
	enabled         bool
}

func TestAccGauge(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := gaugeTestModel{
		id:              testIdGauge,
		gaugeDataSource: "Strong Encryption Not Available",
		enabled:         true,
	}
	updatedResourceModel := gaugeTestModel{
		id:              testIdGauge,
		gaugeDataSource: "Replication Connection Status",
		enabled:         false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckGaugeDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccGaugeResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedGaugeAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccGaugeResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedGaugeAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccGaugeResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_gauge." + resourceName,
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

func testAccGaugeResource(resourceName string, resourceModel gaugeTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_gauge" "%[1]s" {
  type = "indicator"
	 id = "%[2]s"
	 gauge_data_source = "%[3]s"
	 enabled           = %[4]t
}`, resourceName,
		resourceModel.id,
		resourceModel.gaugeDataSource,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedGaugeAttributes(config gaugeTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.GaugeApi.GetGauge(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Gauge"
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
func testAccCheckGaugeDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.GaugeApi.GetGauge(ctx, testIdGauge).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Gauge", testIdGauge)
	}
	return nil
}
