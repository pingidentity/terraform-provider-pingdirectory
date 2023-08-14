package scimschema_test

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

const testIdScimSchema = "urn:com:example:scimschematest"

// Attributes to test with. Add optional properties to test here if desired.
type scimSchemaTestModel struct {
	schemaUrn   string
	description string
}

func TestAccScimSchema(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := scimSchemaTestModel{
		schemaUrn:   testIdScimSchema,
		description: "initial",
	}
	updatedResourceModel := scimSchemaTestModel{
		schemaUrn:   testIdScimSchema,
		description: "updated",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckScimSchemaDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccScimSchemaResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedScimSchemaAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_scim_schema.%s", resourceName), "description", initialResourceModel.description),
					resource.TestCheckResourceAttrSet("data.pingdirectory_scim_schemas.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccScimSchemaResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedScimSchemaAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccScimSchemaResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_scim_schema." + resourceName,
				ImportStateId:     updatedResourceModel.schemaUrn,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.ScimSchemaApi.DeleteScimSchema(ctx, updatedResourceModel.schemaUrn).Execute()
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

func testAccScimSchemaResource(resourceName string, resourceModel scimSchemaTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_scim_schema" "%[1]s" {
  schema_urn  = "%[2]s"
  description = "%[3]s"
}

data "pingdirectory_scim_schema" "%[1]s" {
  schema_urn = "%[2]s"
  depends_on = [
    pingdirectory_scim_schema.%[1]s
  ]
}

data "pingdirectory_scim_schemas" "list" {
  depends_on = [
    pingdirectory_scim_schema.%[1]s
  ]
}`, resourceName,
		resourceModel.schemaUrn,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedScimSchemaAttributes(config scimSchemaTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ScimSchemaApi.GetScimSchema(ctx, config.schemaUrn).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Scim Schema"
		err = acctest.TestAttributesMatchString(resourceType, &config.schemaUrn, "schema-urn",
			config.schemaUrn, response.SchemaURN)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.description, "description",
			config.description, response.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckScimSchemaDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ScimSchemaApi.GetScimSchema(ctx, testIdScimSchema).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Scim Schema", testIdScimSchema)
	}
	return nil
}
