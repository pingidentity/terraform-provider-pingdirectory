// Copyright © 2025 Ping Identity Corporation

package connectionhandler_test

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

// Some attributes to test with
const resourceName = "http"
const configId = "example"

type testModel struct {
	id                   string
	listenPort           int64
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
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckHttpConnectionHandlerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				Config: testAccHttpConnectionHandler(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedHttpConnectionHandlerAttributes(initialResourceModel),
					// Check some computed attributes are set as expected (PingDirectory defaults)
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_connection_handler.%s", resourceName), "use_ssl", "false"),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_connection_handler.%s", resourceName), "listen_port", strconv.FormatInt(initialResourceModel.listenPort, 10)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_connection_handler.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_connection_handlers.list", "objects.0.id"),
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
				Config:            testAccHttpConnectionHandler(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_connection_handler." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.ConnectionHandlerAPI.DeleteConnectionHandler(ctx, updatedResourceModel.id).Execute()
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

func testAccHttpConnectionHandler(resourceName string, resourceModel testModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_connection_handler" "%[1]s" {
  type                   = "http"
  name                   = "%[2]s"
  listen_port            = %[3]d
  enabled                = %[4]t
  http_servlet_extension = %[5]s
}

data "pingdirectory_connection_handler" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_connection_handler.%[1]s
  ]
}

data "pingdirectory_connection_handlers" "list" {
  depends_on = [
    pingdirectory_connection_handler.%[1]s
  ]
}`, resourceName, resourceModel.id, resourceModel.listenPort, resourceModel.enabled,
		acctest.StringSliceToTerraformString(resourceModel.httpServletExtension))
}

// Test that any handlers created by the test are destroyed
func testAccCheckHttpConnectionHandlerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ConnectionHandlerAPI.GetConnectionHandler(ctx, configId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("http connection handler", configId)
	}
	return nil
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedHttpConnectionHandlerAttributes(config testModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "http connection handler"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ConnectionHandlerAPI.GetConnectionHandler(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "listen_port",
			config.listenPort, response.HttpConnectionHandlerResponse.ListenPort)
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
