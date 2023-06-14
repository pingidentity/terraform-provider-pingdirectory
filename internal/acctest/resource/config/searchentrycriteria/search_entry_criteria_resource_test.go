package searchentrycriteria_test

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

const testIdSearchEntryCriteria = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type searchEntryCriteriaTestModel struct {
	id string
}

func TestAccSearchEntryCriteria(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := searchEntryCriteriaTestModel{
		id: testIdSearchEntryCriteria,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckSearchEntryCriteriaDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSearchEntryCriteriaResource(resourceName, initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSearchEntryCriteriaResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_search_entry_criteria." + resourceName,
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

func testAccSearchEntryCriteriaResource(resourceName string, resourceModel searchEntryCriteriaTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_search_entry_criteria" "%[1]s" {
  type = "simple"
	 id = "%[2]s"
}`, resourceName,
		resourceModel.id)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSearchEntryCriteriaAttributes(config searchEntryCriteriaTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		_, _, err := testClient.SearchEntryCriteriaApi.GetSearchEntryCriteria(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSearchEntryCriteriaDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.SearchEntryCriteriaApi.GetSearchEntryCriteria(ctx, testIdSearchEntryCriteria).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Search Entry Criteria", testIdSearchEntryCriteria)
	}
	return nil
}
