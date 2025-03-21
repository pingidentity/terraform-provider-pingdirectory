// Copyright © 2025 Ping Identity Corporation

package serverinstance_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

// Some attributes to test with
type testModel struct {
	jmxPort         int64
	startTlsEnabled bool
}

func TestAccDirectoryServerInstance(t *testing.T) {
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
	resourceName := "instance"

	initialResourceModel := testModel{
		jmxPort:         1111,
		startTlsEnabled: true,
	}
	updatedResourceModel := testModel{
		jmxPort:         1112,
		startTlsEnabled: false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				Config: testAccDirectoryserverInstanceResource(resourceName, instanceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedDirectoryServerInstanceAttributes(instanceName, initialResourceModel),
					// Check some computed attributes are set as expected (PingDirectory defaults)
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_default_server_instance.%s", resourceName), "preferred_security", "ssl"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_default_server_instance.%s", resourceName), "ldap_port", "1389"),
					// Check those attributes on the data source
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_server_instance.%s", resourceName), "preferred_security", "ssl"),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_server_instance.%s", resourceName), "ldap_port", "1389"),
					resource.TestCheckResourceAttrSet("data.pingdirectory_server_instances.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccDirectoryserverInstanceResource(resourceName, instanceName, updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedDirectoryServerInstanceAttributes(instanceName, updatedResourceModel),
				),
			},
			{
				// Test importing the resource
				Config:            testAccDirectoryserverInstanceResource(resourceName, instanceName, updatedResourceModel),
				ResourceName:      "pingdirectory_default_server_instance." + resourceName,
				ImportStateId:     instanceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDirectoryserverInstanceResource(resourceName, instanceName string, resourceModel testModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_server_instance" "%[1]s" {
  name                 = "%[2]s"
  server_instance_name = "%[2]s"
  jmx_port             = %[3]d
  start_tls_enabled    = %[4]t
}

data "pingdirectory_server_instance" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_default_server_instance.%[1]s
  ]
}

data "pingdirectory_server_instances" "list" {
  depends_on = [
    pingdirectory_default_server_instance.%[1]s
  ]
}`, resourceName, instanceName, resourceModel.jmxPort, resourceModel.startTlsEnabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDirectoryServerInstanceAttributes(instanceName string, config testModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "directory server instance"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ServerInstanceAPI.GetServerInstance(ctx, instanceName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchInt(resourceType, &instanceName, "jmx-port",
			config.jmxPort, *response.DirectoryServerInstanceResponse.JmxPort)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &instanceName, "start-tls-enabled",
			config.startTlsEnabled, *response.DirectoryServerInstanceResponse.StartTLSEnabled)
		if err != nil {
			return err
		}
		return nil
	}
}
