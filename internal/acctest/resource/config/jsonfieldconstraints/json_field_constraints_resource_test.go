package jsonfieldconstraints_test

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

const testIdJsonFieldConstraints = "id"
const testJsonAttributeConstraintsName = "ubidEmailJSON"

// Attributes to test with. Add optional properties to test here if desired.
type jsonFieldConstraintsTestModel struct {
	jsonAttributeConstraintsName string
	jsonField                    string
	valueType                    string
}

func TestAccJsonFieldConstraints(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := jsonFieldConstraintsTestModel{
		jsonAttributeConstraintsName: testJsonAttributeConstraintsName,
		jsonField:                    testIdJsonFieldConstraints,
		valueType:                    "string",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckJsonFieldConstraintsDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccJsonFieldConstraintsResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedJsonFieldConstraintsAttributes(initialResourceModel),
			},
		},
	})
}

func testAccJsonFieldConstraintsResource(resourceName string, resourceModel jsonFieldConstraintsTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_json_field_constraints" "%[1]s" {
  json_attribute_constraints_name = "%[2]s"
  json_field                      = "%[3]s"
  value_type                      = "%[4]s"
}`, resourceName,
		resourceModel.jsonAttributeConstraintsName,
		resourceModel.jsonField,
		resourceModel.valueType)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedJsonFieldConstraintsAttributes(config jsonFieldConstraintsTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.JsonFieldConstraintsApi.GetJsonFieldConstraints(ctx, config.jsonField, config.jsonAttributeConstraintsName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Json Field Constraints"
		err = acctest.TestAttributesMatchString(resourceType, &config.jsonField, "json-field",
			config.jsonField, response.JsonField)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.jsonField, "value-type",
			config.valueType, response.ValueType.String())
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckJsonFieldConstraintsDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.JsonFieldConstraintsApi.GetJsonFieldConstraints(ctx, testIdJsonFieldConstraints, testJsonAttributeConstraintsName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Json Field Constraints", testIdJsonFieldConstraints)
	}
	return nil
}
