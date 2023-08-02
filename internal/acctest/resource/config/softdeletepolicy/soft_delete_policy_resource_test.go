package softdeletepolicy_test

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

const testIdSoftDeletePolicy = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type softDeletePolicyTestModel struct {
	id          string
	description string
}

func TestAccSoftDeletePolicy(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := softDeletePolicyTestModel{
		id:          testIdSoftDeletePolicy,
		description: "initial Description",
	}
	updatedResourceModel := softDeletePolicyTestModel{
		id:          testIdSoftDeletePolicy,
		description: "updated Description",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckSoftDeletePolicyDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSoftDeletePolicyResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedSoftDeletePolicyAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_soft_delete_policy.%s", resourceName), "description", initialResourceModel.description),
					resource.TestCheckResourceAttrSet("data.pingdirectory_soft_delete_policies.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccSoftDeletePolicyResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSoftDeletePolicyAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSoftDeletePolicyResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_soft_delete_policy." + resourceName,
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

func testAccSoftDeletePolicyResource(resourceName string, resourceModel softDeletePolicyTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_soft_delete_policy" "%[1]s" {
  name        = "%[2]s"
  description = "%[3]s"
}

data "pingdirectory_soft_delete_policy" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_soft_delete_policy.%[1]s
  ]
}

data "pingdirectory_soft_delete_policies" "list" {
  depends_on = [
    pingdirectory_soft_delete_policy.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSoftDeletePolicyAttributes(config softDeletePolicyTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.SoftDeletePolicyApi.GetSoftDeletePolicy(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify the attributes have expected values
		resourceType := "Soft Delete Policy"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "description",
			config.description, *response.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSoftDeletePolicyDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.SoftDeletePolicyApi.GetSoftDeletePolicy(ctx, testIdSoftDeletePolicy).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Soft Delete Policy", testIdSoftDeletePolicy)
	}
	return nil
}
