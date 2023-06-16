package otpdeliverymechanism_test

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

const testIdOtpDeliveryMechanism = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type otpDeliveryMechanismTestModel struct {
	id            string
	senderAddress string
	enabled       bool
}

func TestAccOtpDeliveryMechanism(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := otpDeliveryMechanismTestModel{
		id:            testIdOtpDeliveryMechanism,
		senderAddress: "sender@example.com",
		enabled:       true,
	}
	updatedResourceModel := otpDeliveryMechanismTestModel{
		id:            testIdOtpDeliveryMechanism,
		senderAddress: "updatedsender@example.com",
		enabled:       false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckOtpDeliveryMechanismDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccOtpDeliveryMechanismResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOtpDeliveryMechanismAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccOtpDeliveryMechanismResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedOtpDeliveryMechanismAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccOtpDeliveryMechanismResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_otp_delivery_mechanism." + resourceName,
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

func testAccOtpDeliveryMechanismResource(resourceName string, resourceModel otpDeliveryMechanismTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_otp_delivery_mechanism" "%[1]s" {
  type           = "email"
  id             = "%[2]s"
  sender_address = "%[3]s"
  enabled        = %[4]t
}`, resourceName,
		resourceModel.id,
		resourceModel.senderAddress,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedOtpDeliveryMechanismAttributes(config otpDeliveryMechanismTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OtpDeliveryMechanismApi.GetOtpDeliveryMechanism(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Otp Delivery Mechanism"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "sender-address",
			config.senderAddress, response.EmailOtpDeliveryMechanismResponse.SenderAddress)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.EmailOtpDeliveryMechanismResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckOtpDeliveryMechanismDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.OtpDeliveryMechanismApi.GetOtpDeliveryMechanism(ctx, testIdOtpDeliveryMechanism).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Otp Delivery Mechanism", testIdOtpDeliveryMechanism)
	}
	return nil
}
