package connectionhandler_test

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

// Some attributes to test with
const resourceName = "http"
const configId = "example"

type testModel struct {
	id                   string
	listenPort           int
	enabled              bool
	httpServletExtension []string
}

func TestAccHttpConnectionHandler(t *testing.T) {

	initialResourceModel := testModel{
		id:                   configId,
		listenPort:           2443,
		enabled:              true,
		httpServletExtension: []string{"Available or Degraded State", "Available State"},
	}
	updatedResourceModel := testModel{
		id:                   configId,
		listenPort:           2444,
		enabled:              false,
		httpServletExtension: []string{"Available or Degraded State"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckHttpConnectionHandlerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				Config: testAccHttpConnectionHandler(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedHttpConnectionHandlerAttributes(initialResourceModel),
					// Check some computed attributes are set as expected (PingDirectory defaults)
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_http_connection_handler.%s", resourceName), "use_ssl", "false"),
				),
			},
			{
				// Test updating some fields
				Config: testAccHttpConnectionHandler(resourceName, updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedHttpConnectionHandlerAttributes(updatedResourceModel),
				),
			},
			{
				// Test importing the resource
				Config:                  testAccHttpConnectionHandler(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_http_connection_handler." + resourceName,
				ImportStateId:           updatedResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccHttpConnectionHandler(resourceName string, resourceModel testModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_http_connection_handler" "%[1]s" {
	id = "%[2]s"
	listen_port = %[3]d
	enabled = %[4]t
	http_servlet_extension = %[5]s
}`, resourceName, resourceModel.id, resourceModel.listenPort, resourceModel.enabled,
		acctest.StringSliceToTerraformString(resourceModel.httpServletExtension))
}

// Test that any handlers created by the test are destroyed
func testAccCheckHttpConnectionHandlerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ConnectionHandlerApi.GetConnectionHandler(ctx, configId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("trust manager provider", configId)
	}
	return nil
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedHttpConnectionHandlerAttributes(config testModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "http connection handler"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ConnectionHandlerApi.GetConnectionHandler(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "listen_port",
			int64(config.listenPort), int64(response.HttpConnectionHandlerResponse.ListenPort))
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.HttpConnectionHandlerResponse.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "http_servlet_extension",
			config.httpServletExtension, response.HttpConnectionHandlerResponse.HttpServletExtension)
		if err != nil {
			return err
		}
		return nil
	}
}
