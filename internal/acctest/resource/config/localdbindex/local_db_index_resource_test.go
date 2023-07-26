package localdbindex_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdLocalDbIndex = "dc"
const testBackendName = "userRoot"

// Attributes to test with. Add optional properties to test here if desired.
type localDbIndexTestModel struct {
	backendName string
	attribute   string
	indexType   []string
}

func TestAccLocalDbIndex(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := localDbIndexTestModel{
		backendName: testBackendName,
		attribute:   testIdLocalDbIndex,
		indexType:   []string{"equality"},
	}
	updatedResourceModel := localDbIndexTestModel{
		backendName: testBackendName,
		attribute:   testIdLocalDbIndex,
		indexType:   []string{"substring"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckLocalDbIndexDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccLocalDbIndexResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedLocalDbIndexAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_local_db_index.%s", resourceName), "attribute", initialResourceModel.attribute),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_local_db_index.%s", resourceName), "index_type.*", initialResourceModel.indexType[0]),
				),
			},
			{
				// Test updating some fields
				Config: testAccLocalDbIndexResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedLocalDbIndexAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccLocalDbIndexResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_local_db_index." + resourceName,
				ImportStateId:     updatedResourceModel.backendName + "/" + updatedResourceModel.attribute,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccLocalDbIndexResource(resourceName string, resourceModel localDbIndexTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_local_db_index" "%[1]s" {
  backend_name = "%[2]s"
  attribute    = "%[3]s"
  index_type   = %[4]s
}

data "pingdirectory_local_db_index" "%[1]s" {
	 backend_name = "%[2]s"
	 attribute = "%[3]s"
  depends_on = [
    pingdirectory_local_db_index.%[1]s
  ]
}`, resourceName,
		resourceModel.backendName,
		resourceModel.attribute,
		acctest.StringSliceToTerraformString(resourceModel.indexType))
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedLocalDbIndexAttributes(config localDbIndexTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LocalDbIndexApi.GetLocalDbIndex(ctx, config.attribute, config.backendName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Local Db Index"
		err = acctest.TestAttributesMatchString(resourceType, &config.attribute, "attribute",
			config.attribute, response.Attribute)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.attribute, "index-type",
			config.indexType, client.StringSliceEnumlocalDbIndexIndexTypeProp(response.IndexType))
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckLocalDbIndexDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.LocalDbIndexApi.GetLocalDbIndex(ctx, testIdLocalDbIndex, testBackendName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Local Db Index", testIdLocalDbIndex)
	}
	return nil
}
