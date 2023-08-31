package delegatedadminresourcerights_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdDelegatedAdminResourceRights = "MyDARightsId"
const testDelegatedAdminRightsName = "myDelegatedAdminRightsName"

// Attributes to test with. Add optional properties to test here if desired.
type delegatedAdminResourceRightsTestModel struct {
	delegatedAdminRightsName string
	enabled                  bool
	restResourceType         string
	adminPermission          []string
}

func TestAccDelegatedAdminResourceRights(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := delegatedAdminResourceRightsTestModel{
		delegatedAdminRightsName: testDelegatedAdminRightsName,
		enabled:                  true,
		restResourceType:         testIdDelegatedAdminResourceRights,
		adminPermission:          []string{"read"},
	}
	updatedResourceModel := delegatedAdminResourceRightsTestModel{
		delegatedAdminRightsName: testDelegatedAdminRightsName,
		enabled:                  false,
		restResourceType:         testIdDelegatedAdminResourceRights,
		adminPermission:          []string{"create", "read"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckDelegatedAdminResourceRightsDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDelegatedAdminResourceRightsResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedDelegatedAdminResourceRightsAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_delegated_admin_resource_rights.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_delegated_admin_resource_rights.%s", resourceName), "rest_resource_type", initialResourceModel.restResourceType),
					resource.TestCheckResourceAttrSet("data.pingdirectory_delegated_admin_resource_rights_list.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccDelegatedAdminResourceRightsResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedDelegatedAdminResourceRightsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccDelegatedAdminResourceRightsResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_delegated_admin_resource_rights." + resourceName,
				ImportStateId:     updatedResourceModel.delegatedAdminRightsName + "/" + updatedResourceModel.restResourceType,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.DelegatedAdminResourceRightsApi.DeleteDelegatedAdminResourceRights(ctx, updatedResourceModel.restResourceType, updatedResourceModel.delegatedAdminRightsName).Execute()
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

func testAccDelegatedAdminResourceRightsResource(resourceName string, resourceModel delegatedAdminResourceRightsTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_rest_resource_type" "%[4]s" {
  type                        = "user"
  name                        = "%[4]s"
  enabled                     = true
  resource_endpoint           = "userRestResourceDelegatedAdminResourceRightsTest"
  structural_ldap_objectclass = "inetOrgPerson"
  search_base_dn              = "cn=users,dc=test,dc=com"
}

resource "pingdirectory_delegated_admin_rights" "%[2]s" {
  name          = "%[2]s"
  enabled       = true
  admin_user_dn = "cn=admin-users,dc=test,dc=com"
}

resource "pingdirectory_delegated_admin_resource_rights" "%[1]s" {
  delegated_admin_rights_name = pingdirectory_delegated_admin_rights.%[2]s.id
  admin_permission            = %[5]s
  enabled                     = %[3]t
  rest_resource_type          = pingdirectory_rest_resource_type.%[4]s.id
}

data "pingdirectory_delegated_admin_resource_rights" "%[1]s" {
  delegated_admin_rights_name = "%[2]s"
  rest_resource_type          = "%[4]s"
  depends_on = [
    pingdirectory_delegated_admin_resource_rights.%[1]s
  ]
}

data "pingdirectory_delegated_admin_resource_rights_list" "list" {
  delegated_admin_rights_name = "%[2]s"
  depends_on = [
    pingdirectory_delegated_admin_resource_rights.%[1]s
  ]
}`, resourceName,
		resourceModel.delegatedAdminRightsName,
		resourceModel.enabled,
		resourceModel.restResourceType,
		acctest.StringSliceToTerraformString(resourceModel.adminPermission))
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDelegatedAdminResourceRightsAttributes(config delegatedAdminResourceRightsTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.DelegatedAdminResourceRightsApi.GetDelegatedAdminResourceRights(ctx, config.restResourceType, config.delegatedAdminRightsName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Delegated Admin Resource Rights"
		err = acctest.TestAttributesMatchBool(resourceType, &config.restResourceType, "enabled",
			config.enabled, response.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.restResourceType, "admin-permission",
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
	_, _, err := testClient.DelegatedAdminResourceRightsApi.GetDelegatedAdminResourceRights(ctx, testIdDelegatedAdminResourceRights, testDelegatedAdminRightsName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Delegated Admin Resource Rights", testIdDelegatedAdminResourceRights)
	}
	return nil
}
