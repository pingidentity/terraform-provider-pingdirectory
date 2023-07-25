package servergroup_test

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

const testIdServerGroup = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type serverGroupTestModel struct {
	id string
}

func TestAccServerGroup(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := serverGroupTestModel{
		id: testIdServerGroup,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccServerGroupResource(resourceName, initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccServerGroupResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_server_group." + resourceName,
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

func testAccServerGroupResource(resourceName string, resourceModel serverGroupTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_server_group" "%[1]s" {
  id = "%[2]s"
}`, resourceName,
		resourceModel.id)
}

// Test that any objects created by the test are destroyed
func testAccCheckServerGroupDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ServerGroupApi.GetServerGroup(ctx, testIdServerGroup).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Server Group", testIdServerGroup)
	}
	return nil
}
