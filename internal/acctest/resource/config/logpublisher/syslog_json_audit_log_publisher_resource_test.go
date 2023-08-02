package logpublisher_test

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

const testIdSyslogJsonAuditLogPublisher = "MyId"
const testIdSyslogExternalServer = "externalServerId"

// Attributes to test with. Add optional properties to test here if desired.
type syslogJsonAuditLogPublisherTestModel struct {
	id                   string
	syslogExternalServer string
	enabled              bool
}

func TestAccSyslogJsonAuditLogPublisher(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := syslogJsonAuditLogPublisherTestModel{
		id:                   testIdSyslogJsonAuditLogPublisher,
		syslogExternalServer: testIdSyslogExternalServer,
		enabled:              false,
	}
	updatedResourceModel := syslogJsonAuditLogPublisherTestModel{
		id:                   testIdSyslogJsonAuditLogPublisher,
		syslogExternalServer: testIdSyslogExternalServer,
		enabled:              true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckSyslogJsonAuditLogPublisherDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSyslogJsonAuditLogPublisherResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSyslogJsonAuditLogPublisherAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccSyslogJsonAuditLogPublisherResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSyslogJsonAuditLogPublisherAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSyslogJsonAuditLogPublisherResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_log_publisher." + resourceName,
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

func testAccSyslogJsonAuditLogPublisherResource(resourceName string, resourceModel syslogJsonAuditLogPublisherTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_external_server" "%[3]s" {
  type                = "syslog"
  name                = "%[3]s"
  server_host_name    = "localhost"
  transport_mechanism = "tls-encrypted-tcp"
}

resource "pingdirectory_log_publisher" "%[1]s" {
  type                   = "syslog-json-audit"
  name                   = "%[2]s"
  syslog_external_server = [pingdirectory_external_server.%[3]s.id]
  enabled                = %[4]t
}`, resourceName,
		resourceModel.id,
		resourceModel.syslogExternalServer,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSyslogJsonAuditLogPublisherAttributes(config syslogJsonAuditLogPublisherTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LogPublisherApi.GetLogPublisher(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Syslog Json Audit Log Publisher"
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "syslog-external-server",
			[]string{config.syslogExternalServer}, response.SyslogJsonAuditLogPublisherResponse.SyslogExternalServer)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.SyslogJsonAuditLogPublisherResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSyslogJsonAuditLogPublisherDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.LogPublisherApi.GetLogPublisher(ctx, testIdSyslogJsonAuditLogPublisher).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Syslog Json Audit Log Publisher", testIdSyslogJsonAuditLogPublisher)
	}
	return nil
}
