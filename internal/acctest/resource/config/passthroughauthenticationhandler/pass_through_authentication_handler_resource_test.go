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
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckLdapPassThroughAuthenticationHandlerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				Config: testAccLdapPassThroughAuthenticationHandlerResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedLdapPassThroughAuthenticationHandlerAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_pass_through_authentication_handler.%s", resourceName), "description", initialResourceModel.description),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_pass_through_authentication_handler.%s", resourceName), "server.*", initialResourceModel.server[0]),
					resource.TestCheckResourceAttrSet("data.pingdirectory_pass_through_authentication_handlers.list", "objects.0.id"),
				),
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
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.PassThroughAuthenticationHandlerApi.DeletePassThroughAuthenticationHandler(ctx, updatedResourceModel.id).Execute()
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

func testAccLdapPassThroughAuthenticationHandlerResource(resourceName string, resourceModel ldapPassThroughAuthenticationHandlerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_external_server" "%[4]s" {
  type                  = "ldap"
  name                  = "%[4]s"
  server_host_name      = "localhost"
  authentication_method = "none"
}

resource "pingdirectory_pass_through_authentication_handler" "%[1]s" {
  type        = "ldap"
  name        = "%[2]s"
  description = "%[3]s"
  server      = [pingdirectory_external_server.%[4]s.id]
}

data "pingdirectory_pass_through_authentication_handler" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_pass_through_authentication_handler.%[1]s
  ]
}

data "pingdirectory_pass_through_authentication_handlers" "list" {
  depends_on = [
    pingdirectory_pass_through_authentication_handler.%[1]s
  ]
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
