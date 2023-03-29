package httpservletextension_test

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

const testIdQuickstartHttpServletExtension = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type quickstartHttpServletExtensionTestModel struct {
	id          string
	description string
}

func TestAccQuickstartHttpServletExtension(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := quickstartHttpServletExtensionTestModel{
		id:          testIdQuickstartHttpServletExtension,
		description: "example description",
	}
	updatedResourceModel := quickstartHttpServletExtensionTestModel{
		id:          testIdQuickstartHttpServletExtension,
		description: "example updated description",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckQuickstartHttpServletExtensionDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccQuickstartHttpServletExtensionResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedQuickstartHttpServletExtensionAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccQuickstartHttpServletExtensionResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedQuickstartHttpServletExtensionAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccQuickstartHttpServletExtensionResource(resourceName, initialResourceModel),
				ResourceName:            "pingdirectory_quickstart_http_servlet_extension." + resourceName,
				ImportStateId:           initialResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccQuickstartHttpServletExtensionResource(resourceName string, resourceModel quickstartHttpServletExtensionTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_quickstart_http_servlet_extension" "%[1]s" {
  id          = "%[2]s"
  description = "%[3]s"
}`, resourceName, resourceModel.id, resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedQuickstartHttpServletExtensionAttributes(config quickstartHttpServletExtensionTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "http servlet extension"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.HttpServletExtensionApi.GetHttpServletExtension(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that description matches expected
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "description",
			config.description, response.QuickstartHttpServletExtensionResponse.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckQuickstartHttpServletExtensionDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.HttpServletExtensionApi.GetHttpServletExtension(ctx, testIdQuickstartHttpServletExtension).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Quickstart Http Servlet Extension", testIdQuickstartHttpServletExtension)
	}
	return nil
}
