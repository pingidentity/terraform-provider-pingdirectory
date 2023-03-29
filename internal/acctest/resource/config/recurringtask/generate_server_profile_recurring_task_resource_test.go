package recurringtask_test

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

const testIdGenerateServerProfileRecurringTask = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type generateServerProfileRecurringTaskTestModel struct {
	id                         string
	profileDirectory           string
	retainPreviousProfileCount int64
}

func TestAccGenerateServerProfileRecurringTask(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := generateServerProfileRecurringTaskTestModel{
		id:                         testIdGenerateServerProfileRecurringTask,
		profileDirectory:           "/opt/out/instance",
		retainPreviousProfileCount: 10,
	}
	updatedResourceModel := generateServerProfileRecurringTaskTestModel{
		id:                         testIdGenerateServerProfileRecurringTask,
		profileDirectory:           "/opt/out",
		retainPreviousProfileCount: 11,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckGenerateServerProfileRecurringTaskDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccGenerateServerProfileRecurringTaskResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedGenerateServerProfileRecurringTaskAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccGenerateServerProfileRecurringTaskResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedGenerateServerProfileRecurringTaskAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccGenerateServerProfileRecurringTaskResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_generate_server_profile_recurring_task." + resourceName,
				ImportStateId:           updatedResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccGenerateServerProfileRecurringTaskResource(resourceName string, resourceModel generateServerProfileRecurringTaskTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_generate_server_profile_recurring_task" "%[1]s" {
  id                            = "%[2]s"
  profile_directory             = "%[3]s"
  retain_previous_profile_count = "%[4]d"
}`, resourceName, resourceModel.id,
		resourceModel.profileDirectory, resourceModel.retainPreviousProfileCount)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedGenerateServerProfileRecurringTaskAttributes(config generateServerProfileRecurringTaskTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.RecurringTaskApi.GetRecurringTask(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Generate Server Profile Recurring Task"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "profile-directory",
			config.profileDirectory, response.GenerateServerProfileRecurringTaskResponse.ProfileDirectory)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckGenerateServerProfileRecurringTaskDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.RecurringTaskApi.GetRecurringTask(ctx, testIdGenerateServerProfileRecurringTask).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Generate Server Profile Recurring Task", testIdGenerateServerProfileRecurringTask)
	}
	return nil
}
