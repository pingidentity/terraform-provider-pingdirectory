package monitoringendpoint_test

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

const testIdStatsdMonitoringEndpoint = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type statsdMonitoringEndpointTestModel struct {
	id       string
	hostname string
	enabled  bool
}

func TestAccStatsdMonitoringEndpoint(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := statsdMonitoringEndpointTestModel{
		id:       testIdStatsdMonitoringEndpoint,
		hostname: "example.com",
		enabled:  true,
	}
	updatedResourceModel := statsdMonitoringEndpointTestModel{
		id:       testIdStatsdMonitoringEndpoint,
		hostname: "example.org",
		enabled:  false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckStatsdMonitoringEndpointDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccStatsdMonitoringEndpointResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedStatsdMonitoringEndpointAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccStatsdMonitoringEndpointResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedStatsdMonitoringEndpointAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccStatsdMonitoringEndpointResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_monitoring_endpoint." + resourceName,
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

func testAccStatsdMonitoringEndpointResource(resourceName string, resourceModel statsdMonitoringEndpointTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_monitoring_endpoint" "%[1]s" {
  id       = "%[2]s"
  hostname = "%[3]s"
  enabled  = %[4]t
}`, resourceName,
		resourceModel.id,
		resourceModel.hostname,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedStatsdMonitoringEndpointAttributes(config statsdMonitoringEndpointTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.MonitoringEndpointApi.GetMonitoringEndpoint(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Statsd Monitoring Endpoint"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "hostname",
			config.hostname, response.Hostname)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckStatsdMonitoringEndpointDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.MonitoringEndpointApi.GetMonitoringEndpoint(ctx, testIdStatsdMonitoringEndpoint).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Statsd Monitoring Endpoint", testIdStatsdMonitoringEndpoint)
	}
	return nil
}
