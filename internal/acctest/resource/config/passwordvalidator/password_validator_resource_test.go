package passwordvalidator_test

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

const testIdPasswordValidator = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type passwordValidatorTestModel struct {
	id                string
	minPasswordLength int64
	maxPasswordLength int64
	enabled           bool
}

func TestAccPasswordValidator(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := passwordValidatorTestModel{
		id:                testIdPasswordValidator,
		minPasswordLength: 8,
		maxPasswordLength: 100,
		enabled:           true,
	}
	updatedResourceModel := passwordValidatorTestModel{
		id:                testIdPasswordValidator,
		minPasswordLength: 6,
		maxPasswordLength: 0,
		enabled:           false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckPasswordValidatorDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccPasswordValidatorResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedPasswordValidatorAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_password_validator.%s", resourceName), "min_password_length", strconv.FormatInt(initialResourceModel.minPasswordLength, 10)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_password_validator.%s", resourceName), "max_password_length", strconv.FormatInt(initialResourceModel.maxPasswordLength, 10)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_password_validator.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
				),
			},
			{
				// Test updating some fields
				Config: testAccPasswordValidatorResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedPasswordValidatorAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccPasswordValidatorResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_password_validator." + resourceName,
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

func testAccPasswordValidatorResource(resourceName string, resourceModel passwordValidatorTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_password_validator" "%[1]s" {
  type                = "length-based"
  id                  = "%[2]s"
  min_password_length = %[3]d
  max_password_length = %[4]d
  enabled             = %[5]t
}

data "pingdirectory_password_validator" "%[1]s" {
  id = "%[2]s"
  depends_on = [
    pingdirectory_password_validator.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.minPasswordLength,
		resourceModel.maxPasswordLength,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedPasswordValidatorAttributes(config passwordValidatorTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.PasswordValidatorApi.GetPasswordValidator(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Password Validator"
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "min-password-length",
			config.minPasswordLength, *response.LengthBasedPasswordValidatorResponse.MinPasswordLength)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "max-password-length",
			config.maxPasswordLength, *response.LengthBasedPasswordValidatorResponse.MaxPasswordLength)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.LengthBasedPasswordValidatorResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckPasswordValidatorDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.PasswordValidatorApi.GetPasswordValidator(ctx, testIdPasswordValidator).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Password Validator", testIdPasswordValidator)
	}
	return nil
}