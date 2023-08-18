package cryptomanager_test

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

// Attributes to test with. Add optional properties to test here if desired.
type cryptoManagerTestModel struct {
	mac_key_length    int64
	cipher_key_length int64
	ssl_cert_nickname string
}

func TestAccCryptoManager(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := cryptoManagerTestModel{
		mac_key_length:    256,
		cipher_key_length: 256,
		ssl_cert_nickname: "ssl-certificate-alias",
	}
	// ads-certificate is the default value for the ssl_cert_nickname attribute
	updatedResourceModel := cryptoManagerTestModel{
		mac_key_length:    192,
		cipher_key_length: 192,
		ssl_cert_nickname: "ads-certificate",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccCryptoManagerResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedCryptoManagerAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_crypto_manager.%s", resourceName), "mac_key_length", strconv.FormatInt(initialResourceModel.mac_key_length, 10)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_crypto_manager.%s", resourceName), "cipher_key_length", strconv.FormatInt(initialResourceModel.cipher_key_length, 10)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_crypto_manager.%s", resourceName), "ssl_cert_nickname", initialResourceModel.ssl_cert_nickname),
				),
			},
			{
				// Test updating some fields
				Config: testAccCryptoManagerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedCryptoManagerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccCryptoManagerResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_default_crypto_manager." + resourceName,
				ImportStateId:     resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccCryptoManagerResource(resourceName string, resourceModel cryptoManagerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_crypto_manager" "%[1]s" {
  mac_key_length    = %[2]d
  cipher_key_length = %[3]d
  ssl_cert_nickname = "%[4]s"
}

data "pingdirectory_crypto_manager" "%[1]s" {
  depends_on = [
    pingdirectory_default_crypto_manager.%[1]s
  ]
}`, resourceName,
		resourceModel.mac_key_length,
		resourceModel.cipher_key_length,
		resourceModel.ssl_cert_nickname)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedCryptoManagerAttributes(config cryptoManagerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "Crypto Manager"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.CryptoManagerApi.GetCryptoManager(ctx).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchInt(resourceType, nil, "mac-key-length",
			config.mac_key_length, *response.MacKeyLength)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, nil, "cipher-key-length",
			config.cipher_key_length, *response.CipherKeyLength)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "ssl-cert-nickname",
			config.ssl_cert_nickname, response.SslCertNickname)
		if err != nil {
			return err
		}

		return nil
	}
}
