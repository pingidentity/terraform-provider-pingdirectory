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
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckVelocityContextProviderDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccVelocityContextProviderResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedVelocityContextProviderAttributes(initialResourceModel),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_velocity_context_provider.%s", resourceName), "included_view.*", initialResourceModel.includedView[0]),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_velocity_context_provider.%s", resourceName), "included_view.*", initialResourceModel.includedView[1]),
					resource.TestCheckResourceAttrSet("data.pingdirectory_velocity_context_providers.list", "objects.0.id"),
				),
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
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.VelocityContextProviderApi.DeleteVelocityContextProvider(ctx, updatedResourceModel.id, updatedResourceModel.httpServletExtensionName).Execute()
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

func testAccVelocityContextProviderResource(resourceName string, resourceModel velocityContextProviderTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_velocity_context_provider" "%[1]s" {
  type                        = "velocity-tools"
  name                        = "%[2]s"
  http_servlet_extension_name = "%[3]s"
  included_view               = %[4]s
}

data "pingdirectory_velocity_context_provider" "%[1]s" {
  name                        = "%[2]s"
  http_servlet_extension_name = "%[3]s"
  depends_on = [
    pingdirectory_velocity_context_provider.%[1]s
  ]
}

data "pingdirectory_velocity_context_providers" "list" {
  http_servlet_extension_name = "%[3]s"
  depends_on = [
    pingdirectory_velocity_context_provider.%[1]s
  ]
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
