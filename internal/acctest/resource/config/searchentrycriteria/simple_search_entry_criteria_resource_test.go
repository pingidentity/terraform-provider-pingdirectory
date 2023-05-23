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

const testIdSimpleSearchEntryCriteria = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type simpleSearchEntryCriteriaTestModel struct {
	id          string
	description string
}

func TestAccSimpleSearchEntryCriteria(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := simpleSearchEntryCriteriaTestModel{
		id:          testIdSimpleSearchEntryCriteria,
		description: "my_description",
	}
	updatedResourceModel := simpleSearchEntryCriteriaTestModel{
		id:          testIdSimpleSearchEntryCriteria,
		description: "my_updated_description",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckSimpleSearchEntryCriteriaDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSimpleSearchEntryCriteriaResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSimpleSearchEntryCriteriaAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccSimpleSearchEntryCriteriaResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSimpleSearchEntryCriteriaAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSimpleSearchEntryCriteriaResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_simple_search_entry_criteria." + resourceName,
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

func testAccSimpleSearchEntryCriteriaResource(resourceName string, resourceModel simpleSearchEntryCriteriaTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_simple_search_entry_criteria" "%[1]s" {
  id          = "%[2]s"
  description = "%[3]s"
}`, resourceName,
		resourceModel.id,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSimpleSearchEntryCriteriaAttributes(config simpleSearchEntryCriteriaTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.SearchEntryCriteriaApi.GetSearchEntryCriteria(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		resourceType := "Search Entry Criteria"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "description",
			config.description, *response.SimpleSearchEntryCriteriaResponse.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSimpleSearchEntryCriteriaDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.SearchEntryCriteriaApi.GetSearchEntryCriteria(ctx, testIdSimpleSearchEntryCriteria).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Simple Search Entry Criteria", testIdSimpleSearchEntryCriteria)
	}
	return nil
}
