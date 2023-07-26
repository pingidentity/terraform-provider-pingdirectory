package scimattribute_test

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

const testIdScimAttribute = "MyId"
const testScimSchemaName = "urn:com:example:scimattributetest"

// Attributes to test with. Add optional properties to test here if desired.
type scimAttributeTestModel struct {
	scimSchemaName string
	name           string
	description    string
}

func TestAccScimAttribute(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := scimAttributeTestModel{
		scimSchemaName: testScimSchemaName,
		name:           testIdScimAttribute,
		description:    "initial",
	}
	updatedResourceModel := scimAttributeTestModel{
		scimSchemaName: testScimSchemaName,
		name:           testIdScimAttribute,
		description:    "updated",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckScimAttributeDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccScimAttributeResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedScimAttributeAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_scim_attribute.%s", resourceName), "description", initialResourceModel.description),
				),
			},
			{
				// Test updating some fields
				Config: testAccScimAttributeResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedScimAttributeAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccScimAttributeResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_scim_attribute." + resourceName,
				ImportStateId:     updatedResourceModel.scimSchemaName + "/" + updatedResourceModel.name,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccScimAttributeResource(resourceName string, resourceModel scimAttributeTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_scim_schema" "myScimSchema" {
  schema_urn = "%[2]s"
}

resource "pingdirectory_scim_attribute" "%[1]s" {
  scim_schema_name = pingdirectory_scim_schema.myScimSchema.schema_urn
  name             = "%[3]s"
  description      = "%[4]s"
}

data "pingdirectory_scim_attribute" "%[1]s" {
	 scim_schema_name = "%[2]s"
	 name = "%[3]s"
  depends_on = [
    pingdirectory_scim_attribute.%[1]s
  ]
}`, resourceName,
		resourceModel.scimSchemaName,
		resourceModel.name,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedScimAttributeAttributes(config scimAttributeTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ScimAttributeApi.GetScimAttribute(ctx, config.name, config.scimSchemaName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Scim Attribute"
		err = acctest.TestAttributesMatchString(resourceType, &config.name, "name",
			config.name, response.Name)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.name, "description",
			config.description, response.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckScimAttributeDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ScimAttributeApi.GetScimAttribute(ctx, testIdScimAttribute, testScimSchemaName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Scim Attribute", testIdScimAttribute)
	}
	return nil
}
