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
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

const tmpName = "mytrustmanagerprovider"
const resourceName = "TestTMP"
const importResourceName = "ImportResource"
const importedThirdPartyTMP = "MyThirdPartyTrustManagerProvider"

func TestAccBlindTrustManagerProvider(t *testing.T) {
	importId := "Blind Trust"
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
				Config:        testAccBlindTrustManagerProviderResourceEmpty(importResourceName, importId),
				ResourceName:  "pingdirectory_blind_trust_manager_provider." + importResourceName,
				ImportStateId: importId,
				ImportState:   true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_blind_trust_manager_provider.%s", importResourceName), "enabled", "false"),
				),
			},
		},
	})
}

func TestAccFileBasedTrustManagerProvider(t *testing.T) {
	importId := "JKS"
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
				Config:        testAccFileBasedManagerProviderResourceEmpty(importResourceName, importId),
				ResourceName:  "pingdirectory_file_based_trust_manager_provider." + importResourceName,
				ImportStateId: importId,
				ImportState:   true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_file_based_trust_manager_provider.%s", importResourceName), "enabled", "true"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_file_based_trust_manager_provider.%s", importResourceName), "trust-store-file", "config/keystore"),
				),
			},
		},
	})
}

func TestAccJvmDefaultTrustManagerProvider(t *testing.T) {
	importId := "JVM-Default"
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
				Config:        testAccJvmDefaultManagerProviderResourceEmpty(importResourceName, importId),
				ResourceName:  "pingdirectory_jvm_default_trust_manager_provider." + importResourceName,
				ImportStateId: importId,
				ImportState:   true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_jvm_default_trust_manager_provider.%s", importResourceName), "enabled", "false"),
				),
			},
		},
	})
}

func TestAccThirdPartyTrustManagerProvider(t *testing.T) {
	initialArguments := []string{"val1=one", "val2=two"}
	updatedArguments := []string{"val3=three"}
	extensionClass := "com.unboundid.directory.sdk.common.api.TrustManagerProvider"
	importArguments := []string{"val1=one", "val2=two", "asdf=jkl;"}
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
				// First, create a third party trust manager provider to import
				PreConfig: func() {
					testAccCreateThirdPartyProviderForImport(t, importedThirdPartyTMP, false, extensionClass, importArguments)
				},
				Config:        testAccThirdPartyTrustManagerProviderResource(importResourceName, importedThirdPartyTMP, false, extensionClass, importArguments),
				ResourceName:  "pingdirectory_third_party_trust_manager_provider." + importResourceName,
				ImportStateId: importedThirdPartyTMP,
				ImportState:   true,
				// Check for expected state values
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_third_party_trust_manager_provider.%s", importResourceName),
						"enabled", "false"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_third_party_trust_manager_provider.%s", importResourceName),
						"extension_class", extensionClass),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("pingdirectory_third_party_trust_manager_provider.%s", importResourceName),
						"extension_argument.*", importArguments[0]),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("pingdirectory_third_party_trust_manager_provider.%s", importResourceName),
						"extension_argument.*", importArguments[1]),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("pingdirectory_third_party_trust_manager_provider.%s", importResourceName),
						"extension_argument.*", importArguments[2]),
				),
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

func testAccBlindTrustManagerProviderResourceEmpty(resourceName, providerName string) string {
	return fmt.Sprintf(`
resource "pingdirectory_blind_trust_manager_provider" "%[1]s" {
	id = "%[2]s"
}`, resourceName, providerName)
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

func testAccFileBasedManagerProviderResourceEmpty(resourceName, providerName string) string {
	return fmt.Sprintf(`
resource "pingdirectory_file_based_trust_manager_provider" "%[1]s" {
	id = "%[2]s"
}`, resourceName, providerName)
}

func testAccJvmDefaultTrustManagerProviderResource(resourceName, providerName string, enabled bool) string {
	return fmt.Sprintf(`
resource "pingdirectory_jvm_default_trust_manager_provider" "%[1]s" {
	id = "%[2]s"
	enabled = %[3]t
}`, resourceName, providerName, enabled)
}

func testAccJvmDefaultManagerProviderResourceEmpty(resourceName, providerName string) string {
	return fmt.Sprintf(`
resource "pingdirectory_jvm_default_trust_manager_provider" "%[1]s" {
	id = "%[2]s"
}`, resourceName, providerName)
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
	// Attempt to destroy the third party trust manager provider created by this test
	// Ignore any error - just make best attempt here
	testAccDestroyThirdPartyProviderForImport(importedThirdPartyTMP)

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

// Create a third party provider on the server to be imported
func testAccCreateThirdPartyProviderForImport(t *testing.T, name string, enabled bool, extensionClass string, arguments []string) {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	addRequest := client.NewAddThirdPartyTrustManagerProviderRequest(name,
		[]client.EnumthirdPartyTrustManagerProviderSchemaUrn{client.ENUMTHIRDPARTYTRUSTMANAGERPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TRUST_MANAGER_PROVIDERTHIRD_PARTY},
		extensionClass,
		enabled)
	addRequest.ExtensionArgument = arguments
	apiAddRequest := testClient.TrustManagerProviderApi.AddTrustManagerProvider(ctx)
	apiAddRequest = apiAddRequest.AddTrustManagerProviderRequest(
		client.AddThirdPartyTrustManagerProviderRequestAsAddTrustManagerProviderRequest(addRequest))
	_, _, err := testClient.TrustManagerProviderApi.AddTrustManagerProviderExecute(apiAddRequest)
	if err != nil {
		t.Error(err)
	}
}

// Destroy the third party provider created for import
func testAccDestroyThirdPartyProviderForImport(name string) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.TrustManagerProviderApi.DeleteTrustManagerProviderExecute(
		testClient.TrustManagerProviderApi.DeleteTrustManagerProvider(ctx, name))
	if err != nil {
		return err
	}
	return nil
}
