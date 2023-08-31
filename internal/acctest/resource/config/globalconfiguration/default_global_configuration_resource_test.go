package globalconfiguration_test

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

// Some global configuration attributes to test with
type testModel struct {
	encryptData         bool
	sensitiveAttribute  []string
	resultCodeMap       string
	sizeLimit           int64
	maximumShutdownTime string
}

func TestAccGlobalConfiguration(t *testing.T) {
	resourceName := "global"
	initialResourceModel := testModel{
		encryptData:         false,
		sensitiveAttribute:  []string{"Delivered One-Time Password", "TOTP Shared Secret"},
		resultCodeMap:       "Sun DS Compatible Behavior",
		sizeLimit:           2000,
		maximumShutdownTime: "4 m",
	}
	updatedResourceModel := testModel{
		encryptData:         true,
		sensitiveAttribute:  []string{"TOTP Shared Secret"},
		resultCodeMap:       "",
		sizeLimit:           1000,
		maximumShutdownTime: "3 m",
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				Config: testAccGlobalConfigurationResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedGlobalConfigurationAttributes(initialResourceModel),
					// Check some computed attributes are set as expected (PingDirectory defaults)
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_default_global_configuration.%s", resourceName), "encrypt_backups_by_default", "true"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_default_global_configuration.%s", resourceName), "default_password_policy", "Default Password Policy"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_default_global_configuration.%s", resourceName), "ldap_join_size_limit", "10000"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_default_global_configuration.%s", resourceName), "replication_set_name", ""),
					// Check those values are visible on the data source
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_global_configuration.%s", resourceName), "encrypt_backups_by_default", "true"),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_global_configuration.%s", resourceName), "default_password_policy", "Default Password Policy"),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_global_configuration.%s", resourceName), "ldap_join_size_limit", "10000"),
				),
			},
			{
				// Test updating some fields
				Config: testAccGlobalConfigurationResource(resourceName, updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedGlobalConfigurationAttributes(updatedResourceModel),
				),
			},
			{
				// Test returning to the original configuration
				Config: testAccGlobalConfigurationResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedGlobalConfigurationAttributes(initialResourceModel),
				),
			},
			{
				// Test importing the global configuration
				Config:       testAccGlobalConfigurationResource(resourceName, initialResourceModel),
				ResourceName: "pingdirectory_default_global_configuration." + resourceName,
				// The id doesn't matter for singleton config objects
				ImportStateId:     resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGlobalConfigurationResource(resourceName string, resourceModel testModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_global_configuration" "%[1]s" {
  encrypt_data          = %[2]t
  sensitive_attribute   = %[3]s
  result_code_map       = "%[4]s"
  size_limit            = %[5]d
  maximum_shutdown_time = "%[6]s"
}

data "pingdirectory_global_configuration" "%[1]s" {
  depends_on = [
    pingdirectory_default_global_configuration.%[1]s
  ]
}`, resourceName, resourceModel.encryptData,
		acctest.StringSliceToTerraformString(resourceModel.sensitiveAttribute),
		resourceModel.resultCodeMap, resourceModel.sizeLimit, resourceModel.maximumShutdownTime)
}

// Test that the expected global configuration attributes are set on the PingDirectory server
func testAccCheckExpectedGlobalConfigurationAttributes(globalConfig testModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "global configuration"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.GlobalConfigurationApi.GetGlobalConfiguration(ctx).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, nil, "encrypt-data", globalConfig.encryptData, *response.EncryptData)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, nil, "sensitive-attribute", globalConfig.sensitiveAttribute, response.SensitiveAttribute)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "result-code-map", globalConfig.resultCodeMap, response.ResultCodeMap)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, nil, "size-limit", globalConfig.sizeLimit, *response.SizeLimit)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "maximum-shutdown-time", globalConfig.maximumShutdownTime, response.MaximumShutdownTime)
		if err != nil {
			return err
		}
		return nil
	}
}
