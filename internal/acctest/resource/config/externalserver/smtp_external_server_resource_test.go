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
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckSmtpExternalServerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSmtpExternalServerResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSmtpExternalServerAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccSmtpExternalServerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSmtpExternalServerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccSmtpExternalServerResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_external_server." + resourceName,
				ImportStateId:           updatedResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccSmtpExternalServerResource(resourceName string, resourceModel smtpExternalServerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_external_server" "%[1]s" {
	type = "smtp"
  id               = "%[2]s"
  server_host_name = "%[3]s"
  server_port      = %[4]d
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
		response, _, err := testClient.ExternalServerApi.GetExternalServer(ctx, config.id).Execute()
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
	_, _, err := testClient.ExternalServerApi.GetExternalServer(ctx, testIdSmtpExternalServer).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Smtp External Server", testIdSmtpExternalServer)
	}
	return nil
}
