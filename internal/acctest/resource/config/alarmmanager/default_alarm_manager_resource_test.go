// Copyright © 2025 Ping Identity Corporation

package alarmmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingdirectory-go-client/v10300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

// Attributes to test with. Add optional properties to test here if desired.
type alarmManagerTestModel struct {
	defaultGaugeAlertLevel string
	generatedAlertTypes    []string
}

func TestAccAlarmManager(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := alarmManagerTestModel{
		defaultGaugeAlertLevel: "critical-only",
		generatedAlertTypes:    []string{"standard"},
	}
	updatedResourceModel := alarmManagerTestModel{
		defaultGaugeAlertLevel: "always",
		generatedAlertTypes:    []string{"alarm"},
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				Config: testAccAlarmManagerResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedAlarmManagerAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_alarm_manager.%s", resourceName), "default_gauge_alert_level", initialResourceModel.defaultGaugeAlertLevel),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_alarm_manager.%s", resourceName), "generated_alert_types.*", initialResourceModel.generatedAlertTypes[0]),
				),
			},
			{
				// Test updating some fields
				Config: testAccAlarmManagerResource(resourceName, updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedAlarmManagerAttributes(updatedResourceModel),
				),
			},
			{
				// Test importing the resource
				Config:       testAccAlarmManagerResource(resourceName, initialResourceModel),
				ResourceName: "pingdirectory_default_alarm_manager." + resourceName,
				// The id doesn't matter for singleton config objects
				ImportStateId:     resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAlarmManagerResource(resourceName string, resourceModel alarmManagerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_alarm_manager" "%[1]s" {
  default_gauge_alert_level = "%[2]s"
  generated_alert_types     = %[3]s
}

data "pingdirectory_alarm_manager" "%[1]s" {
  depends_on = [
    pingdirectory_default_alarm_manager.%[1]s
  ]
}`, resourceName,
		resourceModel.defaultGaugeAlertLevel,
		acctest.StringSliceToTerraformString(resourceModel.generatedAlertTypes))
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedAlarmManagerAttributes(config alarmManagerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "alarm manager"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AlarmManagerAPI.GetAlarmManager(ctx).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, nil, "default-gauge-alert-level", config.defaultGaugeAlertLevel, response.DefaultGaugeAlertLevel.String())
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, nil, "generated-alert-types", config.generatedAlertTypes, client.StringSliceEnumalarmManagerGeneratedAlertTypesProp(response.GeneratedAlertTypes))
		if err != nil {
			return err
		}
		return nil
	}
}
