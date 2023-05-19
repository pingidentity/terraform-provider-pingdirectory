package config_test

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

const testIdScimSubattribute = "MyId"
const testScimAttributeName = "cn"
const testScimSchemaNametest = "urn:com:example"

// Attributes to test with. Add optional properties to test here if desired.
type scimSubattributeTestModel struct {
	id                string
	scimAttributeName string
	scimSchemaName    string
}

func TestAccScimSubattribute(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := scimSubattributeTestModel{
		id:                testIdScimSubattribute,
		scimAttributeName: testScimAttributeName,
		scimSchemaName:    testScimSchemaNametest,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckScimSubattributeDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccScimSubattributeResource(resourceName, initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccScimSubattributeResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_scim_subattribute." + resourceName,
				ImportStateId:     initialResourceModel.scimSchemaName + "/" + initialResourceModel.scimAttributeName + "/" + initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccScimSubattributeResource(resourceName string, resourceModel scimSubattributeTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_scim_attribute" "%[3]s" {
  scim_schema_name = pingdirectory_scim_schema.myScimSchema.schema_urn
  name             = "%[3]s"
}

resource "pingdirectory_scim_schema" "myScimSchema" {
  schema_urn = "urn:com:example"
}

resource "pingdirectory_scim_subattribute" "%[1]s" {
  id                  = "%[2]s"
  scim_attribute_name = pingdirectory_scim_attribute.%[3]s.name
  scim_schema_name    = pingdirectory_scim_schema.myScimSchema.schema_urn
}`, resourceName,
		resourceModel.id,
		resourceModel.scimAttributeName,
		resourceModel.scimSchemaName)
}

// Test that any objects created by the test are destroyed
func testAccCheckScimSubattributeDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ScimSubattributeApi.GetScimSubattribute(ctx, testIdScimSubattribute, testScimAttributeName, testScimSchemaNametest).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Scim Subattribute", testIdScimSubattribute)
	}
	return nil
}
