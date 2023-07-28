package accesscontrolhandler_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

// Attributes to test with. Add optional properties to test here if desired.
type dseeCompatAccessControlHandlerTestModel struct {
	enabled bool
}

func TestAccDseeCompatAccessControlHandler(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := dseeCompatAccessControlHandlerTestModel{
		enabled: false,
	}
	updatedResourceModel := dseeCompatAccessControlHandlerTestModel{
		enabled: true,
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
				Config: testAccDseeCompatAccessControlHandlerResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedDseeCompatAccessControlHandlerAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_access_control_handler.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
				),
			},
			{
				// Test updating some fields
				Config: testAccDseeCompatAccessControlHandlerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedDseeCompatAccessControlHandlerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccDseeCompatAccessControlHandlerResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_default_access_control_handler." + resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccDseeCompatAccessControlHandlerResource(resourceName string, resourceModel dseeCompatAccessControlHandlerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_access_control_handler" "%[1]s" {
  enabled = %[2]t
}

data "pingdirectory_access_control_handler" "%[1]s" {
  depends_on = [
    pingdirectory_default_access_control_handler.%[1]s
  ]
}`, resourceName, resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDseeCompatAccessControlHandlerAttributes(config dseeCompatAccessControlHandlerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AccessControlHandlerApi.GetAccessControlHandler(ctx).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Dsee Compat Access Control Handler"
		err = acctest.TestAttributesMatchBool(resourceType, nil, "enabled",
			config.enabled, response.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}
