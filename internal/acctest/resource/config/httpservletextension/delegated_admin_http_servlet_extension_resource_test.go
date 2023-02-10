package httpservletextension_test

import (
	"fmt"
	"testing"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testIdDelegatedAdminHttpServletExtension = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type delegatedAdminHttpServletExtensionTestModel struct {
	id string
}

func TestAccDelegatedAdminHttpServletExtension(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := delegatedAdminHttpServletExtensionTestModel{
		id: testIdDelegatedAdminHttpServletExtension,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckDelegatedAdminHttpServletExtensionDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDelegatedAdminHttpServletExtensionResource(resourceName, initialResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccDelegatedAdminHttpServletExtensionResource(resourceName, initialResourceModel),
				ResourceName:            "pingdirectory_delegated_admin_http_servlet_extension." + resourceName,
				ImportStateId:           initialResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccDelegatedAdminHttpServletExtensionResource(resourceName string, resourceModel delegatedAdminHttpServletExtensionTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_delegated_admin_http_servlet_extension" "%[1]s" {
	 id = "%[2]s"
}`, resourceName, resourceModel.id)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDelegatedAdminHttpServletExtensionAttributes(config delegatedAdminHttpServletExtensionTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		_, _, err := testClient.HttpServletExtensionApi.GetHttpServletExtension(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		return nil
	}
}
