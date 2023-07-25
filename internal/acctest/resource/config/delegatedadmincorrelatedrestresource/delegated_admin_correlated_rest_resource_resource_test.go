package delegatedadmincorrelatedrestresource_test

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

const testIdDelegatedAdminCorrelatedRestResource = "MyId"
const testRestResourceTypeName = "restresource"

// Attributes to test with. Add optional properties to test here if desired.
type delegatedAdminCorrelatedRestResourceTestModel struct {
	id                                        string
	restResourceTypeName                      string
	displayName                               string
	correlatedRestResource                    string
	primaryRestResourceCorrelationAttribute   string
	secondaryRestResourceCorrelationAttribute string
}

func TestAccDelegatedAdminCorrelatedRestResource(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := delegatedAdminCorrelatedRestResourceTestModel{
		id:                                      testIdDelegatedAdminCorrelatedRestResource,
		restResourceTypeName:                    testRestResourceTypeName,
		displayName:                             "displayname",
		correlatedRestResource:                  testRestResourceTypeName,
		primaryRestResourceCorrelationAttribute: "cn",
		secondaryRestResourceCorrelationAttribute: "sn",
	}
	updatedResourceModel := delegatedAdminCorrelatedRestResourceTestModel{
		id:                                      testIdDelegatedAdminCorrelatedRestResource,
		restResourceTypeName:                    testRestResourceTypeName,
		displayName:                             "newdisplayname",
		correlatedRestResource:                  testRestResourceTypeName,
		primaryRestResourceCorrelationAttribute: "sn",
		secondaryRestResourceCorrelationAttribute: "cn",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckDelegatedAdminCorrelatedRestResourceDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDelegatedAdminCorrelatedRestResourceResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedDelegatedAdminCorrelatedRestResourceAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccDelegatedAdminCorrelatedRestResourceResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedDelegatedAdminCorrelatedRestResourceAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccDelegatedAdminCorrelatedRestResourceResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_delegated_admin_correlated_rest_resource." + resourceName,
				ImportStateId:     updatedResourceModel.restResourceTypeName + "/" + updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccDelegatedAdminCorrelatedRestResourceResource(resourceName string, resourceModel delegatedAdminCorrelatedRestResourceTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_rest_resource_type" "%[3]s" {
  type                        = "user"
  id                          = "%[3]s"
  enabled                     = true
  resource_endpoint           = "userRestResourceDelegatedAdminCorrelatedRestResourceTest"
  structural_ldap_objectclass = "inetOrgPerson"
  search_base_dn              = "cn=users,dc=test,dc=com"
}
resource "pingdirectory_delegated_admin_correlated_rest_resource" "%[1]s" {
  id                                            = "%[2]s"
  rest_resource_type_name                       = pingdirectory_rest_resource_type.%[3]s.id
  display_name                                  = "%[4]s"
  correlated_rest_resource                      = "%[5]s"
  primary_rest_resource_correlation_attribute   = "%[6]s"
  secondary_rest_resource_correlation_attribute = "%[7]s"
}`, resourceName,
		resourceModel.id,
		resourceModel.restResourceTypeName,
		resourceModel.displayName,
		resourceModel.correlatedRestResource,
		resourceModel.primaryRestResourceCorrelationAttribute,
		resourceModel.secondaryRestResourceCorrelationAttribute)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDelegatedAdminCorrelatedRestResourceAttributes(config delegatedAdminCorrelatedRestResourceTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.DelegatedAdminCorrelatedRestResourceApi.GetDelegatedAdminCorrelatedRestResource(ctx, config.id, config.restResourceTypeName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Delegated Admin Correlated Rest Resource"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "display-name",
			config.displayName, response.DisplayName)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "correlated-rest-resource",
			config.correlatedRestResource, response.CorrelatedRESTResource)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "primary-rest-resource-correlation-attribute",
			config.primaryRestResourceCorrelationAttribute, response.PrimaryRESTResourceCorrelationAttribute)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "secondary-rest-resource-correlation-attribute",
			config.secondaryRestResourceCorrelationAttribute, response.SecondaryRESTResourceCorrelationAttribute)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckDelegatedAdminCorrelatedRestResourceDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.DelegatedAdminCorrelatedRestResourceApi.GetDelegatedAdminCorrelatedRestResource(ctx, testIdDelegatedAdminCorrelatedRestResource, testRestResourceTypeName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Delegated Admin Correlated Rest Resource", testIdDelegatedAdminCorrelatedRestResource)
	}
	return nil
}
