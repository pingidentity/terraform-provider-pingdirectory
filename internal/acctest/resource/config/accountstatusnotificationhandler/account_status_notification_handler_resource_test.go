package accountstatusnotificationhandler_test

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

const testIdSmtpAccountStatusNotificationHandler = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type smtpAccountStatusNotificationHandlerTestModel struct {
	id                               string
	sendMessageWithoutEndUserAddress bool
	recipientAddress                 []string
	senderAddress                    string
	messageSubject                   []string
	messageTemplateFile              []string
	enabled                          bool
}

func TestAccSmtpAccountStatusNotificationHandler(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := smtpAccountStatusNotificationHandlerTestModel{
		id:                               testIdSmtpAccountStatusNotificationHandler,
		sendMessageWithoutEndUserAddress: false,
		recipientAddress:                 []string{"test@example.com", "users@example.com"},
		senderAddress:                    "unboundid-notifications@example.com",
		messageSubject:                   []string{"account-disabled:Your directory account has been disabled"},
		messageTemplateFile:              []string{"account-disabled:config/messages/account-disabled.template"},
		enabled:                          false,
	}
	updatedResourceModel := smtpAccountStatusNotificationHandlerTestModel{
		id:                               testIdSmtpAccountStatusNotificationHandler,
		sendMessageWithoutEndUserAddress: true,
		recipientAddress:                 []string{"testing@example.com", "endusers@example.com"},
		senderAddress:                    "unboundid-user-notifications@example.com",
		messageSubject:                   []string{"account-enabled:Your directory account has been re-enabled"},
		messageTemplateFile:              []string{"account-enabled:config/messages/account-enabled.template"},
		enabled:                          false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckSmtpAccountStatusNotificationHandlerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSmtpAccountStatusNotificationHandlerResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedSmtpAccountStatusNotificationHandlerAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_account_status_notification_handler.%s", resourceName), "sender_address", initialResourceModel.senderAddress),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_account_status_notification_handler.%s", resourceName), "message_subject.*", initialResourceModel.messageSubject[0]),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_account_status_notification_handler.%s", resourceName), "message_template_file.*", initialResourceModel.messageTemplateFile[0]),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_account_status_notification_handler.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
				),
			},
			{
				// Test updating some fields
				Config: testAccSmtpAccountStatusNotificationHandlerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSmtpAccountStatusNotificationHandlerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccSmtpAccountStatusNotificationHandlerResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_account_status_notification_handler." + resourceName,
				ImportStateId:           updatedResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccSmtpAccountStatusNotificationHandlerResource(resourceName string, resourceModel smtpAccountStatusNotificationHandlerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_account_status_notification_handler" "%[1]s" {
  type                                  = "smtp"
  id                                    = "%[2]s"
  send_message_without_end_user_address = %[3]t
  recipient_address                     = %[4]s
  sender_address                        = "%[5]s"
  message_subject                       = %[6]s
  message_template_file                 = %[7]s
  enabled                               = %[8]t
}

data "pingdirectory_account_status_notification_handler" "%[1]s" {
  id = "%[2]s"
  depends_on = [
    pingdirectory_account_status_notification_handler.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.sendMessageWithoutEndUserAddress,
		acctest.StringSliceToTerraformString(resourceModel.recipientAddress),
		resourceModel.senderAddress,
		acctest.StringSliceToTerraformString(resourceModel.messageSubject),
		acctest.StringSliceToTerraformString(resourceModel.messageTemplateFile),
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSmtpAccountStatusNotificationHandlerAttributes(config smtpAccountStatusNotificationHandlerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Smtp Account Status Notification Handler"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "send-message-without-end-user-address",
			config.sendMessageWithoutEndUserAddress, response.SmtpAccountStatusNotificationHandlerResponse.SendMessageWithoutEndUserAddress)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "recipient-address",
			config.recipientAddress, response.SmtpAccountStatusNotificationHandlerResponse.RecipientAddress)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "sender-address",
			config.senderAddress, response.SmtpAccountStatusNotificationHandlerResponse.SenderAddress)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "message-subject",
			config.messageSubject, response.SmtpAccountStatusNotificationHandlerResponse.MessageSubject)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "message-template-file",
			config.messageTemplateFile, response.SmtpAccountStatusNotificationHandlerResponse.MessageTemplateFile)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.SmtpAccountStatusNotificationHandlerResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSmtpAccountStatusNotificationHandlerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(ctx, testIdSmtpAccountStatusNotificationHandler).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Smtp Account Status Notification Handler", testIdSmtpAccountStatusNotificationHandler)
	}
	return nil
}
