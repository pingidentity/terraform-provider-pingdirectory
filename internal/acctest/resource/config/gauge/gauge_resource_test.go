// Copyright © 2025 Ping Identity Corporation

package gauge_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

const testIdGauge = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type gaugeTestModel struct {
	id              string
	gaugeDataSource string
	enabled         bool
}

func TestAccGauge(t *testing.T) {
	initialGaugeDataSource := post101InitialDataSource
	pdVersion := os.Getenv("PINGDIRECTORY_PROVIDER_PRODUCT_VERSION")
	compare, err := version.Compare(pdVersion, version.PingDirectory10100)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if compare < 0 {
		// The strong encryption data source was changed in version 10.1
		initialGaugeDataSource = pre101InitialDataSource
	}

	resourceName := "myresource"
	initialResourceModel := gaugeTestModel{
		id:              testIdGauge,
		gaugeDataSource: initialGaugeDataSource,
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
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckGaugeDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccGaugeResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedGaugeAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_gauge.%s", resourceName), "gauge_data_source", initialResourceModel.gaugeDataSource),
					resource.TestCheckResourceAttrSet("data.pingdirectory_gauges.list", "objects.0.id"),
				),
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
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.GaugeAPI.DeleteGauge(ctx, updatedResourceModel.id).Execute()
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

func testAccGaugeResource(resourceName string, resourceModel gaugeTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_gauge" "%[1]s" {
  type              = "indicator"
  name              = "%[2]s"
  gauge_data_source = "%[3]s"
  enabled           = %[4]t
}

data "pingdirectory_gauge" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_gauge.%[1]s
  ]
}

data "pingdirectory_gauges" "list" {
  depends_on = [
    pingdirectory_gauge.%[1]s
  ]
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
		response, _, err := testClient.GaugeAPI.GetGauge(ctx, config.id).Execute()
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
	_, _, err := testClient.GaugeAPI.GetGauge(ctx, testIdGauge).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Gauge", testIdGauge)
	}
	return nil
}
