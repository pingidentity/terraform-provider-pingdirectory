package virtualattribute_test

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

const testIdMirrorVirtualAttribute = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type mirrorVirtualAttributeTestModel struct {
	id              string
	sourceAttribute string
	enabled         bool
	attributeType   string
}

func TestAccMirrorVirtualAttribute(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := mirrorVirtualAttributeTestModel{
		id:              testIdMirrorVirtualAttribute,
		sourceAttribute: "mail",
		enabled:         true,
		attributeType:   "name",
	}
	updatedResourceModel := mirrorVirtualAttributeTestModel{
		id:              testIdMirrorVirtualAttribute,
		sourceAttribute: "cn",
		enabled:         false,
		attributeType:   "name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckMirrorVirtualAttributeDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccMirrorVirtualAttributeResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedMirrorVirtualAttributeAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccMirrorVirtualAttributeResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedMirrorVirtualAttributeAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccMirrorVirtualAttributeResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_virtual_attribute." + resourceName,
				ImportStateId:           updatedResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccMirrorVirtualAttributeResource(resourceName string, resourceModel mirrorVirtualAttributeTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_virtual_attribute" "%[1]s" {
  type             = "mirror"
  id               = "%[2]s"
  source_attribute = "%[3]s"
  enabled          = %[4]t
  attribute_type   = "%[5]s"
}`, resourceName, resourceModel.id,
		resourceModel.sourceAttribute,
		resourceModel.enabled,
		resourceModel.attributeType)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedMirrorVirtualAttributeAttributes(config mirrorVirtualAttributeTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.VirtualAttributeApi.GetVirtualAttribute(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Mirror Virtual Attribute"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "source-attribute",
			config.sourceAttribute, response.MirrorVirtualAttributeResponse.SourceAttribute)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.MirrorVirtualAttributeResponse.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "attribute-type",
			config.attributeType, response.MirrorVirtualAttributeResponse.AttributeType)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckMirrorVirtualAttributeDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.VirtualAttributeApi.GetVirtualAttribute(ctx, testIdMirrorVirtualAttribute).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Mirror Virtual Attribute", testIdMirrorVirtualAttribute)
	}
	return nil
}
