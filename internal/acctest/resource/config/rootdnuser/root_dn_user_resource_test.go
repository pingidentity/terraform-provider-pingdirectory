package rootdnuser_test

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

const testIdRootDnUser = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type rootDnUserTestModel struct {
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

func TestAccRootDnUser(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := rootDnUserTestModel{
		id:                           testIdRootDnUser,
		inheritDefaultRootPrivileges: true,
		searchResultEntryLimit:       0,
		timeLimitSeconds:             0,
		lookThroughEntryLimit:        0,
		idleTimeLimitSeconds:         0,
		passwordPolicy:               "Root Password Policy",
		requireSecureAuthentication:  false,
		requireSecureConnections:     false,
	}
	updatedResourceModel := rootDnUserTestModel{
		id:                           testIdRootDnUser,
		inheritDefaultRootPrivileges: false,
		searchResultEntryLimit:       1,
		timeLimitSeconds:             1,
		lookThroughEntryLimit:        1,
		idleTimeLimitSeconds:         1,
		passwordPolicy:               "Root Password Policy",
		requireSecureAuthentication:  true,
		requireSecureConnections:     true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckRootDnUserDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccRootDnUserResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedRootDnUserAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_root_dn_user.%s", resourceName), "inherit_default_root_privileges", strconv.FormatBool(initialResourceModel.inheritDefaultRootPrivileges)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_root_dn_user.%s", resourceName), "search_result_entry_limit", strconv.FormatInt(initialResourceModel.searchResultEntryLimit, 10)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_root_dn_user.%s", resourceName), "password_policy", initialResourceModel.passwordPolicy),
				),
			},
			{
				// Test updating some fields
				Config: testAccRootDnUserResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedRootDnUserAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccRootDnUserResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_root_dn_user." + resourceName,
				ImportStateId:           updatedResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccRootDnUserResource(resourceName string, resourceModel rootDnUserTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_root_dn_user" "%[1]s" {
  id                              = "%[2]s"
  inherit_default_root_privileges = %[3]t
  search_result_entry_limit       = %[4]d
  time_limit_seconds              = %[5]d
  look_through_entry_limit        = %[6]d
  idle_time_limit_seconds         = %[7]d
  password_policy                 = "%[8]s"
  require_secure_authentication   = %[9]t
  require_secure_connections      = %[10]t
}

data "pingdirectory_root_dn_user" "%[1]s" {
  id = "%[2]s"
  depends_on = [
    pingdirectory_root_dn_user.%[1]s
  ]
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
func testAccCheckExpectedRootDnUserAttributes(config rootDnUserTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.RootDnUserApi.GetRootDnUser(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Root Dn User"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "inherit-default-root-privileges",
			config.inheritDefaultRootPrivileges, response.InheritDefaultRootPrivileges)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "search-result-entry-limit",
			config.searchResultEntryLimit, response.SearchResultEntryLimit)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "time-limit-seconds",
			config.timeLimitSeconds, response.TimeLimitSeconds)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "look-through-entry-limit",
			config.lookThroughEntryLimit, response.LookThroughEntryLimit)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "idle-time-limit-seconds",
			config.idleTimeLimitSeconds, response.IdleTimeLimitSeconds)
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
func testAccCheckRootDnUserDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.RootDnUserApi.GetRootDnUser(ctx, testIdRootDnUser).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Root Dn User", testIdRootDnUser)
	}
	return nil
}
