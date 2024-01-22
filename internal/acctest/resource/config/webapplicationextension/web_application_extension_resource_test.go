package webapplicationextension_test

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

const testIdGenericWebApplicationExtension = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type genericWebApplicationExtensionTestModel struct {
	id                    string
	baseContextPath       string
	documentRootDirectory string
}

func TestAccGenericWebApplicationExtension(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := genericWebApplicationExtensionTestModel{
		id:                    testIdGenericWebApplicationExtension,
		baseContextPath:       "/asdf",
		documentRootDirectory: "/tmp",
	}
	updatedResourceModel := genericWebApplicationExtensionTestModel{
		id:                    testIdGenericWebApplicationExtension,
		baseContextPath:       "/jkl",
		documentRootDirectory: "/opt",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckGenericWebApplicationExtensionDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccGenericWebApplicationExtensionResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedGenericWebApplicationExtensionAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_web_application_extension.%s", resourceName), "base_context_path", initialResourceModel.baseContextPath),
					resource.TestCheckResourceAttrSet("data.pingdirectory_web_application_extensions.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccGenericWebApplicationExtensionResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedGenericWebApplicationExtensionAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccGenericWebApplicationExtensionResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_web_application_extension." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				// Required actions only get returned on the specific request where an attriute is changed
				ImportStateVerifyIgnore: []string{
					"required_actions",
				},
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.WebApplicationExtensionAPI.DeleteWebApplicationExtension(ctx, updatedResourceModel.id).Execute()
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

func testAccGenericWebApplicationExtensionResource(resourceName string, resourceModel genericWebApplicationExtensionTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_web_application_extension" "%[1]s" {
  type                    = "generic"
  name                    = "%[2]s"
  base_context_path       = "%[3]s"
  document_root_directory = "%[4]s"
}

data "pingdirectory_web_application_extension" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_web_application_extension.%[1]s
  ]
}

data "pingdirectory_web_application_extensions" "list" {
  depends_on = [
    pingdirectory_web_application_extension.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.baseContextPath,
		resourceModel.documentRootDirectory)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedGenericWebApplicationExtensionAttributes(config genericWebApplicationExtensionTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.WebApplicationExtensionAPI.GetWebApplicationExtension(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Generic Web Application Extension"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "base-context-path",
			config.baseContextPath, response.GenericWebApplicationExtensionResponse.BaseContextPath)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "document-root-directory",
			config.documentRootDirectory, response.GenericWebApplicationExtensionResponse.DocumentRootDirectory)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckGenericWebApplicationExtensionDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.WebApplicationExtensionAPI.GetWebApplicationExtension(ctx, testIdGenericWebApplicationExtension).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Generic Web Application Extension", testIdGenericWebApplicationExtension)
	}
	return nil
}
