// Copyright © 2025 Ping Identity Corporation

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
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckSimpleSearchEntryCriteriaDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSimpleSearchEntryCriteriaResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedSimpleSearchEntryCriteriaAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_search_entry_criteria.%s", resourceName), "description", initialResourceModel.description),
					resource.TestCheckResourceAttrSet("data.pingdirectory_search_entry_criteria_list.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccSimpleSearchEntryCriteriaResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSimpleSearchEntryCriteriaAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSimpleSearchEntryCriteriaResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_search_entry_criteria." + resourceName,
				ImportStateId:     initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.SearchEntryCriteriaAPI.DeleteSearchEntryCriteria(ctx, updatedResourceModel.id).Execute()
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

func testAccSimpleSearchEntryCriteriaResource(resourceName string, resourceModel simpleSearchEntryCriteriaTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_search_entry_criteria" "%[1]s" {
  type        = "simple"
  name        = "%[2]s"
  description = "%[3]s"
}

data "pingdirectory_search_entry_criteria" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_search_entry_criteria.%[1]s
  ]
}

data "pingdirectory_search_entry_criteria_list" "list" {
  depends_on = [
    pingdirectory_search_entry_criteria.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSimpleSearchEntryCriteriaAttributes(config simpleSearchEntryCriteriaTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.SearchEntryCriteriaAPI.GetSearchEntryCriteria(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		resourceType := "Search Entry Criteria"
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "description",
			config.description, response.SimpleSearchEntryCriteriaResponse.Description)
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
	_, _, err := testClient.SearchEntryCriteriaAPI.GetSearchEntryCriteria(ctx, testIdSimpleSearchEntryCriteria).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Simple Search Entry Criteria", testIdSimpleSearchEntryCriteria)
	}
	return nil
}
