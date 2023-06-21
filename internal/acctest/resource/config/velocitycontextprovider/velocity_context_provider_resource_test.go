package velocitycontextprovider_test

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

const testIdVelocityContextProvider = "MyId"
const testHttpServletExtensionName = "Velocity"

// Attributes to test with. Add optional properties to test here if desired.
type velocityContextProviderTestModel struct {
	id                       string
	httpServletExtensionName string
	includedView             []string
}

func TestAccVelocityContextProvider(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := velocityContextProviderTestModel{
		id:                       testIdVelocityContextProvider,
		httpServletExtensionName: testHttpServletExtensionName,
		includedView:             []string{"path/to/view1", "path/to/view2"},
	}
	updatedResourceModel := velocityContextProviderTestModel{
		id:                       testIdVelocityContextProvider,
		httpServletExtensionName: testHttpServletExtensionName,
		includedView:             []string{"path/to/view3", "path/to/view4"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckVelocityContextProviderDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccVelocityContextProviderResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedVelocityContextProviderAttributes(initialResourceModel),
			},
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccVelocityContextProviderResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedVelocityContextProviderAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccVelocityContextProviderResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_velocity_context_provider." + resourceName,
				ImportStateId:     updatedResourceModel.httpServletExtensionName + "/" + updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccVelocityContextProviderResource(resourceName string, resourceModel velocityContextProviderTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_velocity_context_provider" "%[1]s" {
  type                        = "velocity-tools"
  id                          = "%[2]s"
  http_servlet_extension_name = "%[3]s"
  included_view               = %[4]s
}`, resourceName,
		resourceModel.id,
		resourceModel.httpServletExtensionName,
		acctest.StringSliceToTerraformString(resourceModel.includedView))
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedVelocityContextProviderAttributes(config velocityContextProviderTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.VelocityContextProviderApi.GetVelocityContextProvider(ctx, config.id, config.httpServletExtensionName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Velocity Context Provider"
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "included-view",
			config.includedView, response.VelocityToolsVelocityContextProviderResponse.IncludedView)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckVelocityContextProviderDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.VelocityContextProviderApi.GetVelocityContextProvider(ctx, testIdVelocityContextProvider, testHttpServletExtensionName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Velocity Context Provider", testIdVelocityContextProvider)
	}
	return nil
}
