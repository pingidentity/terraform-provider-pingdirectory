package config_test

import (
	"fmt"
	"testing"

	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testIdDelegatedAdminResourceRights = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type delegatedAdminResourceRightsTestModel struct {
	adminResourceName string
	enabled           bool
	adminPermission   []string
	resourceName      string
}

func TestAccDelegatedAdminResourceRights(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := delegatedAdminResourceRightsTestModel{
		adminResourceName: testIdDelegatedAdminResourceRights,
		enabled:           true,
		adminPermission:   []string{"create", "read"},
	}
	updatedResourceModel := delegatedAdminResourceRightsTestModel{
		adminResourceName: testIdDelegatedAdminResourceRights,
		enabled:           false,
		adminPermission:   []string{"read"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckDelegatedAdminResourceRightsDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDelegatedAdminResourceRightsResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedDelegatedAdminResourceRightsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccDelegatedAdminResourceRightsResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedDelegatedAdminResourceRightsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccDelegatedAdminResourceRightsResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_delegated_admin_resource_rights." + resourceName,
				ImportStateId:           updatedResourceModel.adminResourceName + "/" + resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccDelegatedAdminResourceRightsResource(resourceName string, resourceModel delegatedAdminResourceRightsTestModel) string {
	return fmt.Sprintf(`
	resource "pingdirectory_user_rest_resource_type" "myUserRestResourceType" {
		id                          = "MyUserRestResourceType"
		enabled                     = true
		resource_endpoint           = "userRestResource"
		structural_ldap_objectclass = "inetOrgPerson"
		search_base_dn              = "cn=users,dc=test,dc=com"
}
	resource "pingdirectory_delegated_admin_rights" "myDelegatedAdminRights" {
		id            = "MyDelegatedAdminRights"
		enabled       = true
		admin_user_dn = "cn=admin-users,dc=test,dc=com"
}
	resource "pingdirectory_delegated_admin_resource_rights" "%[1]s" {
		delegated_admin_rights_name = pingdirectory_delegated_admin_rights.myDelegatedAdminRights.id
		admin_permission = %[2]s
		enabled = %[3]t
		rest_resource_type = pingdirectory_user_rest_resource_type.myUserRestResourceType.id
}`, resourceName,
		acctest.StringSliceToTerraformString(resourceModel.adminPermission),
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDelegatedAdminResourceRightsAttributes(config delegatedAdminResourceRightsTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.DelegatedAdminResourceRightsApi.GetDelegatedAdminResourceRights(ctx, config.resourceName, testIdDelegatedAdminResourceRights).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Delegated Admin Resource Rights"
		err = acctest.TestAttributesMatchBool(resourceType, &config.resourceName, "enabled",
			config.enabled, response.Enabled)
		if err != nil {
			return err
		}
		// err = acctest.TestAttributesMatchStringSlice(resourceType, nil, "admin-permission",
		// 	config.adminPermission, client.StringSliceEnumrootDnDefaultRootPrivilegeNameProp(response.AdminPermission))
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.resourceName, "admin-permission",
			config.adminPermission, client.StringSliceEnumdelegatedAdminResourceRightsAdminPermissionProp(response.AdminPermission))
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckDelegatedAdminResourceRightsDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.DelegatedAdminResourceRightsApi.GetDelegatedAdminResourceRights(ctx, testIdDelegatedAdminResourceRights, testIdDelegatedAdminResourceRights).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Delegated Admin Resource Rights", testIdDelegatedAdminResourceRights)
	}
	return nil
}
