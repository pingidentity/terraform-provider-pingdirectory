package connectioncriteria_test

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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedSimpleConnectionCriteriaAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_connection_criteria.%s", resourceName), "description", initialResourceModel.description),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_connection_criteria.%s", resourceName), "user_auth_type.*", initialResourceModel.user_auth_type[0]),
					resource.TestCheckResourceAttrSet("data.pingdirectory_connection_criteria_list.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields.
				Config: testAccSimpleConnectionCriteriaResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSimpleConnectionCriteriaAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccSimpleConnectionCriteriaResource(resourceName, initialResourceModel),
				ResourceName:            "pingdirectory_connection_criteria." + resourceName,
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
resource "pingdirectory_connection_criteria" "%[1]s" {
  type           = "simple"
  id             = "%[2]s"
  description    = "%[3]s"
  user_auth_type = %[4]s
}

data "pingdirectory_connection_criteria" "%[1]s" {
  id = "%[2]s"
  depends_on = [
    pingdirectory_connection_criteria.%[1]s
  ]
}

data "pingdirectory_connection_criteria_list" "list" {
  depends_on = [
    pingdirectory_connection_criteria.%[1]s
  ]
}`, resourceName, resourceModel.id, resourceModel.description, acctest.StringSliceToTerraformString(resourceModel.user_auth_type))
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSimpleConnectionCriteriaAttributes(config simpleConnectionCriteriaTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "connection criteria"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ConnectionCriteriaApi.GetConnectionCriteria(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "description",
			config.description, response.SimpleConnectionCriteriaResponse.Description)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "user_auth_type",
			config.user_auth_type, client.StringSliceEnumconnectionCriteriaUserAuthTypeProp(response.SimpleConnectionCriteriaResponse.UserAuthType))
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
