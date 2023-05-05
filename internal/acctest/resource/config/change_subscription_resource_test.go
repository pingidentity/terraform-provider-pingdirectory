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

const testIdChangeSubscription = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type changeSubscriptionTestModel struct {
	id string
}

func TestAccChangeSubscription(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := changeSubscriptionTestModel{
		id: testIdChangeSubscription,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckChangeSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccChangeSubscriptionResource(resourceName, initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccChangeSubscriptionResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_change_subscription." + resourceName,
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

func testAccChangeSubscriptionResource(resourceName string, resourceModel changeSubscriptionTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_change_subscription" "%[1]s" {
  id = "%[2]s"
}`, resourceName,
		resourceModel.id)
}

// Test that any objects created by the test are destroyed
func testAccCheckChangeSubscriptionDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ChangeSubscriptionApi.GetChangeSubscription(ctx, testIdChangeSubscription).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Change Subscription", testIdChangeSubscription)
	}
	return nil
}
