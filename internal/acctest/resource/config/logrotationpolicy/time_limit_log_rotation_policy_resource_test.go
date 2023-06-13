package logrotationpolicy_test

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

const testIdTimeLimitLogRotationPolicy = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type timeLimitLogRotationPolicyTestModel struct {
	id               string
	rotationInterval string
}

func TestAccTimeLimitLogRotationPolicy(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := timeLimitLogRotationPolicyTestModel{
		id:               testIdTimeLimitLogRotationPolicy,
		rotationInterval: "1 w",
	}
	updatedResourceModel := timeLimitLogRotationPolicyTestModel{
		id:               testIdTimeLimitLogRotationPolicy,
		rotationInterval: "2 w",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckTimeLimitLogRotationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccTimeLimitLogRotationPolicyResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedTimeLimitLogRotationPolicyAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccTimeLimitLogRotationPolicyResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedTimeLimitLogRotationPolicyAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccTimeLimitLogRotationPolicyResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_log_rotation_policy." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccTimeLimitLogRotationPolicyResource(resourceName string, resourceModel timeLimitLogRotationPolicyTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_log_rotation_policy" "%[1]s" {
  type              = "time-limit"
  id                = "%[2]s"
  rotation_interval = "%[3]s"
}`, resourceName,
		resourceModel.id,
		resourceModel.rotationInterval)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedTimeLimitLogRotationPolicyAttributes(config timeLimitLogRotationPolicyTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LogRotationPolicyApi.GetLogRotationPolicy(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Time Limit Log Rotation Policy"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "rotation-interval",
			config.rotationInterval, response.TimeLimitLogRotationPolicyResponse.RotationInterval)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckTimeLimitLogRotationPolicyDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.LogRotationPolicyApi.GetLogRotationPolicy(ctx, testIdTimeLimitLogRotationPolicy).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Time Limit Log Rotation Policy", testIdTimeLimitLogRotationPolicy)
	}
	return nil
}
