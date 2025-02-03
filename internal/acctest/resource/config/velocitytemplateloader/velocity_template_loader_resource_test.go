// Copyright © 2025 Ping Identity Corporation

package velocitytemplateloader_test

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

const testIdVelocityTemplateLoader = "MyId"
const testHttpServletExtensionName = "Velocity"

// Attributes to test with. Add optional properties to test here if desired.
type velocityTemplateLoaderTestModel struct {
	id                       string
	httpServletExtensionName string
	mimeTypeMatcher          string
}

func TestAccVelocityTemplateLoader(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := velocityTemplateLoaderTestModel{
		id:                       testIdVelocityTemplateLoader,
		httpServletExtensionName: testHttpServletExtensionName,
		mimeTypeMatcher:          "text/html",
	}
	updatedResourceModel := velocityTemplateLoaderTestModel{
		id:                       testIdVelocityTemplateLoader,
		httpServletExtensionName: testHttpServletExtensionName,
		mimeTypeMatcher:          "application/json",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckVelocityTemplateLoaderDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccVelocityTemplateLoaderResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedVelocityTemplateLoaderAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_velocity_template_loader.%s", resourceName), "mime_type_matcher", initialResourceModel.mimeTypeMatcher),
					resource.TestCheckResourceAttrSet("data.pingdirectory_velocity_template_loaders.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccVelocityTemplateLoaderResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedVelocityTemplateLoaderAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccVelocityTemplateLoaderResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_velocity_template_loader." + resourceName,
				ImportStateId:     updatedResourceModel.httpServletExtensionName + "/" + updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.VelocityTemplateLoaderAPI.DeleteVelocityTemplateLoader(ctx, updatedResourceModel.id, updatedResourceModel.httpServletExtensionName).Execute()
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

func testAccVelocityTemplateLoaderResource(resourceName string, resourceModel velocityTemplateLoaderTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_velocity_template_loader" "%[1]s" {
  name                        = "%[2]s"
  http_servlet_extension_name = "%[3]s"
  mime_type_matcher           = "%[4]s"
}

data "pingdirectory_velocity_template_loader" "%[1]s" {
  name                        = "%[2]s"
  http_servlet_extension_name = "%[3]s"
  depends_on = [
    pingdirectory_velocity_template_loader.%[1]s
  ]
}

data "pingdirectory_velocity_template_loaders" "list" {
  http_servlet_extension_name = "%[3]s"
  depends_on = [
    pingdirectory_velocity_template_loader.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.httpServletExtensionName,
		resourceModel.mimeTypeMatcher)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedVelocityTemplateLoaderAttributes(config velocityTemplateLoaderTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.VelocityTemplateLoaderAPI.GetVelocityTemplateLoader(ctx, config.id, config.httpServletExtensionName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Velocity Template Loader"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "mime-type-matcher",
			config.mimeTypeMatcher, response.MimeTypeMatcher)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckVelocityTemplateLoaderDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.VelocityTemplateLoaderAPI.GetVelocityTemplateLoader(ctx, testIdVelocityTemplateLoader, testHttpServletExtensionName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Velocity Template Loader", testIdVelocityTemplateLoader)
	}
	return nil
}
