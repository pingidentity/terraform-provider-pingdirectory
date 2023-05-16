package oauthtokenhandler_test

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

const testIdGroovyScriptedOauthTokenHandler = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type groovyScriptedOauthTokenHandlerTestModel struct {
	id          string
	description string
	scriptClass string
}

func TestAccGroovyScriptedOauthTokenHandler(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := groovyScriptedOauthTokenHandlerTestModel{
		id:          testIdGroovyScriptedOauthTokenHandler,
		description: "initial resource model description",
		scriptClass: "com.example",
	}
	updatedResourceModel := groovyScriptedOauthTokenHandlerTestModel{
		id:          testIdGroovyScriptedOauthTokenHandler,
		description: "updated resource model description",
		scriptClass: "com.company",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckGroovyScriptedOauthTokenHandlerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccGroovyScriptedOauthTokenHandlerResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedGroovyScriptedOauthTokenHandlerAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccGroovyScriptedOauthTokenHandlerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedGroovyScriptedOauthTokenHandlerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccGroovyScriptedOauthTokenHandlerResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_groovy_scripted_oauth_token_handler." + resourceName,
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

func testAccGroovyScriptedOauthTokenHandlerResource(resourceName string, resourceModel groovyScriptedOauthTokenHandlerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_groovy_scripted_oauth_token_handler" "%[1]s" {
  id           = "%[2]s"
  description  = "%[3]s"
  script_class = "%[4]s"
}`, resourceName,
		resourceModel.id,
		resourceModel.description,
		resourceModel.scriptClass)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedGroovyScriptedOauthTokenHandlerAttributes(config groovyScriptedOauthTokenHandlerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthTokenHandlerApi.GetOauthTokenHandler(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Groovy Scripted Oauth Token Handler"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "description",
			config.description, *response.GroovyScriptedOauthTokenHandlerResponse.Description)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "script-class",
			config.scriptClass, response.GroovyScriptedOauthTokenHandlerResponse.ScriptClass)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckGroovyScriptedOauthTokenHandlerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.OauthTokenHandlerApi.GetOauthTokenHandler(ctx, testIdGroovyScriptedOauthTokenHandler).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Groovy Scripted Oauth Token Handler", testIdGroovyScriptedOauthTokenHandler)
	}
	return nil
}