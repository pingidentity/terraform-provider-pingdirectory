package backend_test

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

const testIdLocalDbBackend = "MyLocalDbBackend"

// Attributes to test with. Add optional properties to test here if desired.
type localDbBackendTestModel struct {
	backendId           string
	baseDn              []string
	writabilityMode     string
	dbDirectory         string
	importTempDirectory string
	enabled             bool
}

func TestAccLocalDbBackend(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := localDbBackendTestModel{
		backendId:           testIdLocalDbBackend,
		baseDn:              []string{"dc=example1,dc=com"},
		writabilityMode:     "enabled",
		dbDirectory:         "db",
		importTempDirectory: "tmp",
		enabled:             true,
	}
	updatedResourceModel := localDbBackendTestModel{
		backendId:           testIdLocalDbBackend,
		baseDn:              []string{"dc=example2,dc=com"},
		writabilityMode:     "internal-only",
		dbDirectory:         "db",
		importTempDirectory: "tmp/test",
		enabled:             false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckLocalDbBackendDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccLocalDbBackendResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedLocalDbBackendAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccLocalDbBackendResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedLocalDbBackendAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccLocalDbBackendResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_local_db_backend." + resourceName,
				ImportStateId:           updatedResourceModel.backendId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccLocalDbBackendResource(resourceName string, resourceModel localDbBackendTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_local_db_backend" "%[1]s" {
  backend_id            = "%[2]s"
  base_dn               = %[3]s
  writability_mode      = "%[4]s"
  db_directory          = "%[5]s"
  import_temp_directory = "%[6]s"
  enabled               = %[7]t
}`, resourceName,
		resourceModel.backendId,
		acctest.StringSliceToTerraformString(resourceModel.baseDn),
		resourceModel.writabilityMode,
		resourceModel.dbDirectory,
		resourceModel.importTempDirectory,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedLocalDbBackendAttributes(config localDbBackendTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.BackendApi.GetBackend(ctx, config.backendId).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Local Database Backend"
		err = acctest.TestAttributesMatchString(resourceType, &config.backendId, "backend-id",
			config.backendId, response.LocalDbBackendResponse.BackendID)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.backendId, "base-dn",
			config.baseDn, response.LocalDbBackendResponse.BaseDN)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.backendId, "writability-mode",
			config.writabilityMode, response.LocalDbBackendResponse.WritabilityMode.String())
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.backendId, "db-directory",
			config.dbDirectory, response.LocalDbBackendResponse.DbDirectory)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.backendId, "import-temp-directory",
			config.importTempDirectory, response.LocalDbBackendResponse.ImportTempDirectory)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.backendId, "enabled",
			config.enabled, response.LocalDbBackendResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckLocalDbBackendDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.RequestCriteriaApi.GetRequestCriteria(ctx, testIdLocalDbBackend).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Local Database Backend", testIdLocalDbBackend)
	}
	return nil
}
