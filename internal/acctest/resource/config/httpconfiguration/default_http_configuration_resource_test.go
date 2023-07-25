package httpconfiguration_test

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

// Attributes to test with. Add optional properties to test here if desired.
type httpConfigurationTestModel struct {
	include_stack_traces_in_error_pages bool
}

func TestAccHttpConfiguration(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := httpConfigurationTestModel{
		include_stack_traces_in_error_pages: true,
	}
	updatedResourceModel := httpConfigurationTestModel{
		include_stack_traces_in_error_pages: false,
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
				Config: testAccHttpConfigurationResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					// Check the default value
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_default_http_configuration.%s", resourceName), "include_stack_traces_in_error_pages", "true"),
				),
			},
			{
				// Test updating some fields
				Config: testAccHttpConfigurationResource(resourceName, updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedHttpConfigurationAttributes(updatedResourceModel),
				),
			},
			{
				// Test importing the resource
				Config:            testAccHttpConfigurationResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_default_http_configuration." + resourceName,
				ImportStateId:     resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccHttpConfigurationResource(resourceName string, resourceModel httpConfigurationTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_http_configuration" "%[1]s" {
  include_stack_traces_in_error_pages = %[2]t
}`, resourceName,
		resourceModel.include_stack_traces_in_error_pages)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedHttpConfigurationAttributes(config httpConfigurationTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "http configuration"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.HttpConfigurationApi.GetHttpConfiguration(ctx).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, nil, "include-stack-traces-in-error-pages", config.include_stack_traces_in_error_pages, *response.IncludeStackTracesInErrorPages)
		if err != nil {
			return err
		}
		return nil
	}
}
