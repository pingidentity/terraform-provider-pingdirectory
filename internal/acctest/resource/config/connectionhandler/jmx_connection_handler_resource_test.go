package connectionhandler_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdJmxConnectionHandler = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type jmxConnectionHandlerTestModel struct {
	id         string
	listenPort int64
	enabled    bool
}

func TestAccJmxConnectionHandler(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := jmxConnectionHandlerTestModel{
		id:         testIdJmxConnectionHandler,
		listenPort: 1111,
		enabled:    true,
	}
	updatedResourceModel := jmxConnectionHandlerTestModel{
		id:         testIdJmxConnectionHandler,
		listenPort: 1112,
		enabled:    false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckJmxConnectionHandlerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccJmxConnectionHandlerResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedJmxConnectionHandlerAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccJmxConnectionHandlerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedJmxConnectionHandlerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccJmxConnectionHandlerResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_jmx_connection_handler." + resourceName,
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

func testAccJmxConnectionHandlerResource(resourceName string, resourceModel jmxConnectionHandlerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_jmx_connection_handler" "%[1]s" {
  id          = "%[2]s"
  listen_port = %[3]d
  enabled     = %[4]t
}`, resourceName,
		resourceModel.id,
		resourceModel.listenPort,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedJmxConnectionHandlerAttributes(config jmxConnectionHandlerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ConnectionHandlerApi.GetConnectionHandler(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Jmx Connection Handler"
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "listen-port",
			config.listenPort, int64(response.JmxConnectionHandlerResponse.ListenPort))
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.JmxConnectionHandlerResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckJmxConnectionHandlerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ConnectionHandlerApi.GetConnectionHandler(ctx, testIdJmxConnectionHandler).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Jmx Connection Handler", testIdJmxConnectionHandler)
	}
	return nil
}
