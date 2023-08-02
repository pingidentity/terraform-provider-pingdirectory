package logpublisher_test

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

const testIdFileBasedAccessLogPublisher = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type fileBasedAccessLogPublisherTestModel struct {
	id                 string
	logFile            string
	logFilePermissions string
	rotationPolicy     []string
	retentionPolicy    []string
	asynchronous       bool
	enabled            bool
}

func TestAccFileBasedAccessLogPublisher(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := fileBasedAccessLogPublisherTestModel{
		id:                 testIdFileBasedAccessLogPublisher,
		logFile:            "logs/example.log",
		logFilePermissions: "600",
		rotationPolicy:     []string{"Size Limit Rotation Policy"},
		retentionPolicy:    []string{"Never Delete"},
		asynchronous:       false,
		enabled:            true,
	}
	updatedResourceModel := fileBasedAccessLogPublisherTestModel{
		id:                 testIdFileBasedAccessLogPublisher,
		logFile:            "logs/example2.log",
		logFilePermissions: "606",
		rotationPolicy:     []string{"Never Rotate"},
		retentionPolicy:    []string{"File Count Retention Policy"},
		asynchronous:       true,
		enabled:            false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckFileBasedAccessLogPublisherDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccFileBasedAccessLogPublisherResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedFileBasedAccessLogPublisherAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_log_publisher.%s", resourceName), "log_file", initialResourceModel.logFile),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_log_publisher.%s", resourceName), "retention_policy.*", initialResourceModel.retentionPolicy[0]),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_log_publisher.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_log_publishers.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccFileBasedAccessLogPublisherResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedFileBasedAccessLogPublisherAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccFileBasedAccessLogPublisherResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_log_publisher." + resourceName,
				ImportStateId:           updatedResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccFileBasedAccessLogPublisherResource(resourceName string, resourceModel fileBasedAccessLogPublisherTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_log_publisher" "%[1]s" {
  type                 = "file-based-access"
  name                 = "%[2]s"
  log_file             = "%[3]s"
  log_file_permissions = "%[4]s"
  rotation_policy      = %[5]s
  retention_policy     = %[6]s
  asynchronous         = %[7]t
  enabled              = %[8]t
}

data "pingdirectory_log_publisher" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_log_publisher.%[1]s
  ]
}

data "pingdirectory_log_publishers" "list" {
  depends_on = [
    pingdirectory_log_publisher.%[1]s
  ]
}`, resourceName, resourceModel.id,
		resourceModel.logFile,
		resourceModel.logFilePermissions,
		acctest.StringSliceToTerraformString(resourceModel.rotationPolicy),
		acctest.StringSliceToTerraformString(resourceModel.retentionPolicy),
		resourceModel.asynchronous,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedFileBasedAccessLogPublisherAttributes(config fileBasedAccessLogPublisherTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LogPublisherApi.GetLogPublisher(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "File Based Access Log Publisher"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "log-file",
			config.logFile, response.FileBasedAccessLogPublisherResponse.LogFile)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "log-file-permissions",
			config.logFilePermissions, response.FileBasedAccessLogPublisherResponse.LogFilePermissions)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "rotation-policy",
			config.rotationPolicy, response.FileBasedAccessLogPublisherResponse.RotationPolicy)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "retention-policy",
			config.retentionPolicy, response.FileBasedAccessLogPublisherResponse.RetentionPolicy)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "asynchronous",
			config.asynchronous, response.FileBasedAccessLogPublisherResponse.Asynchronous)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.FileBasedAccessLogPublisherResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckFileBasedAccessLogPublisherDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.LogPublisherApi.GetLogPublisher(ctx, testIdFileBasedAccessLogPublisher).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("File Based Access Log Publisher", testIdFileBasedAccessLogPublisher)
	}
	return nil
}
