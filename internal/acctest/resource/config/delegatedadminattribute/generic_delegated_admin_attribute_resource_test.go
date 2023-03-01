package delegatedadminattribute_test

import (
	"fmt"
	"testing"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testIdGenericDelegatedAdminAttribute = "MyId"
const testDelegatedAdminResourceAttributeName = "myGenericDelegatedAdminResourceAttributeName"

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
		restResourceTypeName: testDelegatedAdminResourceAttributeName,
		attributeType:        "cn",
		displayName:          "Device Name",
		displayOrderIndex:    1,
	}
	updatedResourceModel := genericDelegatedAdminAttributeTestModel{
		restResourceTypeName: testDelegatedAdminResourceAttributeName,
		attributeType:        "cn",
		displayName:          "Device Name2",
		displayOrderIndex:    2,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckGenericDelegatedAdminAttributeDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccGenericDelegatedAdminAttributeResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedGenericDelegatedAdminAttributeAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccGenericDelegatedAdminAttributeResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedGenericDelegatedAdminAttributeAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccGenericDelegatedAdminAttributeResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_generic_delegated_admin_attribute." + resourceName,
				ImportStateId:     updatedResourceModel.restResourceTypeName + "/" + updatedResourceModel.attributeType,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccGenericDelegatedAdminAttributeResource(resourceName string, resourceModel genericDelegatedAdminAttributeTestModel) string {
	return fmt.Sprintf(`
		resource "pingdirectory_generic_rest_resource_type" "%[2]s" {
			id                                = "%[2]s"
			enabled                           = true
			resource_endpoint                 = "device"
			display_name                      = "Device"
			structural_ldap_objectclass       = "device"
			search_base_dn                    = "dc=example,dc=com"
			parent_dn                         = "dc=example,dc=com"
			search_filter_pattern             = "(cn=*%%%%*)"
			primary_display_attribute_type    = "cn"
}
		resource "pingdirectory_generic_delegated_admin_attribute" "%[1]s" {
        	rest_resource_type_name = pingdirectory_generic_rest_resource_type.%[2]s.id
	        attribute_type = "%[3]s"
	        display_name = "%[4]s"
	        display_order_index = %[5]d
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
		//response2, _, err := testClient.DelegatedAdminAttributeApi.GetDelegatedAdminAttributeExecute()

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
			config.displayOrderIndex, int64(response.GenericDelegatedAdminAttributeResponse.DisplayOrderIndex))
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
	_, _, err := testClient.DelegatedAdminAttributeApi.GetDelegatedAdminAttribute(ctx, testIdGenericDelegatedAdminAttribute, testDelegatedAdminResourceAttributeName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Generic Delegated Admin Attribute", testIdGenericDelegatedAdminAttribute)
	}
	return nil
}
