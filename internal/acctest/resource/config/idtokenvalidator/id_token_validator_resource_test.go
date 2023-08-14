package idtokenvalidator_test

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

const testIdPingOneIdTokenValidator = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type pingOneIdTokenValidatorTestModel struct {
	id                   string
	issuerUrl            string
	enabled              bool
	identityMapper       string
	evaluationOrderIndex int64
}

func TestAccPingOneIdTokenValidator(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := pingOneIdTokenValidatorTestModel{
		id:                   testIdPingOneIdTokenValidator,
		issuerUrl:            "example.com",
		enabled:              true,
		identityMapper:       "Exact Match",
		evaluationOrderIndex: 1,
	}
	updatedResourceModel := pingOneIdTokenValidatorTestModel{
		id:                   testIdPingOneIdTokenValidator,
		issuerUrl:            "example.org",
		enabled:              false,
		identityMapper:       "Root DN Users",
		evaluationOrderIndex: 2,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckPingOneIdTokenValidatorDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccPingOneIdTokenValidatorResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedPingOneIdTokenValidatorAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_id_token_validator.%s", resourceName), "issuer_url", initialResourceModel.issuerUrl),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_id_token_validator.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_id_token_validator.%s", resourceName), "identity_mapper", initialResourceModel.identityMapper),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_id_token_validator.%s", resourceName), "evaluation_order_index", strconv.FormatInt(initialResourceModel.evaluationOrderIndex, 10)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_id_token_validators.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccPingOneIdTokenValidatorResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedPingOneIdTokenValidatorAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccPingOneIdTokenValidatorResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_id_token_validator." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.IdTokenValidatorApi.DeleteIdTokenValidator(ctx, updatedResourceModel.id).Execute()
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

func testAccPingOneIdTokenValidatorResource(resourceName string, resourceModel pingOneIdTokenValidatorTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_id_token_validator" "%[1]s" {
  type                   = "ping-one"
  name                   = "%[2]s"
  issuer_url             = "%[3]s"
  enabled                = %[4]t
  identity_mapper        = "%[5]s"
  evaluation_order_index = %[6]d
}

data "pingdirectory_id_token_validator" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_id_token_validator.%[1]s
  ]
}

data "pingdirectory_id_token_validators" "list" {
  depends_on = [
    pingdirectory_id_token_validator.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.issuerUrl,
		resourceModel.enabled,
		resourceModel.identityMapper,
		resourceModel.evaluationOrderIndex)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedPingOneIdTokenValidatorAttributes(config pingOneIdTokenValidatorTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.IdTokenValidatorApi.GetIdTokenValidator(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Ping One Id Token Validator"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "issuer-url",
			config.issuerUrl, response.PingOneIdTokenValidatorResponse.IssuerURL)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.PingOneIdTokenValidatorResponse.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "identity-mapper",
			config.identityMapper, response.PingOneIdTokenValidatorResponse.IdentityMapper)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "evaluation-order-index",
			config.evaluationOrderIndex, response.PingOneIdTokenValidatorResponse.EvaluationOrderIndex)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckPingOneIdTokenValidatorDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.IdTokenValidatorApi.GetIdTokenValidator(ctx, testIdPingOneIdTokenValidator).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Ping One Id Token Validator", testIdPingOneIdTokenValidator)
	}
	return nil
}
