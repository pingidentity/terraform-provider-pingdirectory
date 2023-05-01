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

const testIdNameAndOptionalUidAttributeSyntax = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type nameAndOptionalUidAttributeSyntaxTestModel struct {
	   id         string
	   enabled         bool
}

func TestAccNameAndOptionalUidAttributeSyntax(t *testing.T) {
	   resourceName := "myresource"
	   initialResourceModel := nameAndOptionalUidAttributeSyntaxTestModel{
        id: testIdNameAndOptionalUidAttributeSyntax,
	   enabled: //TODO set appropriate value,
    }
	   updatedResourceModel := nameAndOptionalUidAttributeSyntaxTestModel{
        id: testIdNameAndOptionalUidAttributeSyntax,
	   enabled: //TODO set appropriate value,
    }

	   resource.Test(t, resource.TestCase{
	   	   PreCheck: func() { acctest.ConfigurationPreCheck(t) },
	   	   ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
	   	   	   "pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
	   	   },
	   	   CheckDestroy: testAccCheckNameAndOptionalUidAttributeSyntaxDestroy,
	   	   Steps: []resource.TestStep{
	   	   	   {
	   	   	   	   // Test basic resource.
	   	   	   	   // Add checks for computed properties here if desired.
	   	   	   	   Config: testAccNameAndOptionalUidAttributeSyntaxResource(resourceName, initialResourceModel),
	   	   	   	   Check: testAccCheckExpectedNameAndOptionalUidAttributeSyntaxAttributes(initialResourceModel),
	   	   	   },
	   	   	   {
	   	   	   	   // Test updating some fields
	   	   	   	   Config: testAccNameAndOptionalUidAttributeSyntaxResource(resourceName, updatedResourceModel),
	   	   	   	   Check: testAccCheckExpectedNameAndOptionalUidAttributeSyntaxAttributes(updatedResourceModel),
	   	   	   },
	   	   	   {
	   	   	   	   // Test importing the resource
	   	   	   	   Config:                  testAccNameAndOptionalUidAttributeSyntaxResource(resourceName, updatedResourceModel),
	   	   	   	   ResourceName:            "pingdirectory_default_name_and_optional_uid_attribute_syntax." + resourceName,
	   	   	   	   ImportStateId:           updatedResourceModel.id,
	   	   	   	   ImportState:             true,
	   	   	   	   ImportStateVerify:       true,
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
	 id = "%[2]s"
	 enabled = %[3]t
}`, resourceName,
    resourceModel.id,
    resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedNameAndOptionalUidAttributeSyntaxAttributes(config nameAndOptionalUidAttributeSyntaxTestModel) resource.TestCheckFunc {
	   return func(s *terraform.State) error {
	   	   testClient := acctest.TestClient()
	   	   ctx := acctest.TestBasicAuthContext()
	   	   response, _, err := testClient.AttributeSyntaxApi.GetAttributeSyntax(ctx, config.id).Execute()
	   	   if err != nil {
	   	   	   return err
	   	   }
	   	   // Verify that attributes have expected values
	   	   resourceType := "Name And Optional Uid Attribute Syntax"
	   	   err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
	   	   	   config.enabled, response.NameAndOptionalUidAttributeSyntaxResponse.Enabled)
	   	   if err != nil {
	   	   	   return err
	   	   }
	   	   return nil
	   }
}

