package ldapsdkdebuglogger_test

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

// Attributes to test with. Add optional properties to test here if desired.
type ldapSdkDebugLoggerTestModel struct {
	description string
}

func TestAccLdapSdkDebugLogger(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := ldapSdkDebugLoggerTestModel{
		description: "initial resource model description",
	}
	updatedResourceModel := ldapSdkDebugLoggerTestModel{
		description: "updated resource model description",
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
				Config: testAccLdapSdkDebugLoggerResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedLdapSdkDebugLoggerAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_ldap_sdk_debug_logger.%s", resourceName), "description", initialResourceModel.description),
				),
			},
			{
				// Test updating fields
				Config: testAccLdapSdkDebugLoggerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedLdapSdkDebugLoggerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:       testAccLdapSdkDebugLoggerResource(resourceName, initialResourceModel),
				ResourceName: "pingdirectory_default_ldap_sdk_debug_logger." + resourceName,
				// The id doesn't matter for singleton config objects
				ImportStateId:     resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccLdapSdkDebugLoggerResource(resourceName string, resourceModel ldapSdkDebugLoggerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_ldap_sdk_debug_logger" "%[1]s" {
  description = "%[2]s"
}

data "pingdirectory_ldap_sdk_debug_logger" "%[1]s" {
  depends_on = [
    pingdirectory_default_ldap_sdk_debug_logger.%[1]s
  ]
}`, resourceName,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedLdapSdkDebugLoggerAttributes(config ldapSdkDebugLoggerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "LDAP SDK Debug Logger"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LdapSdkDebugLoggerApi.GetLdapSdkDebugLogger(ctx).Execute()
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "description", config.description, response.Description)
		if err != nil {
			return err
		}
		return nil
	}
}
