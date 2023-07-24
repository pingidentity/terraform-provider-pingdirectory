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

const testIdClientConnectionPolicy = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type clientConnectionPolicyTestModel struct {
	policyId             string
	enabled              bool
	evaluationOrderIndex int64
}

func TestAccClientConnectionPolicy(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := clientConnectionPolicyTestModel{
		policyId:             testIdClientConnectionPolicy,
		enabled:              true,
		evaluationOrderIndex: 1,
	}
	updatedResourceModel := clientConnectionPolicyTestModel{
		policyId:             testIdClientConnectionPolicy,
		enabled:              false,
		evaluationOrderIndex: 2,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckClientConnectionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccClientConnectionPolicyResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedClientConnectionPolicyAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccClientConnectionPolicyResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedClientConnectionPolicyAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccClientConnectionPolicyResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_client_connection_policy." + resourceName,
				ImportStateId:     updatedResourceModel.policyId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccClientConnectionPolicyResource(resourceName string, resourceModel clientConnectionPolicyTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_client_connection_policy" "%[1]s" {
  policy_id              = "%[2]s"
  enabled                = %[3]t
  evaluation_order_index = %[4]d
}`, resourceName,
		resourceModel.policyId,
		resourceModel.enabled,
		resourceModel.evaluationOrderIndex)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedClientConnectionPolicyAttributes(config clientConnectionPolicyTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ClientConnectionPolicyApi.GetClientConnectionPolicy(ctx, config.policyId).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Client Connection Policy"
		err = acctest.TestAttributesMatchString(resourceType, &config.policyId, "policy-id",
			config.policyId, response.PolicyID)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.policyId, "enabled",
			config.enabled, response.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.policyId, "evaluation-order-index",
			config.evaluationOrderIndex, response.EvaluationOrderIndex)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckClientConnectionPolicyDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ClientConnectionPolicyApi.GetClientConnectionPolicy(ctx, testIdClientConnectionPolicy).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Client Connection Policy", testIdClientConnectionPolicy)
	}
	return nil
}
