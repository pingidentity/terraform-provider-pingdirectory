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

const testIdScimAttributeMapping = "MyId"
const testScimResourceTypeName = "MyLdapMappingScimResourceType"

// Attributes to test with. Add optional properties to test here if desired.
type scimAttributeMappingTestModel struct {
	id                        string
	scimResourceTypeName      string
	scimResourceTypeAttribute string
	ldapAttribute             string
}

func TestAccScimAttributeMapping(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := scimAttributeMappingTestModel{
		id:                        testIdScimAttributeMapping,
		scimResourceTypeName:      testScimResourceTypeName,
		scimResourceTypeAttribute: "name",
		ldapAttribute:             "name",
	}
	updatedResourceModel := scimAttributeMappingTestModel{
		id:                        testIdScimAttributeMapping,
		scimResourceTypeName:      testScimResourceTypeName,
		scimResourceTypeAttribute: "givenName",
		ldapAttribute:             "givenName",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckScimAttributeMappingDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccScimAttributeMappingResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedScimAttributeMappingAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccScimAttributeMappingResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedScimAttributeMappingAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccScimAttributeMappingResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_scim_attribute_mapping." + resourceName,
				ImportStateId:     updatedResourceModel.scimResourceTypeName + "/" + updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccScimAttributeMappingResource(resourceName string, resourceModel scimAttributeMappingTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_scim_schema" "mySchema" {
     schema_urn = "urn:com:example"
}

resource "pingdirectory_ldap_mapping_scim_resource_type" "myLdapMappingScimResourceType" {
	 id          = "%[3]s"
	 core_schema = pingdirectory_scim_schema.mySchema.schema_urn
	 enabled     = false
	 endpoint    = "myendpoint"
}

resource "pingdirectory_scim_attribute_mapping" "%[1]s" {
	 id = "%[2]s"
	 scim_resource_type_name = pingdirectory_ldap_mapping_scim_resource_type.myLdapMappingScimResourceType.id
	 scim_resource_type_attribute = "%[4]s"
	 ldap_attribute = "%[5]s"
}`, resourceName,
		resourceModel.id,
		resourceModel.scimResourceTypeName,
		resourceModel.scimResourceTypeAttribute,
		resourceModel.ldapAttribute)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedScimAttributeMappingAttributes(config scimAttributeMappingTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ScimAttributeMappingApi.GetScimAttributeMapping(ctx, config.id, config.scimResourceTypeName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Scim Attribute Mapping"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "scim-resource-type-attribute",
			config.scimResourceTypeAttribute, response.ScimResourceTypeAttribute)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "ldap-attribute",
			config.ldapAttribute, response.LdapAttribute)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckScimAttributeMappingDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ScimAttributeMappingApi.GetScimAttributeMapping(ctx, testIdScimAttributeMapping, testScimResourceTypeName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Scim Attribute Mapping", testIdScimAttributeMapping)
	}
	return nil
}
