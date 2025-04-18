// Copyright © 2025 Ping Identity Corporation

package accesstokenvalidator_test

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

const testIdPingFederateAccessTokenValidator = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type pingFederateAccessTokenValidatorTestModel struct {
	id                   string
	clientId             string
	clientSecret         string
	authorizationServer  string
	enabled              bool
	evaluationOrderIndex int64
}

func TestAccPingFederateAccessTokenValidator(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := pingFederateAccessTokenValidatorTestModel{
		id:           testIdPingFederateAccessTokenValidator,
		clientId:     "my-client-id",
		clientSecret: "myclientsecret",
		// In reality you wouldn't use this authorization server, just using it because it's available by default
		// and an authorization server is required to create this access token validator.
		authorizationServer:  "PingOne Auth Service",
		enabled:              true,
		evaluationOrderIndex: 1,
	}
	updatedResourceModel := pingFederateAccessTokenValidatorTestModel{
		id:                   testIdPingFederateAccessTokenValidator,
		clientId:             "my-client-id-updated",
		clientSecret:         "myclientsecretupdated",
		authorizationServer:  "PingOne Auth Service",
		enabled:              false,
		evaluationOrderIndex: 2,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckPingFederateAccessTokenValidatorDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccPingFederateAccessTokenValidatorResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedPingFederateAccessTokenValidatorAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_access_token_validator.%s", resourceName), "type", "ping-federate"),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_access_token_validator.%s", resourceName), "client_id", initialResourceModel.clientId),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_access_token_validator.%s", resourceName), "authorization_server", initialResourceModel.authorizationServer),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_access_token_validator.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_access_token_validators.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccPingFederateAccessTokenValidatorResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedPingFederateAccessTokenValidatorAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccPingFederateAccessTokenValidatorResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_access_token_validator." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				// Can't verify import state for a sensitive attribute that PD won't return
				ImportStateVerifyIgnore: []string{"client_secret"},
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.AccessTokenValidatorAPI.DeleteAccessTokenValidator(ctx, updatedResourceModel.id).Execute()
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

func testAccPingFederateAccessTokenValidatorResource(resourceName string, resourceModel pingFederateAccessTokenValidatorTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_access_token_validator" "%[1]s" {
  type                   = "ping-federate"
  name                   = "%[2]s"
  client_id              = "%[3]s"
  client_secret          = "%[4]s"
  authorization_server   = "%[5]s"
  enabled                = %[6]t
  evaluation_order_index = %[7]d
}

data "pingdirectory_access_token_validator" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_access_token_validator.%[1]s
  ]
}

data "pingdirectory_access_token_validators" "list" {
  depends_on = [
    pingdirectory_access_token_validator.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.clientId,
		resourceModel.clientSecret,
		resourceModel.authorizationServer,
		resourceModel.enabled,
		resourceModel.evaluationOrderIndex)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedPingFederateAccessTokenValidatorAttributes(config pingFederateAccessTokenValidatorTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AccessTokenValidatorAPI.GetAccessTokenValidator(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Ping Federate Access Token Validator"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "client-id",
			config.clientId, response.PingFederateAccessTokenValidatorResponse.ClientID)
		if err != nil {
			return err
		}
		// Unable to check client-secret since it can't be retrieved from the PD configuration API
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "authorization-server",
			config.authorizationServer, response.PingFederateAccessTokenValidatorResponse.AuthorizationServer)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.PingFederateAccessTokenValidatorResponse.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "evaluation-order-index",
			config.evaluationOrderIndex, response.PingFederateAccessTokenValidatorResponse.EvaluationOrderIndex)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckPingFederateAccessTokenValidatorDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.AccessTokenValidatorAPI.GetAccessTokenValidator(ctx, testIdPingFederateAccessTokenValidator).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Ping Federate Access Token Validator", testIdPingFederateAccessTokenValidator)
	}
	return nil
}
