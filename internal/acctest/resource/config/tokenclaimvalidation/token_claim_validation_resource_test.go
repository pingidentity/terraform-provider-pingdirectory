package tokenclaimvalidation_test

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

const testIdStringArrayTokenClaimValidation = "MyId"
const testIdTokenValidatorName = "MyOpenidConnectIdTokenValidator"

// Attributes to test with. Add optional properties to test here if desired.
type stringArrayTokenClaimValidationTestModel struct {
	id                   string
	idTokenValidatorName string
	anyRequiredValue     []string
	claimName            string
}

func TestAccStringArrayTokenClaimValidation(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := stringArrayTokenClaimValidationTestModel{
		id:                   testIdStringArrayTokenClaimValidation,
		idTokenValidatorName: testIdTokenValidatorName,
		anyRequiredValue:     []string{"my_example_value"},
		claimName:            "my_example_name",
	}
	updatedResourceModel := stringArrayTokenClaimValidationTestModel{
		id:                   testIdStringArrayTokenClaimValidation,
		idTokenValidatorName: testIdTokenValidatorName,
		anyRequiredValue:     []string{"my_example_value"},
		claimName:            "my_example_name_update",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckStringArrayTokenClaimValidationDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccStringArrayTokenClaimValidationResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedStringArrayTokenClaimValidationAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_token_claim_validation.%s", resourceName), "claim_name", initialResourceModel.claimName),
					resource.TestCheckResourceAttrSet("data.pingdirectory_token_claim_validations.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccStringArrayTokenClaimValidationResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedStringArrayTokenClaimValidationAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccStringArrayTokenClaimValidationResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_token_claim_validation." + resourceName,
				ImportStateId:     updatedResourceModel.idTokenValidatorName + "/" + updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.TokenClaimValidationApi.DeleteTokenClaimValidation(ctx, updatedResourceModel.id, updatedResourceModel.idTokenValidatorName).Execute()
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

func testAccStringArrayTokenClaimValidationResource(resourceName string, resourceModel stringArrayTokenClaimValidationTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_id_token_validator" "%[3]s" {
  type                   = "ping-one"
  name                   = "%[3]s"
  issuer_url             = "example.com"
  enabled                = false
  identity_mapper        = "Exact Match"
  evaluation_order_index = 1
}

resource "pingdirectory_token_claim_validation" "%[1]s" {
  type                    = "string-array"
  name                    = "%[2]s"
  id_token_validator_name = pingdirectory_id_token_validator.%[3]s.id
  any_required_value      = %[4]s
  claim_name              = "%[5]s"
}

data "pingdirectory_token_claim_validation" "%[1]s" {
  name                    = "%[2]s"
  id_token_validator_name = "%[3]s"
  depends_on = [
    pingdirectory_token_claim_validation.%[1]s
  ]
}

data "pingdirectory_token_claim_validations" "list" {
  id_token_validator_name = "%[3]s"
  depends_on = [
    pingdirectory_token_claim_validation.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.idTokenValidatorName,
		acctest.StringSliceToTerraformString(resourceModel.anyRequiredValue),
		resourceModel.claimName)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedStringArrayTokenClaimValidationAttributes(config stringArrayTokenClaimValidationTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.TokenClaimValidationApi.GetTokenClaimValidation(ctx, config.id, config.idTokenValidatorName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "String Array Token Claim Validation"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "claim-name",
			config.claimName, response.StringArrayTokenClaimValidationResponse.ClaimName)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "any-required-value",
			config.anyRequiredValue, response.StringArrayTokenClaimValidationResponse.AnyRequiredValue)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckStringArrayTokenClaimValidationDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.TokenClaimValidationApi.GetTokenClaimValidation(ctx, testIdStringArrayTokenClaimValidation, testIdTokenValidatorName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("String Array Token Claim Validation", testIdStringArrayTokenClaimValidation)
	}
	return nil
}
