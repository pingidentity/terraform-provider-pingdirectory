package postldifexporttaskprocessor_test

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

const testIdPostLdifExportTaskProcessor = "MyPostLdifProcessor"

// Attributes to test with. Add optional properties to test here if desired.
type postLdifExportTaskProcessorTestModel struct {
	name                 string
	s3BucketName         string
	enabled              bool
	maxFileCountToRetain int64
}

func TestAccPostLdifExportTaskProcessor(t *testing.T) {
	pdVersion := os.Getenv("PINGDIRECTORY_PROVIDER_PRODUCT_VERSION")
	compare, err := version.Compare(pdVersion, version.PingDirectory10000)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if compare < 0 {
		// This resource only exists in PD version 10.0 and later
		return
	}

	resourceName := "MyPostLdifProcessorResource"
	initialResourceModel := postLdifExportTaskProcessorTestModel{
		name:                 testIdPostLdifExportTaskProcessor,
		s3BucketName:         "myInitialBucket",
		enabled:              true,
		maxFileCountToRetain: 20,
	}
	updatedResourceModel := postLdifExportTaskProcessorTestModel{
		name:                 testIdPostLdifExportTaskProcessor,
		s3BucketName:         "myUpdatedBucket",
		enabled:              false,
		maxFileCountToRetain: 25,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckPostLdifExportTaskProcessorDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccPostLdifExportTaskProcessorResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedPostLdifExportTaskProcessorAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_post_ldif_export_task_processor.%s", resourceName), "s3_bucket_name", initialResourceModel.s3BucketName),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_post_ldif_export_task_processor.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_post_ldif_export_task_processors.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccPostLdifExportTaskProcessorResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedPostLdifExportTaskProcessorAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccPostLdifExportTaskProcessorResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_post_ldif_export_task_processor." + resourceName,
				ImportStateId:     updatedResourceModel.name,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.PostLdifExportTaskProcessorAPI.DeletePostLdifExportTaskProcessor(ctx, updatedResourceModel.name).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccPostLdifExportTaskProcessorResource(resourceName string, resourceModel postLdifExportTaskProcessorTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_external_server" "myAwsExternalServer" {
  type            = "amazon-aws"
  name            = "myaws"
  aws_region_name = "us-east-2"
}

resource "pingdirectory_post_ldif_export_task_processor" "%[1]s" {
  type                         = "upload-to-s3"
  name                         = "%[5]s"
  aws_external_server          = pingdirectory_external_server.myAwsExternalServer.name
  s3_bucket_name               = "%[2]s"
  enabled                      = %[3]t
  maximum_file_count_to_retain = %[4]d
}

data "pingdirectory_post_ldif_export_task_processor" "%[1]s" {
  name = "%[5]s"
  depends_on = [
    pingdirectory_post_ldif_export_task_processor.%[1]s
  ]
}

data "pingdirectory_post_ldif_export_task_processors" "list" {
  depends_on = [
    pingdirectory_post_ldif_export_task_processor.%[1]s
  ]
}`, resourceName,
		resourceModel.s3BucketName,
		resourceModel.enabled,
		resourceModel.maxFileCountToRetain,
		resourceModel.name)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedPostLdifExportTaskProcessorAttributes(config postLdifExportTaskProcessorTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, httpResp, err := testClient.PostLdifExportTaskProcessorAPI.GetPostLdifExportTaskProcessor(ctx, config.name).Execute()
		if err != nil {
			println("name", config.name)
			if httpResp != nil {
				body, internalError := io.ReadAll(httpResp.Body)
				if internalError == nil {
					println("Error HTTP response body: " + string(body))
				}
			}
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Post Ldif Export Task Processor"
		err = acctest.TestAttributesMatchString(resourceType, &config.name, "s3-bucket-name",
			config.s3BucketName, response.UploadToS3PostLdifExportTaskProcessorResponse.S3BucketName)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.name, "enabled",
			config.enabled, response.UploadToS3PostLdifExportTaskProcessorResponse.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.name, "max-file-count-to-retain",
			config.maxFileCountToRetain, *response.UploadToS3PostLdifExportTaskProcessorResponse.MaximumFileCountToRetain)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckPostLdifExportTaskProcessorDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.PostLdifExportTaskProcessorAPI.GetPostLdifExportTaskProcessor(ctx, testIdPostLdifExportTaskProcessor).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Post Ldif Export Task Processor", testIdPostLdifExportTaskProcessor)
	}
	return nil
}
