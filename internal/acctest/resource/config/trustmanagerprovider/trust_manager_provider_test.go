package trustmanagerprovider_test

import (
	"fmt"
	"testing"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const tmpName = "mytrustmanagerprovider"
const resourceName = "TestTMP"

func TestAccBlindTrustManagerProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckTrustManagerProviderDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource
				Config: testAccBlindTrustManagerProviderResource(resourceName, tmpName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_blind_trust_manager_provider.%s", resourceName), "include_jvm_default_issuers", "false"),
					testAccCheckExpectedBlindTrustManagerProviderAttributes(true),
				),
			},
			{
				// Test updating the resource
				Config: testAccBlindTrustManagerProviderResource(resourceName, tmpName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedBlindTrustManagerProviderAttributes(false),
				),
			},
			{
				// Test importing the resource
				Config:                  testAccBlindTrustManagerProviderResource(resourceName, tmpName, false),
				ResourceName:            "pingdirectory_blind_trust_manager_provider." + resourceName,
				ImportStateId:           tmpName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func TestAccFileBasedTrustManagerProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckTrustManagerProviderDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource
				Config: testAccFileBasedTrustManagerProviderResource(resourceName, tmpName, false, "config/keystore", "JKS"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_file_based_trust_manager_provider.%s", resourceName), "include_jvm_default_issuers", "false"),
					testAccCheckExpectedFileBasedTrustManagerProviderAttributes(false, "config/keystore", "JKS"),
				),
			},
			{
				// Test updating the resource
				Config: testAccFileBasedTrustManagerProviderResource(resourceName, tmpName, false, "config/truststore", "PKCS12"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedFileBasedTrustManagerProviderAttributes(false, "config/truststore", "PKCS12"),
				),
			},
			{
				// Test importing the resource
				Config:                  testAccFileBasedTrustManagerProviderResource(resourceName, tmpName, false, "config/truststore", "PKCS12"),
				ResourceName:            "pingdirectory_file_based_trust_manager_provider." + resourceName,
				ImportStateId:           tmpName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func TestAccJvmDefaultTrustManagerProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckTrustManagerProviderDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource
				Config: testAccJvmDefaultTrustManagerProviderResource(resourceName, tmpName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedJvmDefaultTrustManagerProviderAttributes(false),
				),
			},
			{
				// Test updating the resource
				Config: testAccJvmDefaultTrustManagerProviderResource(resourceName, tmpName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedJvmDefaultTrustManagerProviderAttributes(true),
				),
			},
			{
				// Test importing the resource
				Config:                  testAccJvmDefaultTrustManagerProviderResource(resourceName, tmpName, true),
				ResourceName:            "pingdirectory_jvm_default_trust_manager_provider." + resourceName,
				ImportStateId:           tmpName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func TestAccThirdPartyTrustManagerProvider(t *testing.T) {
	initialArguments := []string{"val1=one", "val2=two"}
	updatedArguments := []string{"val3=three"}
	extensionClass := "com.unboundid.directory.sdk.common.api.TrustManagerProvider"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckTrustManagerProviderDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource
				Config: testAccThirdPartyTrustManagerProviderResource(resourceName, tmpName, false,
					extensionClass, initialArguments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_third_party_trust_manager_provider.%s", resourceName), "include_jvm_default_issuers", "false"),
					testAccCheckExpectedThirdPartyTrustManagerProviderAttributes(false, extensionClass, initialArguments),
				),
			},
			{
				// Test updating the resource
				Config: testAccThirdPartyTrustManagerProviderResource(resourceName, tmpName, false,
					extensionClass, updatedArguments),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedThirdPartyTrustManagerProviderAttributes(false, extensionClass, updatedArguments),
				),
			},
			{
				// Test importing the resource
				Config: testAccThirdPartyTrustManagerProviderResource(resourceName, tmpName, false,
					extensionClass, updatedArguments),
				ResourceName:            "pingdirectory_third_party_trust_manager_provider." + resourceName,
				ImportStateId:           tmpName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccBlindTrustManagerProviderResource(resourceName, providerName string, enabled bool) string {
	return fmt.Sprintf(`
resource "pingdirectory_blind_trust_manager_provider" "%[1]s" {
	id = "%[2]s"
	enabled = %[3]t
}`, resourceName, providerName, enabled)
}

func testAccFileBasedTrustManagerProviderResource(resourceName, providerName string, enabled bool, trustStoreFile, trustStoreType string) string {
	return fmt.Sprintf(`
resource "pingdirectory_file_based_trust_manager_provider" "%[1]s" {
	id = "%[2]s"
	enabled = %[3]t
	trust_store_file = "%[4]s"
	trust_store_type = "%[5]s"
}`, resourceName, providerName, enabled, trustStoreFile, trustStoreType)
}

func testAccJvmDefaultTrustManagerProviderResource(resourceName, providerName string, enabled bool) string {
	return fmt.Sprintf(`
resource "pingdirectory_jvm_default_trust_manager_provider" "%[1]s" {
	id = "%[2]s"
	enabled = %[3]t
}`, resourceName, providerName, enabled)
}

func testAccThirdPartyTrustManagerProviderResource(resourceName, providerName string, enabled bool, extensionClass string, extensionArgument []string) string {
	return fmt.Sprintf(`
resource "pingdirectory_third_party_trust_manager_provider" "%[1]s" {
	id = "%[2]s"
	enabled = %[3]t
	extension_class = "%[4]s"
	extension_argument = %[5]s
}`, resourceName, providerName, enabled, extensionClass, acctest.StringSliceToTerraformString(extensionArgument))
}

// Test that any resources created by the test are destroyed
func testAccCheckTrustManagerProviderDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.TrustManagerProviderApi.GetTrustManagerProvider(ctx, tmpName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("trust manager provider", tmpName)
	}
	return nil
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedBlindTrustManagerProviderAttributes(enabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "blind trust manager provider"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.TrustManagerProviderApi.GetTrustManagerProvider(ctx, tmpName).Execute()
		if err != nil {
			return err
		}
		name := tmpName
		err = acctest.TestAttributesMatchBool(resourceType, &name, "enabled", enabled, response.BlindTrustManagerProviderResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedFileBasedTrustManagerProviderAttributes(enabled bool, trustStoreFile, trustStoreType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "file based trust manager provider"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.TrustManagerProviderApi.GetTrustManagerProvider(ctx, tmpName).Execute()
		if err != nil {
			return err
		}
		name := tmpName
		err = acctest.TestAttributesMatchBool(resourceType, &name, "enabled", enabled, response.FileBasedTrustManagerProviderResponse.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &name, "trust-store-file", trustStoreFile, response.FileBasedTrustManagerProviderResponse.TrustStoreFile)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, &name, "trust-store-type", trustStoreType, response.FileBasedTrustManagerProviderResponse.TrustStoreType)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedJvmDefaultTrustManagerProviderAttributes(enabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "jvm default trust manager provider"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.TrustManagerProviderApi.GetTrustManagerProvider(ctx, tmpName).Execute()
		if err != nil {
			return err
		}
		name := tmpName
		err = acctest.TestAttributesMatchBool(resourceType, &name, "enabled", enabled, response.JvmDefaultTrustManagerProviderResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedThirdPartyTrustManagerProviderAttributes(enabled bool, extensionClass string, arguments []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "third party trust manager provider"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.TrustManagerProviderApi.GetTrustManagerProvider(ctx, tmpName).Execute()
		if err != nil {
			return err
		}
		name := tmpName
		err = acctest.TestAttributesMatchBool(resourceType, &name, "enabled", enabled, response.ThirdPartyTrustManagerProviderResponse.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &name, "extension-class", extensionClass, response.ThirdPartyTrustManagerProviderResponse.ExtensionClass)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &name, "extension-argument", arguments, response.ThirdPartyTrustManagerProviderResponse.ExtensionArgument)
		if err != nil {
			return err
		}
		return nil
	}
}
