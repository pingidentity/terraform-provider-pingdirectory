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

const testIdResultCriteria = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type resultCriteriaTestModel struct {
	id string
}

func TestAccResultCriteria(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := resultCriteriaTestModel{
		id: testIdResultCriteria,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckResultCriteriaDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccResultCriteriaResource(resourceName, initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccResultCriteriaResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_result_criteria." + resourceName,
				ImportStateId:     initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccResultCriteriaResource(resourceName string, resourceModel resultCriteriaTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_result_criteria" "%[1]s" {
  type = "simple"
	 id = "%[2]s"
}`, resourceName,
		resourceModel.id)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedResultCriteriaAttributes(config resultCriteriaTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		_, _, err := testClient.ResultCriteriaApi.GetResultCriteria(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckResultCriteriaDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ResultCriteriaApi.GetResultCriteria(ctx, testIdResultCriteria).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Result Criteria", testIdResultCriteria)
	}
	return nil
}
