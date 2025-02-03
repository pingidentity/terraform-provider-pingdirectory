// Copyright © 2025 Ping Identity Corporation

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
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccServerGroupResource(resourceName, initialResourceModel),
				Check:  resource.TestCheckResourceAttrSet("data.pingdirectory_server_groups.list", "ids.0"),
			},
			{
				// Test importing the resource
				Config:            testAccServerGroupResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_server_group." + resourceName,
				ImportStateId:     initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.ServerGroupAPI.DeleteServerGroup(ctx, initialResourceModel.id).Execute()
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

func testAccServerGroupResource(resourceName string, resourceModel serverGroupTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_server_group" "%[1]s" {
  name = "%[2]s"
}

data "pingdirectory_server_group" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_server_group.%[1]s
  ]
}

data "pingdirectory_server_groups" "list" {
  depends_on = [
    pingdirectory_server_group.%[1]s
  ]
}`, resourceName,
		resourceModel.id)
}

// Test that any objects created by the test are destroyed
func testAccCheckServerGroupDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ServerGroupAPI.GetServerGroup(ctx, testIdServerGroup).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Server Group", testIdServerGroup)
	}
	return nil
}
