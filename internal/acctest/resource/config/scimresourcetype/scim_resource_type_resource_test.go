// Copyright © 2025 Ping Identity Corporation

package scimresourcetype_test

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

const testIdLdapPassThroughScimResourceType = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type ldapPassThroughScimResourceTypeTestModel struct {
	id          string
	enabled     bool
	endpoint    string
	description string
}

func TestAccLdapPassThroughScimResourceType(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := ldapPassThroughScimResourceTypeTestModel{
		id:          testIdLdapPassThroughScimResourceType,
		enabled:     false,
		endpoint:    "endpoint",
		description: "initial",
	}
	updatedResourceModel := ldapPassThroughScimResourceTypeTestModel{
		id:          testIdLdapPassThroughScimResourceType,
		enabled:     false,
		endpoint:    "endpoint",
		description: "updated",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckLdapPassThroughScimResourceTypeDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccLdapPassThroughScimResourceTypeResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedLdapPassThroughScimResourceTypeAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_scim_resource_type.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_scim_resource_type.%s", resourceName), "endpoint", initialResourceModel.endpoint),
					resource.TestCheckResourceAttrSet("data.pingdirectory_scim_resource_types.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccLdapPassThroughScimResourceTypeResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedLdapPassThroughScimResourceTypeAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccLdapPassThroughScimResourceTypeResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_scim_resource_type." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.ScimResourceTypeAPI.DeleteScimResourceType(ctx, updatedResourceModel.id).Execute()
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

func testAccLdapPassThroughScimResourceTypeResource(resourceName string, resourceModel ldapPassThroughScimResourceTypeTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_scim_resource_type" "%[1]s" {
  type        = "ldap-pass-through"
  name        = "%[2]s"
  enabled     = %[3]t
  endpoint    = "%[4]s"
  description = "%[5]s"
}

data "pingdirectory_scim_resource_type" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_scim_resource_type.%[1]s
  ]
}

data "pingdirectory_scim_resource_types" "list" {
  depends_on = [
    pingdirectory_scim_resource_type.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.enabled,
		resourceModel.endpoint,
		resourceModel.description)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedLdapPassThroughScimResourceTypeAttributes(config ldapPassThroughScimResourceTypeTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ScimResourceTypeAPI.GetScimResourceType(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Ldap Pass Through Scim Resource Type"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.LdapPassThroughScimResourceTypeResponse.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "endpoint",
			config.endpoint, response.LdapPassThroughScimResourceTypeResponse.Endpoint)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.id, "description",
			config.description, response.LdapPassThroughScimResourceTypeResponse.Description)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckLdapPassThroughScimResourceTypeDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ScimResourceTypeAPI.GetScimResourceType(ctx, testIdLdapPassThroughScimResourceType).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Ldap Pass Through Scim Resource Type", testIdLdapPassThroughScimResourceType)
	}
	return nil
}
