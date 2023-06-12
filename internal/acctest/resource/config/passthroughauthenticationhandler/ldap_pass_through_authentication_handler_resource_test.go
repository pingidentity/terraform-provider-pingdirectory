package passthroughauthenticationhandler_test

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

const testIdLdapPassThroughAuthenticationHandler = "MyId"
const ldapExternalServerId = "MyLdapExternalServer"

// Attributes to test with. Add optional properties to test here if desired.
type ldapPassThroughAuthenticationHandlerTestModel struct {
	id          string
	description string
	server      []string
}

func TestAccLdapPassThroughAuthenticationHandler(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := ldapPassThroughAuthenticationHandlerTestModel{
		id:          testIdLdapPassThroughAuthenticationHandler,
		description: "initialDesc",
		server:      []string{ldapExternalServerId},
	}
	updatedResourceModel := ldapPassThroughAuthenticationHandlerTestModel{
		id:          testIdLdapPassThroughAuthenticationHandler,
		description: "updatedDesc",
		server:      []string{ldapExternalServerId},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckLdapPassThroughAuthenticationHandlerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				Config: testAccLdapPassThroughAuthenticationHandlerResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedLdapPassThroughAuthenticationHandlerAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccLdapPassThroughAuthenticationHandlerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedLdapPassThroughAuthenticationHandlerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccLdapPassThroughAuthenticationHandlerResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_pass_through_authentication_handler." + resourceName,
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

func testAccLdapPassThroughAuthenticationHandlerResource(resourceName string, resourceModel ldapPassThroughAuthenticationHandlerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_external_server" "%[4]s" {
	type = "ldap"
  id                    = "%[4]s"
  server_host_name      = "localhost"
  authentication_method = "none"
}

resource "pingdirectory_pass_through_authentication_handler" "%[1]s" {
	type = "ldap"
  id          = "%[2]s"
  description = "%[3]s"
  server      = [pingdirectory_external_server.%[4]s.id]
}`, resourceName,
		resourceModel.id,
		resourceModel.description,
		resourceModel.server[0])
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedLdapPassThroughAuthenticationHandlerAttributes(config ldapPassThroughAuthenticationHandlerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.PassThroughAuthenticationHandlerApi.GetPassThroughAuthenticationHandler(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Ldap Pass Through Authentication Handler"
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "server",
			config.server, response.LdapPassThroughAuthenticationHandlerResponse.Server)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "description",
			config.description, response.LdapPassThroughAuthenticationHandlerResponse.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckLdapPassThroughAuthenticationHandlerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.PassThroughAuthenticationHandlerApi.GetPassThroughAuthenticationHandler(ctx, testIdLdapPassThroughAuthenticationHandler).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Ldap Pass Through Authentication Handler", testIdLdapPassThroughAuthenticationHandler)
	}
	return nil
}
