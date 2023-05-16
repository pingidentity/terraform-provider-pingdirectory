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

const testIdLocalDbVlvIndex = "MyId"
const testBackendNametest = "MyBackend"

// Attributes to test with. Add optional properties to test here if desired.
type localDbVlvIndexTestModel struct {
	backendName string
	baseDn      string
	scope       string
	filter      string
	sortOrder   string
	name        string
}

func TestAccLocalDbVlvIndex(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := localDbVlvIndexTestModel{
		backendName: testBackendNametest,
		baseDn:      "dc=example,dc=com",
		scope:       "base-object",
		filter:      "uid=user.1",
		sortOrder:   "givenName",
		name:        testIdLocalDbVlvIndex,
	}
	updatedResourceModel := localDbVlvIndexTestModel{
		backendName: testBackendNametest,
		baseDn:      "dc=example,dc=com",
		scope:       "base-object",
		filter:      "uid=user.2",
		sortOrder:   "mail",
		name:        testIdLocalDbVlvIndex,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckLocalDbVlvIndexDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccLocalDbVlvIndexResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedLocalDbVlvIndexAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccLocalDbVlvIndexResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedLocalDbVlvIndexAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccLocalDbVlvIndexResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_local_db_vlv_index." + resourceName,
				ImportStateId:     updatedResourceModel.backendName + "/" + updatedResourceModel.name,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccLocalDbVlvIndexResource(resourceName string, resourceModel localDbVlvIndexTestModel) string {
	return fmt.Sprintf(`

resource "pingdirectory_local_db_backend" "%[2]s" {
  backend_id            = "%[2]s"
  base_dn               = ["dc=example1,dc=com"]
  writability_mode      = "enabled"
  db_directory          = "db"
  import_temp_directory = "tmp"
  enabled               = true
}

resource "pingdirectory_local_db_vlv_index" "%[1]s" {
  backend_name = pingdirectory_local_db_backend.%[2]s.backend_id
  base_dn      = "%[3]s"
  scope        = "%[4]s"
  filter       = "%[5]s"
  sort_order   = "%[6]s"
  name         = "%[7]s"
}`, resourceName,
		resourceModel.backendName,
		resourceModel.baseDn,
		resourceModel.scope,
		resourceModel.filter,
		resourceModel.sortOrder,
		resourceModel.name)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedLocalDbVlvIndexAttributes(config localDbVlvIndexTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LocalDbVlvIndexApi.GetLocalDbVlvIndex(ctx, config.name, config.backendName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Local Db Vlv Index"
		err = acctest.TestAttributesMatchString(resourceType, &config.name, "base-dn",
			config.baseDn, response.BaseDN)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.name, "scope",
			config.scope, response.Scope.String())
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.name, "filter",
			config.filter, response.Filter)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.name, "sort-order",
			config.sortOrder, response.SortOrder)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.name, "name",
			config.name, response.Name)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckLocalDbVlvIndexDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.LocalDbVlvIndexApi.GetLocalDbVlvIndex(ctx, testIdLocalDbVlvIndex, testBackendNametest).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Local Db Vlv Index", testIdLocalDbVlvIndex)
	}
	return nil
}
