package logfieldmapping_test

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

const testIdAccessLogFieldMapping = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type accessLogFieldMappingTestModel struct {
	id                string
	description       string
	log_field_message string
}

func TestAccAccessLogFieldMapping(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := accessLogFieldMappingTestModel{
		id:                testIdAccessLogFieldMapping,
		description:       "My error log field mapping",
		log_field_message: "message",
	}
	updatedResourceModel := accessLogFieldMappingTestModel{
		id:                testIdAccessLogFieldMapping,
		description:       "Updated error log field mapping",
		log_field_message: "updatedMessage",
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckAccessLogFieldMappingDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccAccessLogFieldMappingResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedAccessLogFieldMappingAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_log_field_mapping.%s", resourceName), "description", initialResourceModel.description),
					resource.TestCheckResourceAttrSet("data.pingdirectory_log_field_mappings.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccAccessLogFieldMappingResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedAccessLogFieldMappingAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccAccessLogFieldMappingResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_log_field_mapping." + resourceName,
				ImportStateId:     initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccAccessLogFieldMappingResource(resourceName string, resourceModel accessLogFieldMappingTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_log_field_mapping" "%[1]s" {
  type              = "access"
  name              = "%[2]s"
  description       = "%[3]s"
  log_field_message = "%[4]s"
}

data "pingdirectory_log_field_mapping" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_log_field_mapping.%[1]s
  ]
}

data "pingdirectory_log_field_mappings" "list" {
  depends_on = [
    pingdirectory_log_field_mapping.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.description,
		resourceModel.log_field_message)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedAccessLogFieldMappingAttributes(config accessLogFieldMappingTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LogFieldMappingApi.GetLogFieldMapping(ctx, config.id).Execute()
		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		resourceType := "Access Log Field Mapping"
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.description, "description",
			config.description, response.AccessLogFieldMappingResponse.Description)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.log_field_message, "log-field-message",
			config.log_field_message, response.AccessLogFieldMappingResponse.LogFieldMessage)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckAccessLogFieldMappingDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.LogFieldMappingApi.GetLogFieldMapping(ctx, testIdAccessLogFieldMapping).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Access Log Field Mapping", testIdAccessLogFieldMapping)
	}
	return nil
}
