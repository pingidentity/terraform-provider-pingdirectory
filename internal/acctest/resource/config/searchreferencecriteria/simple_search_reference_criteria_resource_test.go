package searchreferencecriteria_test

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

const testIdSimpleSearchReferenceCriteria = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type simpleSearchReferenceCriteriaTestModel struct {
	id          string
	description string
}

func TestAccSimpleSearchReferenceCriteria(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := simpleSearchReferenceCriteriaTestModel{
		id:          testIdSimpleSearchReferenceCriteria,
		description: "my_description",
	}
	updatedResourceModel := simpleSearchReferenceCriteriaTestModel{
		id:          testIdSimpleSearchReferenceCriteria,
		description: "my_updated_description",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckSimpleSearchReferenceCriteriaDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSimpleSearchReferenceCriteriaResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSimpleSearchReferenceCriteriaAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccSimpleSearchReferenceCriteriaResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSimpleSearchReferenceCriteriaAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSimpleSearchReferenceCriteriaResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_search_reference_criteria." + resourceName,
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

func testAccSimpleSearchReferenceCriteriaResource(resourceName string, resourceModel simpleSearchReferenceCriteriaTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_search_reference_criteria" "%[1]s" {
	type = "simple"
  id          = "%[2]s"
  description = "%[3]s"
}`, resourceName,
		resourceModel.id,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSimpleSearchReferenceCriteriaAttributes(config simpleSearchReferenceCriteriaTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.SearchReferenceCriteriaApi.GetSearchReferenceCriteria(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		resourceType := "Search Reference Criteria"
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "description",
			config.description, response.SimpleSearchReferenceCriteriaResponse.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSimpleSearchReferenceCriteriaDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.SearchReferenceCriteriaApi.GetSearchReferenceCriteria(ctx, testIdSimpleSearchReferenceCriteria).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Simple Search Reference Criteria", testIdSimpleSearchReferenceCriteria)
	}
	return nil
}
