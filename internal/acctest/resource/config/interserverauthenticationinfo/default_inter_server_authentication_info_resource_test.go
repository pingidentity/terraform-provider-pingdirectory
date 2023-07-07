package interserverauthenticationinfo_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdInterServerAuthenticationInfo = "certificate-auth-mirrored-config"
const testServerInstanceListenerName = "ldap-listener-mirrored-config"

// Attributes to test with. Add optional properties to test here if desired.
type interServerAuthenticationInfoTestModel struct {
	id                         string
	serverInstanceListenerName string
	serverInstanceName         string
	purpose                    []string
}

func TestAccInterServerAuthenticationInfo(t *testing.T) {
	// Figure out the instance name of the test server, so we can use that instance
	var instanceName string
	// Only run for acceptance tests
	if os.Getenv("TF_ACC") == "1" {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.GlobalConfigurationApi.GetGlobalConfiguration(ctx).Execute()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		instanceName = response.InstanceName
	}

	resourceName := "myresource"
	initialResourceModel := interServerAuthenticationInfoTestModel{
		id:                         testIdInterServerAuthenticationInfo,
		serverInstanceListenerName: testServerInstanceListenerName,
		serverInstanceName:         instanceName,
		purpose:                    []string{"mirrored-config"},
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
				Config: testAccInterServerAuthenticationInfoResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedInterServerAuthenticationInfoAttributes(initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccInterServerAuthenticationInfoResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_default_inter_server_authentication_info." + resourceName,
				ImportStateId:     initialResourceModel.serverInstanceName + "/" + initialResourceModel.serverInstanceListenerName + "/" + initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
					"password",
				},
			},
		},
	})
}

func testAccInterServerAuthenticationInfoResource(resourceName string, resourceModel interServerAuthenticationInfoTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_inter_server_authentication_info" "%[1]s" {
  type                          = "certificate"
  id                            = "%[2]s"
  server_instance_listener_name = "%[3]s"
  server_instance_name          = "%[4]s"
  purpose                       = %[5]s
}`, resourceName,
		resourceModel.id,
		resourceModel.serverInstanceListenerName,
		resourceModel.serverInstanceName,
		acctest.StringSliceToTerraformString(resourceModel.purpose))
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedInterServerAuthenticationInfoAttributes(config interServerAuthenticationInfoTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.InterServerAuthenticationInfoApi.GetInterServerAuthenticationInfo(ctx, config.id, config.serverInstanceListenerName, config.serverInstanceName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Inter Server Authentication Info"
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "purpose",
			config.purpose, client.StringSliceEnuminterServerAuthenticationInfoPurposeProp(response.CertificateInterServerAuthenticationInfoResponse.Purpose))
		if err != nil {
			return err
		}
		return nil
	}
}
