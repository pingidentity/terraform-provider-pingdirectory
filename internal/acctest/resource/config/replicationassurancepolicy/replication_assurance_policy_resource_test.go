package replicationassurancepolicy_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdReplicationAssurancePolicy = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type replicationAssurancePolicyTestModel struct {
	id                   string
	description          string
	evaluationOrderIndex int64
	timeout              string
}

func TestAccReplicationAssurancePolicy(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := replicationAssurancePolicyTestModel{
		id:                   testIdReplicationAssurancePolicy,
		description:          "Initial replication assurance policy",
		evaluationOrderIndex: 3,
		timeout:              "3 s",
	}
	updatedResourceModel := replicationAssurancePolicyTestModel{
		id:                   testIdReplicationAssurancePolicy,
		description:          "Updated replication assurance policy",
		evaluationOrderIndex: 4,
		timeout:              "20 ms",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckReplicationAssurancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccReplicationAssurancePolicyResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedReplicationAssurancePolicyAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_replication_assurance_policy.%s", resourceName), "evaluation_order_index", strconv.FormatInt(initialResourceModel.evaluationOrderIndex, 10)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_replication_assurance_policy.%s", resourceName), "timeout", initialResourceModel.timeout),
					resource.TestCheckResourceAttrSet("data.pingdirectory_replication_assurance_policies.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccReplicationAssurancePolicyResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedReplicationAssurancePolicyAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccReplicationAssurancePolicyResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_replication_assurance_policy." + resourceName,
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
					_, err := testClient.ReplicationAssurancePolicyApi.DeleteReplicationAssurancePolicy(ctx, updatedResourceModel.id).Execute()
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

func testAccReplicationAssurancePolicyResource(resourceName string, resourceModel replicationAssurancePolicyTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_replication_assurance_policy" "%[1]s" {
  name                   = "%[2]s"
  description            = "%[3]s"
  evaluation_order_index = %[4]d
  timeout                = "%[5]s"
}

data "pingdirectory_replication_assurance_policy" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_replication_assurance_policy.%[1]s
  ]
}

data "pingdirectory_replication_assurance_policies" "list" {
  depends_on = [
    pingdirectory_replication_assurance_policy.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.description,
		resourceModel.evaluationOrderIndex,
		resourceModel.timeout)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedReplicationAssurancePolicyAttributes(config replicationAssurancePolicyTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ReplicationAssurancePolicyApi.GetReplicationAssurancePolicy(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Replication Assurance Policy"
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "description",
			config.description, response.Description)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "evaluation-order-index",
			config.evaluationOrderIndex, response.EvaluationOrderIndex)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "timeout",
			config.timeout, response.Timeout)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckReplicationAssurancePolicyDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ReplicationAssurancePolicyApi.GetReplicationAssurancePolicy(ctx, testIdReplicationAssurancePolicy).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Replication Assurance Policy", testIdReplicationAssurancePolicy)
	}
	return nil
}
