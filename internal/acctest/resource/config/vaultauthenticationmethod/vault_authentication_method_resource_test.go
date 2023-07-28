package vaultauthenticationmethod_test

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

const testIdVaultAuthenticationMethod = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type vaultAuthenticationMethodTestModel struct {
	id               string
	vaultAccessToken string
	description      string
}

func TestAccVaultAuthenticationMethod(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := vaultAuthenticationMethodTestModel{
		id:               testIdVaultAuthenticationMethod,
		vaultAccessToken: "myfirsttoken1234",
		description:      "mydescription",
	}
	updatedResourceModel := vaultAuthenticationMethodTestModel{
		id:               testIdVaultAuthenticationMethod,
		vaultAccessToken: "mysecondtoken5678",
		description:      "anotherdescription",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckVaultAuthenticationMethodDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccVaultAuthenticationMethodResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedVaultAuthenticationMethodAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_vault_authentication_method.%s", resourceName), "description", initialResourceModel.description),
				),
			},
			{
				// Test updating some fields
				Config: testAccVaultAuthenticationMethodResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedVaultAuthenticationMethodAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccVaultAuthenticationMethodResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_vault_authentication_method." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
					"vault_access_token",
				},
			},
		},
	})
}

func testAccVaultAuthenticationMethodResource(resourceName string, resourceModel vaultAuthenticationMethodTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_vault_authentication_method" "%[1]s" {
  type               = "static-token"
  id                 = "%[2]s"
  vault_access_token = "%[3]s"
  description        = "%[4]s"
}

data "pingdirectory_vault_authentication_method" "%[1]s" {
  id = "%[2]s"
  depends_on = [
    pingdirectory_vault_authentication_method.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.vaultAccessToken,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedVaultAuthenticationMethodAttributes(config vaultAuthenticationMethodTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.VaultAuthenticationMethodApi.GetVaultAuthenticationMethod(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		// Note that we can't check the obscured vault-access-token attribute because PD won't return it in plaintext
		resourceType := "Vault Authentication Method"
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "description",
			config.description, response.StaticTokenVaultAuthenticationMethodResponse.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckVaultAuthenticationMethodDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.VaultAuthenticationMethodApi.GetVaultAuthenticationMethod(ctx, testIdVaultAuthenticationMethod).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Vault Authentication Method", testIdVaultAuthenticationMethod)
	}
	return nil
}
