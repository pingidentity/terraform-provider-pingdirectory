package localdbcompositeindex_test

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

const testIdLocalDbCompositeIndex = "MyId"
const testBackendNameComposite = "userRoot"

// Attributes to test with. Add optional properties to test here if desired.
type localDbCompositeIndexTestModel struct {
	id                 string
	backendName        string
	description        string
	indexFilterPattern string
}

func TestAccLocalDbCompositeIndex(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := localDbCompositeIndexTestModel{
		id:                 testIdLocalDbCompositeIndex,
		backendName:        testBackendNameComposite,
		description:        "initial resource model description",
		indexFilterPattern: "(sn=?)",
	}
	// indexFilterPattern cannot be modified after creation
	updatedResourceModel := localDbCompositeIndexTestModel{
		id:                 testIdLocalDbCompositeIndex,
		backendName:        testBackendNameComposite,
		description:        "updated resource model description",
		indexFilterPattern: "(sn=?)",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckLocalDbCompositeIndexDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccLocalDbCompositeIndexResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedLocalDbCompositeIndexAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_local_db_composite_index.%s", resourceName), "index_filter_pattern", initialResourceModel.indexFilterPattern),
				),
			},
			{
				// Test updating some fields
				Config: testAccLocalDbCompositeIndexResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedLocalDbCompositeIndexAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccLocalDbCompositeIndexResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_local_db_composite_index." + resourceName,
				ImportStateId:     updatedResourceModel.backendName + "/" + updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccLocalDbCompositeIndexResource(resourceName string, resourceModel localDbCompositeIndexTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_local_db_composite_index" "%[1]s" {
  id                   = "%[2]s"
  backend_name         = "%[3]s"
  description          = "%[4]s"
  index_filter_pattern = "%[5]s"
}

data "pingdirectory_local_db_composite_index" "%[1]s" {
  id           = "%[2]s"
  backend_name = "%[3]s"
  depends_on = [
    pingdirectory_local_db_composite_index.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.backendName,
		resourceModel.description,
		resourceModel.indexFilterPattern)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedLocalDbCompositeIndexAttributes(config localDbCompositeIndexTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LocalDbCompositeIndexApi.GetLocalDbCompositeIndex(ctx, config.id, config.backendName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Local Db Composite Index"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "description",
			config.description, *response.Description)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "index-filter-pattern",
			config.indexFilterPattern, response.IndexFilterPattern)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckLocalDbCompositeIndexDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.LocalDbCompositeIndexApi.GetLocalDbCompositeIndex(ctx, testIdLocalDbCompositeIndex, testBackendNameComposite).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Local Db Composite Index", testIdLocalDbCompositeIndex)
	}
	return nil
}