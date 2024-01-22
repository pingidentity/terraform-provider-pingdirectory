package resultcriteria_test

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

const testIdSimpleResultCriteria = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type simpleResultCriteriaTestModel struct {
	id          string
	description string
}

func TestAccSimpleResultCriteria(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := simpleResultCriteriaTestModel{
		id:          testIdSimpleResultCriteria,
		description: "my_description",
	}
	updatedResourceModel := simpleResultCriteriaTestModel{
		id:          testIdSimpleResultCriteria,
		description: "my_updated_description",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckSimpleResultCriteriaDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSimpleResultCriteriaResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedSimpleResultCriteriaAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_result_criteria.%s", resourceName), "description", initialResourceModel.description),
					resource.TestCheckResourceAttrSet("data.pingdirectory_result_criteria_list.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccSimpleResultCriteriaResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSimpleResultCriteriaAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSimpleResultCriteriaResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_result_criteria." + resourceName,
				ImportStateId:     initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.ResultCriteriaAPI.DeleteResultCriteria(ctx, updatedResourceModel.id).Execute()
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

func testAccSimpleResultCriteriaResource(resourceName string, resourceModel simpleResultCriteriaTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_result_criteria" "%[1]s" {
  type        = "simple"
  name        = "%[2]s"
  description = "%[3]s"
}

data "pingdirectory_result_criteria" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_result_criteria.%[1]s
  ]
}

data "pingdirectory_result_criteria_list" "list" {
  depends_on = [
    pingdirectory_result_criteria.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSimpleResultCriteriaAttributes(config simpleResultCriteriaTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ResultCriteriaAPI.GetResultCriteria(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		resourceType := "Result Criteria"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "description",
			config.description, *response.SimpleResultCriteriaResponse.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSimpleResultCriteriaDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ResultCriteriaAPI.GetResultCriteria(ctx, testIdSimpleResultCriteria).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Simple Result Criteria", testIdSimpleResultCriteria)
	}
	return nil
}
