package passwordpolicy_test

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

const testIdPasswordPolicy = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type passwordPolicyTestModel struct {
	id                           string
	passwordAttribute            string
	defaultPasswordStorageScheme []string
	allowPreEncodedPasswords     string
}

func TestAccPasswordPolicy(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := passwordPolicyTestModel{
		id:                           testIdPasswordPolicy,
		passwordAttribute:            "userPassword",
		defaultPasswordStorageScheme: []string{"Blowfish"},
		allowPreEncodedPasswords:     "false",
	}
	updatedResourceModel := passwordPolicyTestModel{
		id:                           testIdPasswordPolicy,
		passwordAttribute:            "userPassword",
		defaultPasswordStorageScheme: []string{"Salted SHA-512"},
		allowPreEncodedPasswords:     "true",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckPasswordPolicyDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccPasswordPolicyResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedPasswordPolicyAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_password_policy.%s", resourceName), "password_attribute", initialResourceModel.passwordAttribute),
					resource.TestCheckTypeSetElemAttr(fmt.Sprintf("data.pingdirectory_password_policy.%s", resourceName), "default_password_storage_scheme.*", initialResourceModel.defaultPasswordStorageScheme[0]),
					resource.TestCheckResourceAttrSet("data.pingdirectory_password_policies.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccPasswordPolicyResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedPasswordPolicyAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccPasswordPolicyResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_password_policy." + resourceName,
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

func testAccPasswordPolicyResource(resourceName string, resourceModel passwordPolicyTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_password_policy" "%[1]s" {
  id                              = "%[2]s"
  password_attribute              = "%[3]s"
  default_password_storage_scheme = %[4]s
  allow_pre_encoded_passwords     = "%[5]s"
}

data "pingdirectory_password_policy" "%[1]s" {
  id = "%[2]s"
  depends_on = [
    pingdirectory_password_policy.%[1]s
  ]
}

data "pingdirectory_password_policies" "list" {
  depends_on = [
    pingdirectory_password_policy.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.passwordAttribute,
		acctest.StringSliceToTerraformString(resourceModel.defaultPasswordStorageScheme),
		resourceModel.allowPreEncodedPasswords)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedPasswordPolicyAttributes(config passwordPolicyTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.PasswordPolicyApi.GetPasswordPolicy(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Password Policy"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "password-attribute",
			config.passwordAttribute, response.PasswordAttribute)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "default-password-storage-scheme",
			config.defaultPasswordStorageScheme, response.DefaultPasswordStorageScheme)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "allow-pre-encoded-passwords",
			config.allowPreEncodedPasswords, response.AllowPreEncodedPasswords.String())
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckPasswordPolicyDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.PasswordPolicyApi.GetPasswordPolicy(ctx, testIdPasswordPolicy).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Password Policy", testIdPasswordPolicy)
	}
	return nil
}
