// Copyright © 2025 Ping Identity Corporation

package externalserver_test

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

const testIdSmtpExternalServer = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type smtpExternalServerTestModel struct {
	id             string
	serverHostName string
	serverPort     int64
}

func TestAccSmtpExternalServer(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := smtpExternalServerTestModel{
		id:             testIdSmtpExternalServer,
		serverHostName: "initial.mailserver.com",
		serverPort:     25,
	}

	updatedResourceModel := smtpExternalServerTestModel{
		id:             testIdSmtpExternalServer,
		serverHostName: "modified.mailserver.com",
		serverPort:     6225,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckSmtpExternalServerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSmtpExternalServerResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedSmtpExternalServerAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_external_server.%s", resourceName), "server_host_name", initialResourceModel.serverHostName),
					resource.TestCheckResourceAttrSet("data.pingdirectory_external_servers.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccSmtpExternalServerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSmtpExternalServerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSmtpExternalServerResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_external_server." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.ExternalServerAPI.DeleteExternalServer(ctx, updatedResourceModel.id).Execute()
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

func testAccSmtpExternalServerResource(resourceName string, resourceModel smtpExternalServerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_external_server" "%[1]s" {
  type             = "smtp"
  name             = "%[2]s"
  server_host_name = "%[3]s"
  server_port      = %[4]d
}

data "pingdirectory_external_server" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_external_server.%[1]s
  ]
}

data "pingdirectory_external_servers" "list" {
  depends_on = [
    pingdirectory_external_server.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.serverHostName,
		resourceModel.serverPort)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSmtpExternalServerAttributes(config smtpExternalServerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ExternalServerAPI.GetExternalServer(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Smtp External Server"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "serverHostName",
			config.serverHostName, response.SmtpExternalServerResponse.ServerHostName)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "serverPort",
			config.serverPort, *response.SmtpExternalServerResponse.ServerPort)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSmtpExternalServerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ExternalServerAPI.GetExternalServer(ctx, testIdSmtpExternalServer).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Smtp External Server", testIdSmtpExternalServer)
	}
	return nil
}
