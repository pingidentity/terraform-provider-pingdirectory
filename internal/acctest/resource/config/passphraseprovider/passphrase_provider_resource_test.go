package passphraseprovider_test

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

const testIdPassphraseProvider = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type passphraseProviderTestModel struct {
	id                  string
	environmentVariable string
	enabled             bool
}

func TestAccPassphraseProvider(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := passphraseProviderTestModel{
		id:                  testIdPassphraseProvider,
		environmentVariable: "PING_IDENTITY_DEVOPS_USER",
		enabled:             true,
	}
	updatedResourceModel := passphraseProviderTestModel{
		id:                  testIdPassphraseProvider,
		environmentVariable: "PING_IDENTITY_DEVOPS_KEY",
		enabled:             false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckPassphraseProviderDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccPassphraseProviderResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedPassphraseProviderAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccPassphraseProviderResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedPassphraseProviderAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccPassphraseProviderResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_passphrase_provider." + resourceName,
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

func testAccPassphraseProviderResource(resourceName string, resourceModel passphraseProviderTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_passphrase_provider" "%[1]s" {
  type                 = "environment-variable"
  id                   = "%[2]s"
  environment_variable = "%[3]s"
  enabled              = %[4]t
}`, resourceName,
		resourceModel.id,
		resourceModel.environmentVariable,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedPassphraseProviderAttributes(config passphraseProviderTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.PassphraseProviderApi.GetPassphraseProvider(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Passphrase Provider"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "environment-variable",
			config.environmentVariable, response.EnvironmentVariablePassphraseProviderResponse.EnvironmentVariable)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.EnvironmentVariablePassphraseProviderResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckPassphraseProviderDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.PassphraseProviderApi.GetPassphraseProvider(ctx, testIdPassphraseProvider).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Passphrase Provider", testIdPassphraseProvider)
	}
	return nil
}
