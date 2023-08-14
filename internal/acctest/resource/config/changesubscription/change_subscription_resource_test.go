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
				Check:  resource.TestCheckResourceAttrSet("data.pingdirectory_change_subscriptions.list", "ids.0"),
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
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.ChangeSubscriptionApi.DeleteChangeSubscription(ctx, initialResourceModel.id).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccChangeSubscriptionResource(resourceName string, resourceModel changeSubscriptionTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_change_subscription" "%[1]s" {
  name = "%[2]s"
}

data "pingdirectory_change_subscription" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_change_subscription.%[1]s
  ]
}

data "pingdirectory_change_subscriptions" "list" {
  depends_on = [
    pingdirectory_change_subscription.%[1]s
  ]
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
