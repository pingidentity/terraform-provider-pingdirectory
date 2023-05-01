package attributesyntax_test

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

//const testIdNameAndOptionalUidAttributeSyntax = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type nameAndOptionalUidAttributeSyntaxTestModel struct {
	id                      string
	enable_compaction       bool
	require_binary_transfer bool
}

func TestAccNameAndOptionalUidAttributeSyntax(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := nameAndOptionalUidAttributeSyntaxTestModel{
		id:                      "Name and Optional UID",
		enable_compaction:       true,
		require_binary_transfer: true,
	}
	updatedResourceModel := nameAndOptionalUidAttributeSyntaxTestModel{
		id:                      "Name and Optional UID",
		enable_compaction:       false,
		require_binary_transfer: false,
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
				Config: testAccNameAndOptionalUidAttributeSyntaxResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedNameAndOptionalUidAttributeSyntaxAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccNameAndOptionalUidAttributeSyntaxResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedNameAndOptionalUidAttributeSyntaxAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccNameAndOptionalUidAttributeSyntaxResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_default_name_and_optional_uid_attribute_syntax." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccNameAndOptionalUidAttributeSyntaxResource(resourceName string, resourceModel nameAndOptionalUidAttributeSyntaxTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_name_and_optional_uid_attribute_syntax" "%[1]s" {
  id                      = "%[2]s"
  enable_compaction       = %[3]t
  require_binary_transfer = %[4]t
}`, resourceName,
		resourceModel.id,
		resourceModel.enable_compaction,
		resourceModel.require_binary_transfer)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedNameAndOptionalUidAttributeSyntaxAttributes(config nameAndOptionalUidAttributeSyntaxTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "Name and Optional UID"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AttributeSyntaxApi.GetAttributeSyntax(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enable-compaction",
			config.enable_compaction, *response.NameAndOptionalUidAttributeSyntaxResponse.EnableCompaction)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "require-binary-transfer",
			config.require_binary_transfer, *response.NameAndOptionalUidAttributeSyntaxResponse.RequireBinaryTransfer)
		if err != nil {
			return err
		}
		return nil
	}
}
