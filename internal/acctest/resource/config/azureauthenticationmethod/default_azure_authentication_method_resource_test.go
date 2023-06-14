package azureauthenticationmethod_test

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

const testIdAzureAuthenticationMethod = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type azureAuthenticationMethodTestModel struct {
	id string
}

func TestAccAzureAuthenticationMethod(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := azureAuthenticationMethodTestModel{
		id: testIdAzureAuthenticationMethod,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckAzureAuthenticationMethodDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccAzureAuthenticationMethodResource(resourceName, initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccAzureAuthenticationMethodResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_azure_authentication_method." + resourceName,
				ImportStateId:     initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccAzureAuthenticationMethodResource(resourceName string, resourceModel azureAuthenticationMethodTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_azure_authentication_method" "%[1]s" {
  type = "default"
	 id = "%[2]s"
}`, resourceName,
		resourceModel.id)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedAzureAuthenticationMethodAttributes(config azureAuthenticationMethodTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		_, _, err := testClient.AzureAuthenticationMethodApi.GetAzureAuthenticationMethod(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckAzureAuthenticationMethodDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.AzureAuthenticationMethodApi.GetAzureAuthenticationMethod(ctx, testIdAzureAuthenticationMethod).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Azure Authentication Method", testIdAzureAuthenticationMethod)
	}
	return nil
}
