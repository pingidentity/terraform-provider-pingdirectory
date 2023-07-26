package uncachedentrycriteria_test

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

const testIdDefaultUncachedEntryCriteria = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type defaultUncachedEntryCriteriaTestModel struct {
	id          string
	description string
	enabled     bool
}

func TestAccDefaultUncachedEntryCriteria(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := defaultUncachedEntryCriteriaTestModel{
		id:          testIdDefaultUncachedEntryCriteria,
		description: "initial description",
		enabled:     false,
	}
	updatedResourceModel := defaultUncachedEntryCriteriaTestModel{
		id:          testIdDefaultUncachedEntryCriteria,
		description: "updated description",
		enabled:     true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckDefaultUncachedEntryCriteriaDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDefaultUncachedEntryCriteriaResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedDefaultUncachedEntryCriteriaAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_uncached_entry_criteria.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
				),
			},
			{
				// Test updating some fields
				Config: testAccDefaultUncachedEntryCriteriaResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedDefaultUncachedEntryCriteriaAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccDefaultUncachedEntryCriteriaResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_uncached_entry_criteria." + resourceName,
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

func testAccDefaultUncachedEntryCriteriaResource(resourceName string, resourceModel defaultUncachedEntryCriteriaTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_uncached_entry_criteria" "%[1]s" {
  type        = "default"
  id          = "%[2]s"
  description = "%[3]s"
  enabled     = %[4]t
}

data "pingdirectory_uncached_entry_criteria" "%[1]s" {
	 id = "%[2]s"
  depends_on = [
    pingdirectory_uncached_entry_criteria.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.description,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDefaultUncachedEntryCriteriaAttributes(config defaultUncachedEntryCriteriaTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.UncachedEntryCriteriaApi.GetUncachedEntryCriteria(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Default Uncached Entry Criteria"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.DefaultUncachedEntryCriteriaResponse.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "description",
			config.description, response.DefaultUncachedEntryCriteriaResponse.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckDefaultUncachedEntryCriteriaDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.UncachedEntryCriteriaApi.GetUncachedEntryCriteria(ctx, testIdDefaultUncachedEntryCriteria).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Default Uncached Entry Criteria", testIdDefaultUncachedEntryCriteria)
	}
	return nil
}
