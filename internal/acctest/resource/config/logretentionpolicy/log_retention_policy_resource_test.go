package logretentionpolicy_test

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

const testIdTimeLimitLogRetentionPolicy = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type timeLimitLogRetentionPolicyTestModel struct {
	id             string
	retainDuration string
}

func TestAccTimeLimitLogRetentionPolicy(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := timeLimitLogRetentionPolicyTestModel{
		id:             testIdTimeLimitLogRetentionPolicy,
		retainDuration: "3 d",
	}
	updatedResourceModel := timeLimitLogRetentionPolicyTestModel{
		id:             testIdTimeLimitLogRetentionPolicy,
		retainDuration: "1 w",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckTimeLimitLogRetentionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccTimeLimitLogRetentionPolicyResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedTimeLimitLogRetentionPolicyAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_log_retention_policy.%s", resourceName), "retain_duration", initialResourceModel.retainDuration),
					resource.TestCheckResourceAttrSet("data.pingdirectory_log_retention_policies.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccTimeLimitLogRetentionPolicyResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedTimeLimitLogRetentionPolicyAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccTimeLimitLogRetentionPolicyResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_log_retention_policy." + resourceName,
				ImportStateId:     updatedResourceModel.id,
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
					_, err := testClient.LogRetentionPolicyApi.DeleteLogRetentionPolicy(ctx, updatedResourceModel.id).Execute()
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

func testAccTimeLimitLogRetentionPolicyResource(resourceName string, resourceModel timeLimitLogRetentionPolicyTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_log_retention_policy" "%[1]s" {
  type            = "time-limit"
  name            = "%[2]s"
  retain_duration = "%[3]s"
}

data "pingdirectory_log_retention_policy" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_log_retention_policy.%[1]s
  ]
}

data "pingdirectory_log_retention_policies" "list" {
  depends_on = [
    pingdirectory_log_retention_policy.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.retainDuration)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedTimeLimitLogRetentionPolicyAttributes(config timeLimitLogRetentionPolicyTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LogRetentionPolicyApi.GetLogRetentionPolicy(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Time Limit Log Retention Policy"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "retain-duration",
			config.retainDuration, response.TimeLimitLogRetentionPolicyResponse.RetainDuration)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckTimeLimitLogRetentionPolicyDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.LogRetentionPolicyApi.GetLogRetentionPolicy(ctx, testIdTimeLimitLogRetentionPolicy).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Time Limit Log Retention Policy", testIdTimeLimitLogRetentionPolicy)
	}
	return nil
}
