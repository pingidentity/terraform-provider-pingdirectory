package config_test

import (
	"fmt"
	"testing"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testIdTopologyAdminUser = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type topologyAdminUserTestModel struct {
	id string
}

func TestAccTopologyAdminUser(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := topologyAdminUserTestModel{
		id: testIdTopologyAdminUser,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckTopologyAdminUserDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccTopologyAdminUserResource(resourceName, initialResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccTopologyAdminUserResource(resourceName, initialResourceModel),
				ResourceName:            "pingdirectory_topology_admin_user." + resourceName,
				ImportStateId:           initialResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccTopologyAdminUserResource(resourceName string, resourceModel topologyAdminUserTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_topology_admin_user" "%[1]s" {
	 id = "%[2]s"
}`, resourceName, resourceModel.id)
}

// Test that any objects created by the test are destroyed
func testAccCheckTopologyAdminUserDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.TopologyAdminUserApi.GetTopologyAdminUser(ctx, testIdTopologyAdminUser).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Topology Admin User", testIdTopologyAdminUser)
	}
	return nil
}
