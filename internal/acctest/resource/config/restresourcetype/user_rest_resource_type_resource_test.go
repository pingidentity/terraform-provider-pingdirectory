package restresourcetype_test

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

const testIdUserRestResourceType = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type userRestResourceTypeTestModel struct {
	id                        string
	enabled                   bool
	resourceEndpoint          string
	structuralLdapObjectclass string
	searchBaseDn              string
}

func TestAccUserRestResourceType(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := userRestResourceTypeTestModel{
		id:                        testIdUserRestResourceType,
		enabled:                   true,
		resourceEndpoint:          "userRestResourceTest",
		structuralLdapObjectclass: "inetOrgPerson",
		searchBaseDn:              "cn=users,dc=test,dc=com",
	}
	updatedResourceModel := userRestResourceTypeTestModel{
		id:                        testIdUserRestResourceType,
		enabled:                   true,
		resourceEndpoint:          "userRestResourceTest",
		structuralLdapObjectclass: "inetOrgPerson",
		searchBaseDn:              "cn=users1,dc=test,dc=com",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckUserRestResourceTypeDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccUserRestResourceTypeResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedUserRestResourceTypeAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccUserRestResourceTypeResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedUserRestResourceTypeAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccUserRestResourceTypeResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_user_rest_resource_type." + resourceName,
				ImportStateId:           updatedResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccUserRestResourceTypeResource(resourceName string, resourceModel userRestResourceTypeTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_user_rest_resource_type" "%[1]s" {
	 id = "%[2]s"
	 enabled = %[3]t
	 resource_endpoint = "%[4]s"
	 structural_ldap_objectclass = "%[5]s"
	 search_base_dn = "%[6]s"
}`, resourceName, resourceModel.id,
		resourceModel.enabled,
		resourceModel.resourceEndpoint,
		resourceModel.structuralLdapObjectclass,
		resourceModel.searchBaseDn)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedUserRestResourceTypeAttributes(config userRestResourceTypeTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.RestResourceTypeApi.GetRestResourceType(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "User Rest Resource Type"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.UserRestResourceTypeResponse.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "resource-endpoint",
			config.resourceEndpoint, response.UserRestResourceTypeResponse.ResourceEndpoint)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "structural-ldap-objectclass",
			config.structuralLdapObjectclass, response.UserRestResourceTypeResponse.StructuralLDAPObjectclass)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "search-base-dn",
			config.searchBaseDn, response.UserRestResourceTypeResponse.SearchBaseDN)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckUserRestResourceTypeDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.RestResourceTypeApi.GetRestResourceType(ctx, testIdUserRestResourceType).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("User Rest Resource Type", testIdUserRestResourceType)
	}
	return nil
}
