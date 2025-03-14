// Copyright © 2025 Ping Identity Corporation

package serverinstancelistener_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdServerInstanceListener = "ldap-listener-mirrored-config"

// Attributes to test with. Add optional properties to test here if desired.
type serverInstanceListenerTestModel struct {
	id                 string
	serverInstanceName string
}

func TestAccServerInstanceListener(t *testing.T) {
	resourceName := "myresource"
	// Figure out the instance name of the test server, so we can use that instance
	var instanceName string
	// Only run for acceptance tests
	if os.Getenv("TF_ACC") == "1" {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.GlobalConfigurationAPI.GetGlobalConfiguration(ctx).Execute()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		instanceName = response.InstanceName
	}
	initialResourceModel := serverInstanceListenerTestModel{
		id:                 testIdServerInstanceListener,
		serverInstanceName: instanceName,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccServerInstanceListenerResource(resourceName, initialResourceModel),
				Check:  resource.TestCheckResourceAttrSet("data.pingdirectory_server_instance_listeners.list", "objects.0.id"),
			},
			{
				// Test importing the resource
				Config:            testAccServerInstanceListenerResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_default_server_instance_listener." + resourceName,
				ImportStateId:     initialResourceModel.serverInstanceName + "/" + initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccServerInstanceListenerResource(resourceName string, resourceModel serverInstanceListenerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_server_instance_listener" "%[1]s" {
  name                 = "%[2]s"
  server_instance_name = "%[3]s"
}

data "pingdirectory_server_instance_listener" "%[1]s" {
  name                 = "%[2]s"
  server_instance_name = "%[3]s"
  depends_on = [
    pingdirectory_default_server_instance_listener.%[1]s
  ]
}

data "pingdirectory_server_instance_listeners" "list" {
  server_instance_name = "%[3]s"
  depends_on = [
    pingdirectory_default_server_instance_listener.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.serverInstanceName)
}
