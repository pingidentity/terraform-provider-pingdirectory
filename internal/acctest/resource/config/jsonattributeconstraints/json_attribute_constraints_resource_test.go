package jsonattributeconstraints_test

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

const testIdJsonAttributeConstraints = "ubidEntitlement"

// Attributes to test with. Add optional properties to test here if desired.
type jsonAttributeConstraintsTestModel struct {
	attributeType        string
	description          string
	allow_unnamed_fields bool
}

func TestAccJsonAttributeConstraints(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := jsonAttributeConstraintsTestModel{
		attributeType:        testIdJsonAttributeConstraints,
		description:          "Initial JSON attribute constraint",
		allow_unnamed_fields: false,
	}
	updatedResourceModel := jsonAttributeConstraintsTestModel{
		attributeType:        testIdJsonAttributeConstraints,
		description:          "Updated JSON attribute constraint",
		allow_unnamed_fields: true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckJsonAttributeConstraintsDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccJsonAttributeConstraintsResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedJsonAttributeConstraintsAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_json_attribute_constraints.%s", resourceName), "description", initialResourceModel.description),
					resource.TestCheckResourceAttrSet("data.pingdirectory_json_attribute_constraints_list.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccJsonAttributeConstraintsResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedJsonAttributeConstraintsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccJsonAttributeConstraintsResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_json_attribute_constraints." + resourceName,
				ImportStateId:     updatedResourceModel.attributeType,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccJsonAttributeConstraintsResource(resourceName string, resourceModel jsonAttributeConstraintsTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_json_attribute_constraints" "%[1]s" {
  attribute_type       = "%[2]s"
  description          = "%[3]s"
  allow_unnamed_fields = %[4]t
}

data "pingdirectory_json_attribute_constraints" "%[1]s" {
  attribute_type = "%[2]s"
  depends_on = [
    pingdirectory_json_attribute_constraints.%[1]s
  ]
}

data "pingdirectory_json_attribute_constraints_list" "list" {
  depends_on = [
    pingdirectory_json_attribute_constraints.%[1]s
  ]
}`, resourceName,
		resourceModel.attributeType,
		resourceModel.description,
		resourceModel.allow_unnamed_fields)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedJsonAttributeConstraintsAttributes(config jsonAttributeConstraintsTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.JsonAttributeConstraintsApi.GetJsonAttributeConstraints(ctx, config.attributeType).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Json Attribute Constraints"
		err = acctest.TestAttributesMatchString(resourceType, &config.attributeType, "attribute-type",
			config.attributeType, response.AttributeType)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.attributeType, "description",
			config.description, *response.Description)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.attributeType, "allow-unnamed-fields",
			config.allow_unnamed_fields, *response.AllowUnnamedFields)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckJsonAttributeConstraintsDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.JsonAttributeConstraintsApi.GetJsonAttributeConstraints(ctx, testIdJsonAttributeConstraints).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Json Attribute Constraints", testIdJsonAttributeConstraints)
	}
	return nil
}
