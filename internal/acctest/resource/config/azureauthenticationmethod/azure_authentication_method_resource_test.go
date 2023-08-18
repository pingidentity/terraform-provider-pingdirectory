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

const testIdDefaultAzureAuthenticationMethod = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type defaultAzureAuthenticationMethodTestModel struct {
	id string
}

func TestAccDefaultAzureAuthenticationMethod(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := defaultAzureAuthenticationMethodTestModel{
		id: testIdDefaultAzureAuthenticationMethod,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckDefaultAzureAuthenticationMethodDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDefaultAzureAuthenticationMethodResource(resourceName, initialResourceModel),
				Check:  resource.TestCheckResourceAttrSet("data.pingdirectory_azure_authentication_methods.list", "objects.0.id"),
			},
			{
				// Test importing the resource
				Config:            testAccDefaultAzureAuthenticationMethodResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_azure_authentication_method." + resourceName,
				ImportStateId:     initialResourceModel.id,
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
					_, err := testClient.AzureAuthenticationMethodApi.DeleteAzureAuthenticationMethod(ctx, initialResourceModel.id).Execute()
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

func testAccDefaultAzureAuthenticationMethodResource(resourceName string, resourceModel defaultAzureAuthenticationMethodTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_azure_authentication_method" "%[1]s" {
  type = "default"
  name = "%[2]s"
}

data "pingdirectory_azure_authentication_method" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_azure_authentication_method.%[1]s
  ]
}

data "pingdirectory_azure_authentication_methods" "list" {
  depends_on = [
    pingdirectory_azure_authentication_method.%[1]s
  ]
}`, resourceName,
		resourceModel.id)
}

// Test that any objects created by the test are destroyed
func testAccCheckDefaultAzureAuthenticationMethodDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.AzureAuthenticationMethodApi.GetAzureAuthenticationMethod(ctx, testIdDefaultAzureAuthenticationMethod).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Default Azure Authentication Method", testIdDefaultAzureAuthenticationMethod)
	}
	return nil
}
