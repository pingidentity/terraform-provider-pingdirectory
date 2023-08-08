package groupimplementation_test

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

const testIdStaticGroupImplementation = "Static"

// Attributes to test with. Add optional properties to test here if desired.
type staticGroupImplementationTestModel struct {
	id      string
	enabled bool
}

func TestAccStaticGroupImplementation(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := staticGroupImplementationTestModel{
		id:      testIdStaticGroupImplementation,
		enabled: true,
	}
	updatedResourceModel := staticGroupImplementationTestModel{
		id:      testIdStaticGroupImplementation,
		enabled: false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccStaticGroupImplementationResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedStaticGroupImplementationAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_group_implementation.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_group_implementations.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccStaticGroupImplementationResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedStaticGroupImplementationAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccStaticGroupImplementationResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_default_group_implementation." + resourceName,
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

func testAccStaticGroupImplementationResource(resourceName string, resourceModel staticGroupImplementationTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_group_implementation" "%[1]s" {
  name    = "%[2]s"
  enabled = %[3]t
}

data "pingdirectory_group_implementation" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_default_group_implementation.%[1]s
  ]
}

data "pingdirectory_group_implementations" "list" {
  depends_on = [
    pingdirectory_default_group_implementation.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedStaticGroupImplementationAttributes(config staticGroupImplementationTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.GroupImplementationApi.GetGroupImplementation(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Static Group Implementation"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.StaticGroupImplementationResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}
