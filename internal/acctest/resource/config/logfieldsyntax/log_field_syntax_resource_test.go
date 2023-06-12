package logfieldsyntax_test

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

const testResource = "Generalized Time"

// Attributes to test with. Add optional properties to test here if desired.
type genericLogFieldSyntaxTestModel struct {
	id               string
	default_behavior string
}

func TestAccGenericLogFieldSyntax(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := genericLogFieldSyntaxTestModel{
		id:               testResource,
		default_behavior: "tokenize-entire-value",
	}
	// set field back to default value of 'preserve'
	updatedResourceModel := genericLogFieldSyntaxTestModel{
		id:               testResource,
		default_behavior: "preserve",
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
				Config: testAccGenericLogFieldSyntaxResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedGenericLogFieldSyntaxAttributes(initialResourceModel),
			},

			{
				// Test updating some fields
				Config: testAccGenericLogFieldSyntaxResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedGenericLogFieldSyntaxAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccGenericLogFieldSyntaxResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_default_log_field_syntax." + resourceName,
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

func testAccGenericLogFieldSyntaxResource(resourceName string, resourceModel genericLogFieldSyntaxTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_log_field_syntax" "%[1]s" {
	type = "generic"
  id               = "%[2]s"
  default_behavior = "%[3]s"
}`, resourceName,
		resourceModel.id,
		resourceModel.default_behavior)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedGenericLogFieldSyntaxAttributes(config genericLogFieldSyntaxTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := testResource
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LogFieldSyntaxApi.GetLogFieldSyntax(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.default_behavior, "default-behavior",
			config.default_behavior, response.GenericLogFieldSyntaxResponse.DefaultBehavior.String())
		if err != nil {
			return err
		}
		return nil
	}
}