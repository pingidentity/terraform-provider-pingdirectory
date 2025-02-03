// Copyright © 2025 Ping Identity Corporation

package otpdeliverymechanism_test

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
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckOtpDeliveryMechanismDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccOtpDeliveryMechanismResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedOtpDeliveryMechanismAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_otp_delivery_mechanism.%s", resourceName), "sender_address", initialResourceModel.senderAddress),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_otp_delivery_mechanism.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_otp_delivery_mechanisms.list", "objects.0.id"),
				),
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
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.OtpDeliveryMechanismAPI.DeleteOtpDeliveryMechanism(ctx, updatedResourceModel.id).Execute()
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

func testAccOtpDeliveryMechanismResource(resourceName string, resourceModel otpDeliveryMechanismTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_otp_delivery_mechanism" "%[1]s" {
  type           = "email"
  name           = "%[2]s"
  sender_address = "%[3]s"
  enabled        = %[4]t
}

data "pingdirectory_otp_delivery_mechanism" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_otp_delivery_mechanism.%[1]s
  ]
}

data "pingdirectory_otp_delivery_mechanisms" "list" {
  depends_on = [
    pingdirectory_otp_delivery_mechanism.%[1]s
  ]
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
		response, _, err := testClient.OtpDeliveryMechanismAPI.GetOtpDeliveryMechanism(ctx, config.id).Execute()
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
	_, _, err := testClient.OtpDeliveryMechanismAPI.GetOtpDeliveryMechanism(ctx, testIdOtpDeliveryMechanism).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Otp Delivery Mechanism", testIdOtpDeliveryMechanism)
	}
	return nil
}
