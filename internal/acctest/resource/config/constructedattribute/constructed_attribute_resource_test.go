package constructedattribute_test

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

const testIdConstructedAttribute = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type constructedAttributeTestModel struct {
	id            string
	attributeType string
	valuePattern  []string
}

func TestAccConstructedAttribute(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := constructedAttributeTestModel{
		id:            testIdConstructedAttribute,
		attributeType: "cn",
		valuePattern:  []string{"{attr-name}"},
	}
	updatedResourceModel := constructedAttributeTestModel{
		id:            testIdConstructedAttribute,
		attributeType: "mail",
		valuePattern:  []string{"{userMail}"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckConstructedAttributeDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccConstructedAttributeResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedConstructedAttributeAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_constructed_attribute.%s", resourceName), "attribute_type", initialResourceModel.attributeType),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_constructed_attribute.%s", resourceName), "value_pattern.*", initialResourceModel.valuePattern[0]),
					resource.TestCheckResourceAttrSet("data.pingdirectory_constructed_attributes.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccConstructedAttributeResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedConstructedAttributeAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccConstructedAttributeResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_constructed_attribute." + resourceName,
				ImportStateId:     updatedResourceModel.id,
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
					_, err := testClient.ConstructedAttributeApi.DeleteConstructedAttribute(ctx, updatedResourceModel.id).Execute()
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

func testAccConstructedAttributeResource(resourceName string, resourceModel constructedAttributeTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_constructed_attribute" "%[1]s" {
  name           = "%[2]s"
  attribute_type = "%[3]s"
  value_pattern  = %[4]s
}

data "pingdirectory_constructed_attribute" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_constructed_attribute.%[1]s
  ]
}

data "pingdirectory_constructed_attributes" "list" {
  depends_on = [
    pingdirectory_constructed_attribute.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.attributeType,
		acctest.StringSliceToTerraformString(resourceModel.valuePattern))
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedConstructedAttributeAttributes(config constructedAttributeTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ConstructedAttributeApi.GetConstructedAttribute(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Constructed Attribute"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "attribute-type",
			config.attributeType, response.AttributeType)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "value-pattern",
			config.valuePattern, response.ValuePattern)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckConstructedAttributeDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ConstructedAttributeApi.GetConstructedAttribute(ctx, testIdConstructedAttribute).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Constructed Attribute", testIdConstructedAttribute)
	}
	return nil
}
