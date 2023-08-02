package resultcodemap_test

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

const testIdResultCodeMap = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type resultCodeMapTestModel struct {
	id          string
	description string
}

func TestAccResultCodeMap(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := resultCodeMapTestModel{
		id:          testIdResultCodeMap,
		description: "mapping my codes",
	}
	updatedResourceModel := resultCodeMapTestModel{
		id:          testIdResultCodeMap,
		description: "mapping my codes again",
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckResultCodeMapDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccResultCodeMapResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedResultCodeMapAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_result_code_map.%s", resourceName), "description", initialResourceModel.description),
					resource.TestCheckResourceAttrSet("data.pingdirectory_result_code_maps.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccResultCodeMapResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedResultCodeMapAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccResultCodeMapResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_result_code_map." + resourceName,
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

func testAccResultCodeMapResource(resourceName string, resourceModel resultCodeMapTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_result_code_map" "%[1]s" {
  name        = "%[2]s"
  description = "%[3]s"
}

data "pingdirectory_result_code_map" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_result_code_map.%[1]s
  ]
}

data "pingdirectory_result_code_maps" "list" {
  depends_on = [
    pingdirectory_result_code_map.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedResultCodeMapAttributes(config resultCodeMapTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "Result Code Map"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ResultCodeMapApi.GetResultCodeMap(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, nil, "description", config.description, *response.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckResultCodeMapDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ResultCodeMapApi.GetResultCodeMap(ctx, testIdResultCodeMap).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Result Code Map", testIdResultCodeMap)
	}
	return nil
}
