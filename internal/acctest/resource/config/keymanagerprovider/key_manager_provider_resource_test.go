package keymanagerprovider_test

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

const testIdFileBasedKeyManagerProvider = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type fileBasedKeyManagerProviderTestModel struct {
	id           string
	keyStoreFile string
	enabled      bool
	description  string
}

func TestAccFileBasedKeyManagerProvider(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := fileBasedKeyManagerProviderTestModel{
		id:           testIdFileBasedKeyManagerProvider,
		keyStoreFile: "/tmp/initial-key-store-file",
		enabled:      false,
		description:  "Initial resource model description",
	}
	updatedResourceModel := fileBasedKeyManagerProviderTestModel{
		id:           testIdFileBasedKeyManagerProvider,
		keyStoreFile: "/tmp/updated-key-store-file",
		enabled:      false,
		description:  "Updated resource model description",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckFileBasedKeyManagerProviderDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccFileBasedKeyManagerProviderResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedFileBasedKeyManagerProviderAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_key_manager_provider.%s", resourceName), "key_store_file", initialResourceModel.keyStoreFile),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_key_manager_provider.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
				),
			},
			{
				// Test updating some fields
				Config: testAccFileBasedKeyManagerProviderResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedFileBasedKeyManagerProviderAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccFileBasedKeyManagerProviderResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_key_manager_provider." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
					"key_store_pin",
					"private_key_pin",
				},
			},
		},
	})
}

func testAccFileBasedKeyManagerProviderResource(resourceName string, resourceModel fileBasedKeyManagerProviderTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_key_manager_provider" "%[1]s" {
  type           = "file-based"
  id             = "%[2]s"
  key_store_file = "%[3]s"
  enabled        = %[4]t
  description    = "%[5]s"
}

data "pingdirectory_key_manager_provider" "%[1]s" {
  id = "%[2]s"
  depends_on = [
    pingdirectory_key_manager_provider.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.keyStoreFile,
		resourceModel.enabled,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedFileBasedKeyManagerProviderAttributes(config fileBasedKeyManagerProviderTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.KeyManagerProviderApi.GetKeyManagerProvider(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "File Based Key Manager Provider"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "key-store-file",
			config.keyStoreFile, response.FileBasedKeyManagerProviderResponse.KeyStoreFile)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.FileBasedKeyManagerProviderResponse.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "description",
			config.description, *response.FileBasedKeyManagerProviderResponse.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckFileBasedKeyManagerProviderDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.KeyManagerProviderApi.GetKeyManagerProvider(ctx, testIdFileBasedKeyManagerProvider).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("File Based Key Manager Provider", testIdFileBasedKeyManagerProvider)
	}
	return nil
}
