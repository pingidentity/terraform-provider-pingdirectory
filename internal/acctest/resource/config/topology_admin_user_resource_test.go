package config_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdTopologyAdminUser = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type topologyAdminUserTestModel struct {
	id                           string
	inheritDefaultRootPrivileges bool
	searchResultEntryLimit       int64
	timeLimitSeconds             int64
	lookThroughEntryLimit        int64
	idleTimeLimitSeconds         int64
	passwordPolicy               string
	requireSecureAuthentication  bool
	requireSecureConnections     bool
}

func TestAccTopologyAdminUser(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := topologyAdminUserTestModel{
		id:                           testIdTopologyAdminUser,
		inheritDefaultRootPrivileges: true,
		searchResultEntryLimit:       100,
		timeLimitSeconds:             60,
		lookThroughEntryLimit:        20,
		idleTimeLimitSeconds:         120,
		passwordPolicy:               "Default Password Policy",
		requireSecureAuthentication:  true,
		requireSecureConnections:     false,
	}
	updatedResourceModel := topologyAdminUserTestModel{
		id:                           testIdTopologyAdminUser,
		inheritDefaultRootPrivileges: false,
		searchResultEntryLimit:       101,
		timeLimitSeconds:             61,
		lookThroughEntryLimit:        21,
		idleTimeLimitSeconds:         121,
		passwordPolicy:               "Root Password Policy",
		requireSecureAuthentication:  false,
		requireSecureConnections:     true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckTopologyAdminUserDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccTopologyAdminUserResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedTopologyAdminUserAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccTopologyAdminUserResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedTopologyAdminUserAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccTopologyAdminUserResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_topology_admin_user." + resourceName,
				ImportStateId:           updatedResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccTopologyAdminUserResource(resourceName string, resourceModel topologyAdminUserTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_topology_admin_user" "%[1]s" {
  id                              = "%[2]s"
  inherit_default_root_privileges = %[3]t
  search_result_entry_limit       = %[4]d
  time_limit_seconds              = %[5]d
  look_through_entry_limit        = %[6]d
  idle_time_limit_seconds         = %[7]d
  password_policy                 = "%[8]s"
  require_secure_authentication   = %[9]t
  require_secure_connections      = %[10]t
}`, resourceName, resourceModel.id,
		resourceModel.inheritDefaultRootPrivileges,
		resourceModel.searchResultEntryLimit,
		resourceModel.timeLimitSeconds,
		resourceModel.lookThroughEntryLimit,
		resourceModel.idleTimeLimitSeconds,
		resourceModel.passwordPolicy,
		resourceModel.requireSecureAuthentication,
		resourceModel.requireSecureConnections)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedTopologyAdminUserAttributes(config topologyAdminUserTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.TopologyAdminUserApi.GetTopologyAdminUser(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Topology Admin User"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "inherit-default-root-privileges",
			config.inheritDefaultRootPrivileges, response.InheritDefaultRootPrivileges)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "search-result-entry-limit",
			config.searchResultEntryLimit, int64(response.SearchResultEntryLimit))
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "time-limit-seconds",
			config.timeLimitSeconds, int64(response.TimeLimitSeconds))
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "look-through-entry-limit",
			config.lookThroughEntryLimit, int64(response.LookThroughEntryLimit))
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "idle-time-limit-seconds",
			config.idleTimeLimitSeconds, int64(response.IdleTimeLimitSeconds))
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "password-policy",
			config.passwordPolicy, response.PasswordPolicy)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "require-secure-authentication",
			config.requireSecureAuthentication, response.RequireSecureAuthentication)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "require-secure-connections",
			config.requireSecureConnections, response.RequireSecureConnections)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckTopologyAdminUserDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.TopologyAdminUserApi.GetTopologyAdminUser(ctx, testIdTopologyAdminUser).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Topology Admin User", testIdTopologyAdminUser)
	}
	return nil
}
