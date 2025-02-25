// Copyright © 2025 Ping Identity Corporation

package conjurauthenticationmethod_test

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

const testIdConjurAuthenticationMethod = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type conjurAuthenticationMethodTestModel struct {
	id       string
	username string
	password string
}

func TestAccConjurAuthenticationMethod(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := conjurAuthenticationMethodTestModel{
		id:       testIdConjurAuthenticationMethod,
		username: "firstusername",
		password: "password",
	}
	updatedResourceModel := conjurAuthenticationMethodTestModel{
		id:       testIdConjurAuthenticationMethod,
		username: "secondusername",
		password: "password2",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckConjurAuthenticationMethodDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccConjurAuthenticationMethodResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedConjurAuthenticationMethodAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_conjur_authentication_method.%s", resourceName), "username", initialResourceModel.username),
					resource.TestCheckResourceAttrSet("data.pingdirectory_conjur_authentication_methods.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccConjurAuthenticationMethodResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedConjurAuthenticationMethodAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccConjurAuthenticationMethodResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_conjur_authentication_method." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"api_key",
				},
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.ConjurAuthenticationMethodAPI.DeleteConjurAuthenticationMethod(ctx, updatedResourceModel.id).Execute()
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

func testAccConjurAuthenticationMethodResource(resourceName string, resourceModel conjurAuthenticationMethodTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_conjur_authentication_method" "%[1]s" {
  name     = "%[2]s"
  username = "%[3]s"
  password = "%[4]s"
}

data "pingdirectory_conjur_authentication_method" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_conjur_authentication_method.%[1]s
  ]
}

data "pingdirectory_conjur_authentication_methods" "list" {
  depends_on = [
    pingdirectory_conjur_authentication_method.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.username,
		resourceModel.password)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedConjurAuthenticationMethodAttributes(config conjurAuthenticationMethodTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ConjurAuthenticationMethodAPI.GetConjurAuthenticationMethod(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Conjur Authentication Method"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "username",
			config.username, response.Username)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckConjurAuthenticationMethodDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ConjurAuthenticationMethodAPI.GetConjurAuthenticationMethod(ctx, testIdConjurAuthenticationMethod).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Conjur Authentication Method", testIdConjurAuthenticationMethod)
	}
	return nil
}
