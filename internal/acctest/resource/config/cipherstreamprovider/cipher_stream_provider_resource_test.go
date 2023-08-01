package cipherstreamprovider_test

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

const testIdWaitForPassphraseCipherStreamProvider = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type waitForPassphraseCipherStreamProviderTestModel struct {
	id      string
	enabled bool
}

func TestAccWaitForPassphraseCipherStreamProvider(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := waitForPassphraseCipherStreamProviderTestModel{
		id:      testIdWaitForPassphraseCipherStreamProvider,
		enabled: true,
	}
	updatedResourceModel := waitForPassphraseCipherStreamProviderTestModel{
		id:      testIdWaitForPassphraseCipherStreamProvider,
		enabled: false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckWaitForPassphraseCipherStreamProviderDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccWaitForPassphraseCipherStreamProviderResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedWaitForPassphraseCipherStreamProviderAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_cipher_stream_provider.%s", resourceName), "type", "wait-for-passphrase"),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_cipher_stream_provider.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_cipher_stream_providers.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccWaitForPassphraseCipherStreamProviderResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedWaitForPassphraseCipherStreamProviderAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccWaitForPassphraseCipherStreamProviderResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_cipher_stream_provider." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
					"aws_secret_access_key",
				},
			},
		},
	})
}

func testAccWaitForPassphraseCipherStreamProviderResource(resourceName string, resourceModel waitForPassphraseCipherStreamProviderTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_cipher_stream_provider" "%[1]s" {
  type    = "wait-for-passphrase"
  id      = "%[2]s"
  enabled = %[3]t
}

data "pingdirectory_cipher_stream_provider" "%[1]s" {
  id = "%[2]s"
  depends_on = [
    pingdirectory_cipher_stream_provider.%[1]s
  ]
}

data "pingdirectory_cipher_stream_providers" "list" {
  depends_on = [
    pingdirectory_cipher_stream_provider.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedWaitForPassphraseCipherStreamProviderAttributes(config waitForPassphraseCipherStreamProviderTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.CipherStreamProviderApi.GetCipherStreamProvider(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Wait for Passphrase Cipher Stream Provider"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.WaitForPassphraseCipherStreamProviderResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckWaitForPassphraseCipherStreamProviderDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.CipherStreamProviderApi.GetCipherStreamProvider(ctx, testIdWaitForPassphraseCipherStreamProvider).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Wait for Passphrase Cipher Stream Provider", testIdWaitForPassphraseCipherStreamProvider)
	}
	return nil
}
