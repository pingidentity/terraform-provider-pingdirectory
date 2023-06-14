package requestcriteria_test

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

const testIdRequestCriteria = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type requestCriteriaTestModel struct {
	id string
}

func TestAccRequestCriteria(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := requestCriteriaTestModel{
		id: testIdRequestCriteria,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckRequestCriteriaDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccRequestCriteriaResource(resourceName, initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccRequestCriteriaResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_request_criteria." + resourceName,
				ImportStateId:     initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccRequestCriteriaResource(resourceName string, resourceModel requestCriteriaTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_request_criteria" "%[1]s" {
  type = "root-dse"
	 id = "%[2]s"
}`, resourceName,
		resourceModel.id)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedRequestCriteriaAttributes(config requestCriteriaTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		_, _, err := testClient.RequestCriteriaApi.GetRequestCriteria(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckRequestCriteriaDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.RequestCriteriaApi.GetRequestCriteria(ctx, testIdRequestCriteria).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Request Criteria", testIdRequestCriteria)
	}
	return nil
}
