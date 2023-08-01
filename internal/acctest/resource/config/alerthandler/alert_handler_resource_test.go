package alerthandler_test

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

const testIdSmtpAlertHandler = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type smtpAlertHandlerTestModel struct {
	id               string
	senderAddress    string
	recipientAddress []string
	enabled          bool
}

func TestAccSmtpAlertHandler(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := smtpAlertHandlerTestModel{
		id:               testIdSmtpAlertHandler,
		senderAddress:    "unboundid-notifications@example.com",
		recipientAddress: []string{"test@example.com", "users@example.com"},
		enabled:          false,
	}
	updatedResourceModel := smtpAlertHandlerTestModel{
		id:               testIdSmtpAlertHandler,
		senderAddress:    "unboundid-user-notifications@example.com",
		recipientAddress: []string{"testing@example.com", "endusers@example.com"},
		enabled:          false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckSmtpAlertHandlerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSmtpAlertHandlerResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedSmtpAlertHandlerAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_alert_handler.%s", resourceName), "sender_address", initialResourceModel.senderAddress),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_alert_handler.%s", resourceName), "recipient_address.*", initialResourceModel.recipientAddress[0]),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_alert_handler.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_alert_handlers.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccSmtpAlertHandlerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSmtpAlertHandlerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSmtpAlertHandlerResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_alert_handler." + resourceName,
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

func testAccSmtpAlertHandlerResource(resourceName string, resourceModel smtpAlertHandlerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_alert_handler" "%[1]s" {
  type              = "smtp"
  id                = "%[2]s"
  sender_address    = "%[3]s"
  recipient_address = %[4]s
  enabled           = %[5]t
}

data "pingdirectory_alert_handler" "%[1]s" {
  id = "%[2]s"
  depends_on = [
    pingdirectory_alert_handler.%[1]s
  ]
}

data "pingdirectory_alert_handlers" "list" {
  depends_on = [
    pingdirectory_alert_handler.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.senderAddress,
		acctest.StringSliceToTerraformString(resourceModel.recipientAddress),
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSmtpAlertHandlerAttributes(config smtpAlertHandlerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AlertHandlerApi.GetAlertHandler(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Smtp Alert Handler"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "sender-address",
			config.senderAddress, response.SmtpAlertHandlerResponse.SenderAddress)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "recipient-address",
			config.recipientAddress, response.SmtpAlertHandlerResponse.RecipientAddress)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.SmtpAlertHandlerResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSmtpAlertHandlerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.AlertHandlerApi.GetAlertHandler(ctx, testIdSmtpAlertHandler).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Smtp Alert Handler", testIdSmtpAlertHandler)
	}
	return nil
}
