package prometheusmonitorattributemetric_test

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

const testIdPrometheusMonitorAttributeMetric = "MyId"
const testPrometheusHttpServletExtensionName = "Prometheus Monitoring"

// Attributes to test with. Add optional properties to test here if desired.
type prometheusMonitorAttributeMetricTestModel struct {
	httpServletExtensionName string
	metricName               string
	monitorAttributeName     string
	monitorObjectClassName   string
	metricType               string
}

func TestAccPrometheusMonitorAttributeMetric(t *testing.T) {
	pdVersion := os.Getenv("PINGDIRECTORY_PROVIDER_PRODUCT_VERSION")
	compare, err := version.Compare(pdVersion, version.PingDirectory9200)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if compare < 0 {
		// This resource only exists in PD version 9.2 and later
		return
	}

	resourceName := "myresource"
	initialResourceModel := prometheusMonitorAttributeMetricTestModel{
		httpServletExtensionName: testPrometheusHttpServletExtensionName,
		metricName:               testIdPrometheusMonitorAttributeMetric,
		monitorAttributeName:     "mymonitorattr",
		monitorObjectClassName:   "ds-cfg-monitor",
		metricType:               "gauge",
	}
	updatedResourceModel := prometheusMonitorAttributeMetricTestModel{
		httpServletExtensionName: testPrometheusHttpServletExtensionName,
		metricName:               testIdPrometheusMonitorAttributeMetric,
		monitorAttributeName:     "mymonitorattr",
		monitorObjectClassName:   "ds-cfg-monitor",
		metricType:               "counter",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckPrometheusMonitorAttributeMetricDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccPrometheusMonitorAttributeMetricResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedPrometheusMonitorAttributeMetricAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_prometheus_monitor_attribute_metric.%s", resourceName), "metric_name", initialResourceModel.metricName),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_prometheus_monitor_attribute_metric.%s", resourceName), "monitor_attribute_name", initialResourceModel.monitorAttributeName),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_prometheus_monitor_attribute_metric.%s", resourceName), "monitor_object_class_name", initialResourceModel.monitorObjectClassName),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_prometheus_monitor_attribute_metric.%s", resourceName), "metric_type", initialResourceModel.metricType),
					resource.TestCheckResourceAttrSet("data.pingdirectory_prometheus_monitor_attribute_metrics.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccPrometheusMonitorAttributeMetricResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedPrometheusMonitorAttributeMetricAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccPrometheusMonitorAttributeMetricResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_prometheus_monitor_attribute_metric." + resourceName,
				ImportStateId:     updatedResourceModel.httpServletExtensionName + "/" + updatedResourceModel.metricName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.PrometheusMonitorAttributeMetricApi.DeletePrometheusMonitorAttributeMetric(ctx, updatedResourceModel.metricName, updatedResourceModel.httpServletExtensionName).Execute()
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

func testAccPrometheusMonitorAttributeMetricResource(resourceName string, resourceModel prometheusMonitorAttributeMetricTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_prometheus_monitor_attribute_metric" "%[1]s" {
  http_servlet_extension_name = "%[2]s"
  metric_name                 = "%[3]s"
  monitor_attribute_name      = "%[4]s"
  monitor_object_class_name   = "%[5]s"
  metric_type                 = "%[6]s"
}

data "pingdirectory_prometheus_monitor_attribute_metric" "%[1]s" {
  http_servlet_extension_name = "%[2]s"
  metric_name                 = "%[3]s"
  depends_on = [
    pingdirectory_prometheus_monitor_attribute_metric.%[1]s
  ]
}

data "pingdirectory_prometheus_monitor_attribute_metrics" "list" {
  http_servlet_extension_name = "%[2]s"
  depends_on = [
    pingdirectory_prometheus_monitor_attribute_metric.%[1]s
  ]
}`, resourceName,
		resourceModel.httpServletExtensionName,
		resourceModel.metricName,
		resourceModel.monitorAttributeName,
		resourceModel.monitorObjectClassName,
		resourceModel.metricType)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedPrometheusMonitorAttributeMetricAttributes(config prometheusMonitorAttributeMetricTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.PrometheusMonitorAttributeMetricApi.GetPrometheusMonitorAttributeMetric(ctx, config.metricName, config.httpServletExtensionName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Prometheus Monitor Attribute Metric"
		err = acctest.TestAttributesMatchString(resourceType, &config.metricName, "metric-name",
			config.metricName, response.MetricName)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.metricName, "monitor-attribute-name",
			config.monitorAttributeName, response.MonitorAttributeName)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.metricName, "monitor-object-class-name",
			config.monitorObjectClassName, response.MonitorObjectClassName)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.metricName, "metric-type",
			config.metricType, response.MetricType.String())
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckPrometheusMonitorAttributeMetricDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.PrometheusMonitorAttributeMetricApi.GetPrometheusMonitorAttributeMetric(ctx, testIdPrometheusMonitorAttributeMetric, testPrometheusHttpServletExtensionName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Prometheus Monitor Attribute Metric", testIdPrometheusMonitorAttributeMetric)
	}
	return nil
}
