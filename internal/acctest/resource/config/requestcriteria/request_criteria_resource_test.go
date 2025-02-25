// Copyright © 2025 Ping Identity Corporation

package requestcriteria_test

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

const testIdRootDseRequestCriteria = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type rootDseRequestCriteriaTestModel struct {
	id          string
	description string
}

func TestAccRootDseRequestCriteria(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := rootDseRequestCriteriaTestModel{
		id:          testIdRootDseRequestCriteria,
		description: "test description",
	}

	updatedResourceModel := rootDseRequestCriteriaTestModel{
		id:          testIdRootDseRequestCriteria,
		description: "updated test description",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckRootDseRequestCriteriaDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccRootDseRequestCriteriaResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedRootDseRequestCriteriaAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_request_criteria.%s", resourceName), "description", initialResourceModel.description),
					resource.TestCheckResourceAttrSet("data.pingdirectory_request_criteria_list.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccRootDseRequestCriteriaResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedRootDseRequestCriteriaAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccRootDseRequestCriteriaResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_request_criteria." + resourceName,
				ImportStateId:     initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.RequestCriteriaAPI.DeleteRequestCriteria(ctx, updatedResourceModel.id).Execute()
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

func testAccRootDseRequestCriteriaResource(resourceName string, resourceModel rootDseRequestCriteriaTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_request_criteria" "%[1]s" {
  type        = "root-dse"
  name        = "%[2]s"
  description = "%[3]s"
}

data "pingdirectory_request_criteria" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_request_criteria.%[1]s
  ]
}

data "pingdirectory_request_criteria_list" "list" {
  depends_on = [
    pingdirectory_request_criteria.%[1]s
  ]
}`, resourceName, resourceModel.id, resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedRootDseRequestCriteriaAttributes(config rootDseRequestCriteriaTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "root dse request criteria"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.RequestCriteriaAPI.GetRequestCriteria(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that description matches expected
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "description",
			config.description, response.RootDseRequestCriteriaResponse.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckRootDseRequestCriteriaDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.RequestCriteriaAPI.GetRequestCriteria(ctx, testIdRootDseRequestCriteria).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Root Dse Request Criteria", testIdRootDseRequestCriteria)
	}
	return nil
}
