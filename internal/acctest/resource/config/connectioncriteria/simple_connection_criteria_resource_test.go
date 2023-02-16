package connectioncriteria_test

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

const testIdSimpleConnectionCriteria = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type simpleConnectionCriteriaTestModel struct {
	id             string
	description    string
	user_auth_type []string
}

func TestAccSimpleConnectionCriteria(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := simpleConnectionCriteriaTestModel{
		id:             testIdSimpleConnectionCriteria,
		description:    "Test simple connection example",
		user_auth_type: []string{"internal", "sasl"},
	}

	updatedResourceModel := simpleConnectionCriteriaTestModel{
		id:             testIdSimpleConnectionCriteria,
		description:    "Test simple connection modified",
		user_auth_type: []string{"internal", "sasl", "simple"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckSimpleConnectionCriteriaDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSimpleConnectionCriteriaResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSimpleConnectionCriteriaAttributes(initialResourceModel),
			},
			{
				// Test updating some fields.
				Config: testAccSimpleConnectionCriteriaResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSimpleConnectionCriteriaAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccSimpleConnectionCriteriaResource(resourceName, initialResourceModel),
				ResourceName:            "pingdirectory_simple_connection_criteria." + resourceName,
				ImportStateId:           initialResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccSimpleConnectionCriteriaResource(resourceName string, resourceModel simpleConnectionCriteriaTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_simple_connection_criteria" "%[1]s" {
	 id = "%[2]s"
	 description  = "%[3]s"
	 user_auth_type = %[4]s
}`, resourceName, resourceModel.id, resourceModel.description, acctest.StringSliceToTerraformString(resourceModel.user_auth_type))
}

//acctest.StringSliceToTerraformString(permissionsList)

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSimpleConnectionCriteriaAttributes(config simpleConnectionCriteriaTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		_, _, err := testClient.ConnectionCriteriaApi.GetConnectionCriteria(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSimpleConnectionCriteriaDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ConnectionCriteriaApi.GetConnectionCriteria(ctx, testIdSimpleConnectionCriteria).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Simple Connection Criteria", testIdSimpleConnectionCriteria)
	}
	return nil
}
