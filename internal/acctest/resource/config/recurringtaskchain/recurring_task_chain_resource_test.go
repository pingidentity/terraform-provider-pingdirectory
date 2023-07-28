package recurringtaskchain_test

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

const testIdRecurringTaskChain = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type recurringTaskChainTestModel struct {
	id                         string
	recurringTask              []string
	scheduledDateSelectionType string
	scheduledTimeOfDay         []string
}

func TestAccRecurringTaskChain(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := recurringTaskChainTestModel{
		id: testIdRecurringTaskChain,
		recurringTask: []string{
			"Delete Old Expensive Operation Dumps",
		},
		scheduledDateSelectionType: "every-day",
		scheduledTimeOfDay: []string{
			"02:00", "03:15",
		},
	}
	updatedResourceModel := recurringTaskChainTestModel{
		id: testIdRecurringTaskChain,
		recurringTask: []string{
			"Delete Old Lock Conflict Details Log Files",
			"Delete Old Work Queue Backlog Thread Dumps",
		},
		scheduledDateSelectionType: "every-day",
		scheduledTimeOfDay: []string{
			"01:00",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckRecurringTaskChainDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccRecurringTaskChainResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedRecurringTaskChainAttributes(initialResourceModel),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_recurring_task_chain.%s", resourceName), "recurring_task.*", initialResourceModel.recurringTask[0]),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_recurring_task_chain.%s", resourceName), "scheduled_date_selection_type", initialResourceModel.scheduledDateSelectionType),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_recurring_task_chain.%s", resourceName), "scheduled_time_of_day.*", initialResourceModel.scheduledTimeOfDay[0]),
				),
			},
			{
				// Test updating some fields
				Config: testAccRecurringTaskChainResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedRecurringTaskChainAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccRecurringTaskChainResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_recurring_task_chain." + resourceName,
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

func testAccRecurringTaskChainResource(resourceName string, resourceModel recurringTaskChainTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_recurring_task_chain" "%[1]s" {
  id                            = "%[2]s"
  recurring_task                = %[3]s
  scheduled_date_selection_type = "%[4]s"
  scheduled_time_of_day         = %[5]s
}

data "pingdirectory_recurring_task_chain" "%[1]s" {
  id = "%[2]s"
  depends_on = [
    pingdirectory_recurring_task_chain.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		acctest.StringSliceToTerraformString(resourceModel.recurringTask),
		resourceModel.scheduledDateSelectionType,
		acctest.StringSliceToTerraformString(resourceModel.scheduledTimeOfDay))
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedRecurringTaskChainAttributes(config recurringTaskChainTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.RecurringTaskChainApi.GetRecurringTaskChain(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Recurring Task Chain"
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "recurring-task",
			config.recurringTask, response.RecurringTask)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "scheduled-date-selection-type",
			config.scheduledDateSelectionType, response.ScheduledDateSelectionType.String())
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "scheduled-time-of-day",
			config.scheduledTimeOfDay, response.ScheduledTimeOfDay)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckRecurringTaskChainDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.RecurringTaskChainApi.GetRecurringTaskChain(ctx, testIdRecurringTaskChain).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Recurring Task Chain", testIdRecurringTaskChain)
	}
	return nil
}
