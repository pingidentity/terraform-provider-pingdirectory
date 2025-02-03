// Copyright © 2025 Ping Identity Corporation

package uncachedattributecriteria_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdDefaultUncachedAttributeCriteria = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type defaultUncachedAttributeCriteriaTestModel struct {
	id          string
	description string
	enabled     bool
}

func TestAccDefaultUncachedAttributeCriteria(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := defaultUncachedAttributeCriteriaTestModel{
		id:          testIdDefaultUncachedAttributeCriteria,
		description: "initial description",
		enabled:     false,
	}
	updatedResourceModel := defaultUncachedAttributeCriteriaTestModel{
		id:          testIdDefaultUncachedAttributeCriteria,
		description: "updated description",
		enabled:     false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckDefaultUncachedAttributeCriteriaDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDefaultUncachedAttributeCriteriaResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedDefaultUncachedAttributeCriteriaAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_uncached_attribute_criteria.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_uncached_attribute_criteria_list.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccDefaultUncachedAttributeCriteriaResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedDefaultUncachedAttributeCriteriaAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccDefaultUncachedAttributeCriteriaResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_uncached_attribute_criteria." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.UncachedAttributeCriteriaAPI.DeleteUncachedAttributeCriteria(ctx, updatedResourceModel.id).Execute()
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

func testAccDefaultUncachedAttributeCriteriaResource(resourceName string, resourceModel defaultUncachedAttributeCriteriaTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_uncached_attribute_criteria" "%[1]s" {
  type        = "default"
  name        = "%[2]s"
  description = "%[3]s"
  enabled     = %[4]t
}

data "pingdirectory_uncached_attribute_criteria" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_uncached_attribute_criteria.%[1]s
  ]
}

data "pingdirectory_uncached_attribute_criteria_list" "list" {
  depends_on = [
    pingdirectory_uncached_attribute_criteria.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.description,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDefaultUncachedAttributeCriteriaAttributes(config defaultUncachedAttributeCriteriaTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.UncachedAttributeCriteriaAPI.GetUncachedAttributeCriteria(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Default Uncached Attribute Criteria"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.DefaultUncachedAttributeCriteriaResponse.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "description",
			config.description, *response.DefaultUncachedAttributeCriteriaResponse.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckDefaultUncachedAttributeCriteriaDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.UncachedAttributeCriteriaAPI.GetUncachedAttributeCriteria(ctx, testIdDefaultUncachedAttributeCriteria).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Default Uncached Attribute Criteria", testIdDefaultUncachedAttributeCriteria)
	}
	return nil
}
