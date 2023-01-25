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

// Testing will do three things
// First, read the state and check if default permissions match the slice I provide (unmodified PD instance)
// Second, apply the minimum permissions set and confirm that only those are there (expected = provided)
// Third, apply updated permissions set and confirm that "backend-restore" has been added back (expected = provided)

var defaultPermissionOne = "backend-backup"
var defaultPermissionTwo = "metrics-read"

// Overkill to have a struct, but makes the logic cleaner later
type rootDnTestModel struct {
	permissionsList []string
}

func TestAccRootDn(t *testing.T) {
	resourceName := "testrootdn"
	defaultResourceModel := rootDnTestModel{
		permissionsList: []string{"audit-data-security", "backend-backup", "backend-restore", "bypass-acl", "collect-support-data", "config-read", "config-write", "disconnect-client", "file-servlet-access", "ldif-export", "ldif-import", "lockdown-mode", "manage-topology", "metrics-read", "modify-acl", "password-reset", "permit-get-password-policy-state-issues", "privilege-change", "server-restart", "server-shutdown", "soft-delete-read", "stream-values", "third-party-task", "unindexed-search", "update-schema", "use-admin-session"},
	}
	minimumResourceModel := rootDnTestModel{
		permissionsList: []string{"bypass-acl", "config-read", "config-write", "modify-acl", "privilege-change", "use-admin-session"},
	}
	updatedResourceModel := rootDnTestModel{
		permissionsList: []string{"bypass-acl", "backend-restore", "config-read", "config-write", "modify-acl", "privilege-change", "use-admin-session"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				// Test defaults
				// load defaults from server
				Config: testAccRootDnResourceDefault(resourceName),
				Check: resource.ComposeTestCheckFunc(
					// check if the test default permissions are present in the state file
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("pingdirectory_root_dn.%s", resourceName), "default_root_privilege_name.*", defaultPermissionOne),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("pingdirectory_root_dn.%s", resourceName), "default_root_privilege_name.*", defaultPermissionTwo),

					// check if the permissions reported by PingDirectory match the state file
					testAccCheckExpectedRootDnPermissions(fmt.Sprintf("pingdirectory_root_dn.%s", resourceName), defaultResourceModel.permissionsList),
				),
			},

			{
				// Test after applying the minimum set
				Config: testAccRootDnResource(resourceName, minimumResourceModel),
				Check: resource.ComposeTestCheckFunc(
					// check that the permissions reported by PingDirectory match what was sent
					testAccCheckExpectedRootDnPermissions(fmt.Sprintf("pingdirectory_root_dn.%s", resourceName), minimumResourceModel.permissionsList),
				),
			},
			{
				// Test after applying updated permissions
				Config: testAccRootDnResource(resourceName, updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					// check that the permissions reported by PingDirectory match what was sent
					testAccCheckExpectedRootDnPermissions(fmt.Sprintf("pingdirectory_root_dn.%s", resourceName), updatedResourceModel.permissionsList),
				),
			},
			{
				// Set permissions back to default for other tests
				Config: testAccRootDnResource(resourceName, defaultResourceModel),
				Check: resource.ComposeTestCheckFunc(
					// check if the permissions reported by PingDirectory match the state file
					testAccCheckExpectedRootDnPermissions(fmt.Sprintf("pingdirectory_root_dn.%s", resourceName), defaultResourceModel.permissionsList),
				),
			},
		},
	})
}

// empty resource object means all values are computed, so it will retreive defaults from PD
func testAccRootDnResourceDefault(resourceName string) string {
	return fmt.Sprintf(`
resource "pingdirectory_root_dn" "%[1]s" {
}`, resourceName)
}

// apply a list of permissions to the default_root_privilege_name object
// only what is in the list should exist after applied
func testAccRootDnResource(resourceName string, resourceModel rootDnTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_root_dn" "%[1]s" {
	default_root_privilege_name = %[2]s
}`, resourceName, acctest.StringSliceToTerraformString(resourceModel.permissionsList))
}

// Test that the expected RootDN permissions are set on the PingDirectory server
func testAccCheckExpectedRootDnPermissions(resourceName string, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//		resourceType := "rootdn"
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
