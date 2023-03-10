package config_test

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

const defaultLocationId = "Docker"
const defaultLogPublisherId = "File-Based Audit Logger"

// Attributes to test with. Add optional properties to test here if desired.
type defaultLocationTestModel struct {
	id          string
	description string
}

type defaultLogPublisherTestModel struct {
	id      string
	enabled bool
}

func TestAccDefaultLocation(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := defaultLocationTestModel{
		id:          defaultLocationId,
		description: "test",
	}
	updatedResourceModel := defaultLocationTestModel{
		id:          defaultLocationId,
		description: "updated",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckDefaultLocationDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDefaultLocationResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedDefaultLocationAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccDefaultLocationResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedDefaultLocationAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccDefaultLocationResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_default_location." + resourceName,
				ImportStateId:           updatedResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccDefaultLocationResource(resourceName string, resourceModel defaultLocationTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_location" "%[1]s" {
  id          = "%[2]s"
  description = "%[3]s"
}`, resourceName, resourceModel.id,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDefaultLocationAttributes(config defaultLocationTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LocationApi.GetLocation(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Default Location"
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "description",
			config.description, response.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Verify that the default resources are NOT destroyed
func testAccCheckDefaultLocationDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.LocationApi.GetLocation(ctx, defaultLocationId).Execute()
	if err != nil {
		return err
	}
	return nil
}

func TestAccDefaultLogPublisher(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := defaultLogPublisherTestModel{
		id:      defaultLogPublisherId,
		enabled: true,
	}
	updatedResourceModel := defaultLogPublisherTestModel{
		id:      defaultLogPublisherId,
		enabled: false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckDefaultLogPublisherDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDefaultLogPublisherResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedDefaultLogPublisherAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccDefaultLogPublisherResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedDefaultLogPublisherAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccDefaultLogPublisherResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_default_file_based_audit_log_publisher." + resourceName,
				ImportStateId:           updatedResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccDefaultLogPublisherResource(resourceName string, resourceModel defaultLogPublisherTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_file_based_audit_log_publisher" "%[1]s" {
  id      = "%[2]s"
  enabled = "%[3]t"
}`, resourceName, resourceModel.id,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDefaultLogPublisherAttributes(config defaultLogPublisherTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LogPublisherApi.GetLogPublisher(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Default File-Based Audit Log Publisher"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.FileBasedAuditLogPublisherResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Verify that the default resources are NOT destroyed
func testAccCheckDefaultLogPublisherDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.LogPublisherApi.GetLogPublisher(ctx, defaultLogPublisherId).Execute()
	if err != nil {
		return err
	}
	return nil
}
