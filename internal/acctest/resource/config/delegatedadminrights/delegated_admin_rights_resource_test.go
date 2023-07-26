package delegatedadminrights_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdDelegatedAdminRights = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type delegatedAdminRightsTestModel struct {
	id          string
	enabled     bool
	adminUserDn string
}

func TestAccDelegatedAdminRights(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := delegatedAdminRightsTestModel{
		id:          testIdDelegatedAdminRights,
		enabled:     true,
		adminUserDn: "cn=admin-users,dc=test,dc=com",
	}
	updatedResourceModel := delegatedAdminRightsTestModel{
		id:          testIdDelegatedAdminRights,
		enabled:     false,
		adminUserDn: "cn=other-admin-users,dc=test,dc=com",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckDelegatedAdminRightsDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDelegatedAdminRightsResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedDelegatedAdminRightsAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_delegated_admin_rights.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
				),
			},
			{
				// Test updating some fields
				Config: testAccDelegatedAdminRightsResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedDelegatedAdminRightsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccDelegatedAdminRightsResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_delegated_admin_rights." + resourceName,
				ImportStateId:           updatedResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccDelegatedAdminRightsResource(resourceName string, resourceModel delegatedAdminRightsTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_delegated_admin_rights" "%[1]s" {
  id            = "%[2]s"
  enabled       = %[3]t
  admin_user_dn = "%[4]s"
}

data "pingdirectory_delegated_admin_rights" "%[1]s" {
	 id = "%[2]s"
  depends_on = [
    pingdirectory_delegated_admin_rights.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.enabled,
		resourceModel.adminUserDn)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDelegatedAdminRightsAttributes(config delegatedAdminRightsTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.DelegatedAdminRightsApi.GetDelegatedAdminRights(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Delegated Admin Rights"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "admin-user-dn",
			config.adminUserDn, *response.AdminUserDN)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckDelegatedAdminRightsDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.DelegatedAdminRightsApi.GetDelegatedAdminRights(ctx, testIdDelegatedAdminRights).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Delegated Admin Rights", testIdDelegatedAdminRights)
	}
	return nil
}
