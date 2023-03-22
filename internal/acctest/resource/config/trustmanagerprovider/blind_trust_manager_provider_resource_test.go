package trustmanagerprovider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdBlindTrustManagerProvider = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type blindTrustManagerProviderTestModel struct {
	id      string
	enabled bool
}

func TestAccBlindTrustManagerProvider(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := blindTrustManagerProviderTestModel{
		id:      testIdBlindTrustManagerProvider,
		enabled: true,
	}
	updatedResourceModel := blindTrustManagerProviderTestModel{
		id:      testIdBlindTrustManagerProvider,
		enabled: false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckBlindTrustManagerProviderDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccBlindTrustManagerProviderResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedBlindTrustManagerProviderAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccBlindTrustManagerProviderResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedBlindTrustManagerProviderAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccBlindTrustManagerProviderResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_blind_trust_manager_provider." + resourceName,
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

func testAccBlindTrustManagerProviderResource(resourceName string, resourceModel blindTrustManagerProviderTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_blind_trust_manager_provider" "%[1]s" {
  id      = "%[2]s"
  enabled = %[3]t
}`, resourceName,
		resourceModel.id,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedBlindTrustManagerProviderAttributes(config blindTrustManagerProviderTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.TrustManagerProviderApi.GetTrustManagerProvider(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Blind Trust Manager Provider"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.BlindTrustManagerProviderResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckBlindTrustManagerProviderDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.TrustManagerProviderApi.GetTrustManagerProvider(ctx, testIdBlindTrustManagerProvider).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Blind Trust Manager Provider", testIdBlindTrustManagerProvider)
	}
	return nil
}
