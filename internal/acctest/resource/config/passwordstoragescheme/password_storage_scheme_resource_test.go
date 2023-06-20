package passwordstoragescheme_test

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

const testIdPasswordStorageScheme = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type passwordStorageSchemeTestModel struct {
	id                    string
	iterationCount        int64
	parallelismFactor     int64
	memoryUsageKb         int64
	saltLengthBytes       int64
	derivedKeyLengthBytes int64
	enabled               bool
}

func TestAccPasswordStorageScheme(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := passwordStorageSchemeTestModel{
		id:                    testIdPasswordStorageScheme,
		iterationCount:        10,
		parallelismFactor:     1,
		memoryUsageKb:         16,
		saltLengthBytes:       16,
		derivedKeyLengthBytes: 16,
		enabled:               true,
	}
	updatedResourceModel := passwordStorageSchemeTestModel{
		id:                    testIdPasswordStorageScheme,
		iterationCount:        20,
		parallelismFactor:     2,
		memoryUsageKb:         32,
		saltLengthBytes:       8,
		derivedKeyLengthBytes: 8,
		enabled:               false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckPasswordStorageSchemeDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccPasswordStorageSchemeResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedPasswordStorageSchemeAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccPasswordStorageSchemeResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedPasswordStorageSchemeAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccPasswordStorageSchemeResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_password_storage_scheme." + resourceName,
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

func testAccPasswordStorageSchemeResource(resourceName string, resourceModel passwordStorageSchemeTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_password_storage_scheme" "%[1]s" {
  type                     = "argon2d"
  id                       = "%[2]s"
  iteration_count          = %[3]d
  parallelism_factor       = %[4]d
  memory_usage_kb          = %[5]d
  salt_length_bytes        = %[6]d
  derived_key_length_bytes = %[7]d
  enabled                  = %[8]t
}`, resourceName,
		resourceModel.id,
		resourceModel.iterationCount,
		resourceModel.parallelismFactor,
		resourceModel.memoryUsageKb,
		resourceModel.saltLengthBytes,
		resourceModel.derivedKeyLengthBytes,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedPasswordStorageSchemeAttributes(config passwordStorageSchemeTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.PasswordStorageSchemeApi.GetPasswordStorageScheme(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Password Storage Scheme"
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "iteration-count",
			config.iterationCount, response.Argon2dPasswordStorageSchemeResponse.IterationCount)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "parallelism-factor",
			config.parallelismFactor, response.Argon2dPasswordStorageSchemeResponse.ParallelismFactor)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "memory-usage-kb",
			config.memoryUsageKb, response.Argon2dPasswordStorageSchemeResponse.MemoryUsageKb)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "salt-length-bytes",
			config.saltLengthBytes, response.Argon2dPasswordStorageSchemeResponse.SaltLengthBytes)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "derived-key-length-bytes",
			config.derivedKeyLengthBytes, response.Argon2dPasswordStorageSchemeResponse.DerivedKeyLengthBytes)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.Argon2dPasswordStorageSchemeResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckPasswordStorageSchemeDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.PasswordStorageSchemeApi.GetPasswordStorageScheme(ctx, testIdPasswordStorageScheme).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Password Storage Scheme", testIdPasswordStorageScheme)
	}
	return nil
}
