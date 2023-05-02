package changesubscriptionhandler_test

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

const testIdGroovyScriptedChangeSubscriptionHandler = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type groovyScriptedChangeSubscriptionHandlerTestModel struct {
	id          string
	scriptClass string
	enabled     bool
}

func TestAccGroovyScriptedChangeSubscriptionHandler(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := groovyScriptedChangeSubscriptionHandlerTestModel{
		id:          testIdGroovyScriptedChangeSubscriptionHandler,
		scriptClass: "com.unboundid.directory.sdk.ds.api.ChangeSubscriptionHandler",
		enabled:     false,
	}
	updatedResourceModel := groovyScriptedChangeSubscriptionHandlerTestModel{
		id:          testIdGroovyScriptedChangeSubscriptionHandler,
		scriptClass: "com.unboundid.directory.sdk.ds.api.ChangeSubscriptionHandler",
		enabled:     false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckGroovyScriptedChangeSubscriptionHandlerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccGroovyScriptedChangeSubscriptionHandlerResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedGroovyScriptedChangeSubscriptionHandlerAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccGroovyScriptedChangeSubscriptionHandlerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedGroovyScriptedChangeSubscriptionHandlerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccGroovyScriptedChangeSubscriptionHandlerResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_groovy_scripted_change_subscription_handler." + resourceName,
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

func testAccGroovyScriptedChangeSubscriptionHandlerResource(resourceName string, resourceModel groovyScriptedChangeSubscriptionHandlerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_groovy_scripted_change_subscription_handler" "%[1]s" {
  id           = "%[2]s"
  script_class = "%[3]s"
  enabled      = %[4]t
}`, resourceName,
		resourceModel.id,
		resourceModel.scriptClass,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedGroovyScriptedChangeSubscriptionHandlerAttributes(config groovyScriptedChangeSubscriptionHandlerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ChangeSubscriptionHandlerApi.GetChangeSubscriptionHandler(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Groovy Scripted Change Subscription Handler"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "script-class",
			config.scriptClass, response.GroovyScriptedChangeSubscriptionHandlerResponse.ScriptClass)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.GroovyScriptedChangeSubscriptionHandlerResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckGroovyScriptedChangeSubscriptionHandlerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ChangeSubscriptionHandlerApi.GetChangeSubscriptionHandler(ctx, testIdGroovyScriptedChangeSubscriptionHandler).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Groovy Scripted Change Subscription Handler", testIdGroovyScriptedChangeSubscriptionHandler)
	}
	return nil
}
