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

const testIdSensitiveAttribute = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type sensitiveAttributeTestModel struct {
	id            string
	attributeType []string
}

func TestAccSensitiveAttribute(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := sensitiveAttributeTestModel{
		id:            testIdSensitiveAttribute,
		attributeType: []string{"userPassword", "ds-pwp-retired-password"},
	}
	updatedResourceModel := sensitiveAttributeTestModel{
		id:            testIdSensitiveAttribute,
		attributeType: []string{"pwdHistory"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckSensitiveAttributeDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSensitiveAttributeResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSensitiveAttributeAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccSensitiveAttributeResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSensitiveAttributeAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSensitiveAttributeResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_sensitive_attribute." + resourceName,
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

func testAccSensitiveAttributeResource(resourceName string, resourceModel sensitiveAttributeTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_sensitive_attribute" "%[1]s" {
  id             = "%[2]s"
  attribute_type = %[3]s
}`, resourceName,
		resourceModel.id,
		acctest.StringSliceToTerraformString(resourceModel.attributeType))
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSensitiveAttributeAttributes(config sensitiveAttributeTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.SensitiveAttributeApi.GetSensitiveAttribute(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Sensitive Attribute"
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "attribute-type",
			config.attributeType, response.AttributeType)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSensitiveAttributeDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.SensitiveAttributeApi.GetSensitiveAttribute(ctx, testIdSensitiveAttribute).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Sensitive Attribute", testIdSensitiveAttribute)
	}
	return nil
}
