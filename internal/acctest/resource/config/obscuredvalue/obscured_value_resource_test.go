// Copyright © 2025 Ping Identity Corporation

package obscuredvalue_test

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

const testIdObscuredValue = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type obscuredValueTestModel struct {
	id            string
	obscuredValue string
	description   string
}

func TestAccObscuredValue(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := obscuredValueTestModel{
		id:            testIdObscuredValue,
		obscuredValue: "myobscuredvalue",
		description:   "mydescription",
	}
	updatedResourceModel := obscuredValueTestModel{
		id:            testIdObscuredValue,
		obscuredValue: "mychangedobscuredvalue",
		description:   "mychangeddescription",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckObscuredValueDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccObscuredValueResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedObscuredValueAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_obscured_value.%s", resourceName), "description", initialResourceModel.description),
					resource.TestCheckResourceAttrSet("data.pingdirectory_obscured_values.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccObscuredValueResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedObscuredValueAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccObscuredValueResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_obscured_value." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"obscured_value",
				},
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.ObscuredValueAPI.DeleteObscuredValue(ctx, updatedResourceModel.id).Execute()
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

func testAccObscuredValueResource(resourceName string, resourceModel obscuredValueTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_obscured_value" "%[1]s" {
  name           = "%[2]s"
  obscured_value = "%[3]s"
  description    = "%[4]s"
}

data "pingdirectory_obscured_value" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_obscured_value.%[1]s
  ]
}

data "pingdirectory_obscured_values" "list" {
  depends_on = [
    pingdirectory_obscured_value.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.obscuredValue,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedObscuredValueAttributes(config obscuredValueTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ObscuredValueAPI.GetObscuredValue(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that description matches expected
		err = acctest.TestAttributesMatchStringPointer("obscured-value", &config.id, "description", config.description, response.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckObscuredValueDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ObscuredValueAPI.GetObscuredValue(ctx, testIdObscuredValue).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Obscured Value", testIdObscuredValue)
	}
	return nil
}
