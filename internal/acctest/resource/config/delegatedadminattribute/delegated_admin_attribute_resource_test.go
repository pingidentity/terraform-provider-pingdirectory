package delegatedadminattribute_test

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

const testAttributeType = "cn"
const testParentResourceType = "myParentGenericResource"

// Attributes to test with. Add optional properties to test here if desired.
type genericDelegatedAdminAttributeTestModel struct {
	restResourceTypeName string
	attributeType        string
	displayName          string
	displayOrderIndex    int64
}

func TestAccGenericDelegatedAdminAttribute(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := genericDelegatedAdminAttributeTestModel{
		restResourceTypeName: testParentResourceType,
		attributeType:        testAttributeType,
		displayName:          "Device Name",
		displayOrderIndex:    1,
	}
	updatedResourceModel := genericDelegatedAdminAttributeTestModel{
		restResourceTypeName: testParentResourceType,
		attributeType:        testAttributeType,
		displayName:          "Device Name2",
		displayOrderIndex:    2,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckGenericDelegatedAdminAttributeDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccGenericDelegatedAdminAttributeResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedGenericDelegatedAdminAttributeAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_delegated_admin_attribute.%s", resourceName), "attribute_type", initialResourceModel.attributeType),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_delegated_admin_attribute.%s", resourceName), "display_name", initialResourceModel.displayName),
					resource.TestCheckResourceAttrSet("data.pingdirectory_delegated_admin_attributes.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccGenericDelegatedAdminAttributeResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedGenericDelegatedAdminAttributeAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccGenericDelegatedAdminAttributeResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_delegated_admin_attribute." + resourceName,
				ImportStateId:     updatedResourceModel.restResourceTypeName + "/" + updatedResourceModel.attributeType,
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
					_, err := testClient.DelegatedAdminAttributeApi.DeleteDelegatedAdminAttribute(ctx, updatedResourceModel.attributeType, updatedResourceModel.restResourceTypeName).Execute()
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

func testAccGenericDelegatedAdminAttributeResource(resourceName string, resourceModel genericDelegatedAdminAttributeTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_rest_resource_type" "%[2]s" {
  type                           = "generic"
  name                           = "%[2]s"
  enabled                        = true
  resource_endpoint              = "device"
  display_name                   = "Device"
  structural_ldap_objectclass    = "device"
  search_base_dn                 = "dc=example,dc=com"
  parent_dn                      = "dc=example,dc=com"
  search_filter_pattern          = "(cn=*%%%%*)"
  primary_display_attribute_type = "cn"
}

resource "pingdirectory_delegated_admin_attribute" "%[1]s" {
  type                    = "generic"
  rest_resource_type_name = pingdirectory_rest_resource_type.%[2]s.id
  attribute_type          = "%[3]s"
  display_name            = "%[4]s"
  display_order_index     = %[5]d
}

data "pingdirectory_delegated_admin_attribute" "%[1]s" {
  rest_resource_type_name = "%[2]s"
  attribute_type          = "%[3]s"
  depends_on = [
    pingdirectory_delegated_admin_attribute.%[1]s
  ]
}

data "pingdirectory_delegated_admin_attributes" "list" {
  rest_resource_type_name = "%[2]s"
  depends_on = [
    pingdirectory_delegated_admin_attribute.%[1]s
  ]
}`, resourceName,
		resourceModel.restResourceTypeName,
		resourceModel.attributeType,
		resourceModel.displayName,
		resourceModel.displayOrderIndex)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedGenericDelegatedAdminAttributeAttributes(config genericDelegatedAdminAttributeTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.DelegatedAdminAttributeApi.GetDelegatedAdminAttribute(ctx, config.attributeType, config.restResourceTypeName).Execute()

		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Generic Delegated Admin Attribute"
		err = acctest.TestAttributesMatchString(resourceType, &config.attributeType, "attribute-type",
			config.attributeType, response.GenericDelegatedAdminAttributeResponse.AttributeType)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.attributeType, "display-name",
			config.displayName, response.GenericDelegatedAdminAttributeResponse.DisplayName)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.attributeType, "display-order-index",
			config.displayOrderIndex, response.GenericDelegatedAdminAttributeResponse.DisplayOrderIndex)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckGenericDelegatedAdminAttributeDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.DelegatedAdminAttributeApi.GetDelegatedAdminAttribute(ctx, testAttributeType, testParentResourceType).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Generic Delegated Admin Attribute", testAttributeType)
	}
	return nil
}
