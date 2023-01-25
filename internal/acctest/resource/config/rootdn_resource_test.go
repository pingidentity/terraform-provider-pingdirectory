package config_test

import (
	"fmt"
	"testing"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Testing will do four things
//   1) Read the state prior to making changes (unmodified PD instance) and check if default permissions match expected defaults (expected = provided)
//   2) Apply the minimum permissions set and confirm that only those are there (expected = provided)
//   3) Apply updated permissions set and confirm that "backend-restore" has been added back (expected = provided)
//   4) Apply the default permissions just in case they might impact other tests

var defaultPermissionOne = "backend-backup"
var defaultPermissionTwo = "metrics-read"

func TestAccRootDn(t *testing.T) {
	resourceName := "testrootdn"
	// default permissions as of January 2023, PingDirectory 9.1.0.0
	defaultPermissionsList := []string{"audit-data-security", "backend-backup", "backend-restore", "bypass-acl", "collect-support-data", "config-read", "config-write", "disconnect-client", "file-servlet-access", "ldif-export", "ldif-import", "lockdown-mode", "manage-topology", "metrics-read", "modify-acl", "password-reset", "permit-get-password-policy-state-issues", "privilege-change", "server-restart", "server-shutdown", "soft-delete-read", "stream-values", "third-party-task", "unindexed-search", "update-schema", "use-admin-session"}
	//
	minimumPermissionsList := []string{"bypass-acl", "config-read", "config-write", "modify-acl", "privilege-change", "use-admin-session"}
	updatedPermissionsList := []string{"bypass-acl", "backend-restore", "config-read", "config-write", "modify-acl", "privilege-change", "use-admin-session"}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				// Test defaults
				// load defaults from server with empty resource definition call
				Config: testAccRootDnResourceDefault(resourceName),
				Check: resource.ComposeTestCheckFunc(
					// check if the sample set of test default permissions are present in the state file
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("pingdirectory_root_dn.%s", resourceName), "default_root_privilege_name.*", defaultPermissionOne),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("pingdirectory_root_dn.%s", resourceName), "default_root_privilege_name.*", defaultPermissionTwo),

					// check if the permissions reported by PingDirectory match the state file
					testAccCheckExpectedRootDnPermissions(fmt.Sprintf("pingdirectory_root_dn.%s", resourceName), defaultPermissionsList),
				),
			},

			{
				// Test after applying the minimum set
				Config: testAccRootDnResource(resourceName, minimumPermissionsList),
				Check: resource.ComposeTestCheckFunc(
					// check that the permissions reported by PingDirectory match what was sent
					testAccCheckExpectedRootDnPermissions(fmt.Sprintf("pingdirectory_root_dn.%s", resourceName), minimumPermissionsList),
				),
			},
			{
				// Test after applying updated permissions
				Config: testAccRootDnResource(resourceName, updatedPermissionsList),
				Check: resource.ComposeTestCheckFunc(
					// check that the permissions reported by PingDirectory match what was sent
					testAccCheckExpectedRootDnPermissions(fmt.Sprintf("pingdirectory_root_dn.%s", resourceName), updatedPermissionsList),
				),
			},
			{
				// Set permissions back to default for other tests
				Config: testAccRootDnResource(resourceName, defaultPermissionsList),
				Check: resource.ComposeTestCheckFunc(
					// check if the permissions reported by PingDirectory match the state file
					testAccCheckExpectedRootDnPermissions(fmt.Sprintf("pingdirectory_root_dn.%s", resourceName), defaultPermissionsList),
				),
			},
		},
	})
}

// empty resource object means all values are computed, so it will retrieve defaults from PD
func testAccRootDnResourceDefault(resourceName string) string {
	return fmt.Sprintf(`
resource "pingdirectory_root_dn" "%[1]s" {
}`, resourceName)
}

// apply a list of permissions to the default_root_privilege_name object
// only what is in the list should exist after applied
func testAccRootDnResource(resourceName string, permissionsList []string) string {
	return fmt.Sprintf(`
resource "pingdirectory_root_dn" "%[1]s" {
	default_root_privilege_name = %[2]s
}`, resourceName, acctest.StringSliceToTerraformString(permissionsList))
}

// Test that the expected RootDN permissions are set on the PingDirectory server
func testAccCheckExpectedRootDnPermissions(resourceName string, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		rootDnResponse, _, err := testClient.RootDnApi.GetRootDn(ctx).Execute()
		if err != nil {
			return err
		}
		// Verify that permission matches expected
		err = acctest.TestAttributesMatchStringSlice("rootDn", nil, "default_root_privilege_name", expected, internaltypes.GetEnumStringSlice(rootDnResponse.DefaultRootPrivilegeName))
		if err != nil {
			return err
		}

		return nil

	}

}
