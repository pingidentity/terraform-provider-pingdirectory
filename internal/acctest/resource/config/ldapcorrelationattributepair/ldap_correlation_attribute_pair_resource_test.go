package ldapcorrelationattributepair_test

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

const testIdLdapCorrelationAttributePair = "MyId"
const testCorrelatedLdapDataViewName = "MyCorrelatedLdapDataView"
const testScimResourceTypeNameTest = "MyScimResourceType"

// Attributes to test with. Add optional properties to test here if desired.
type ldapCorrelationAttributePairTestModel struct {
	id                            string
	correlatedLdapDataViewName    string
	scimResourceTypeName          string
	primaryCorrelationAttribute   string
	secondaryCorrelationAttribute string
}

func TestAccLdapCorrelationAttributePair(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := ldapCorrelationAttributePairTestModel{
		id:                            testIdLdapCorrelationAttributePair,
		correlatedLdapDataViewName:    testCorrelatedLdapDataViewName,
		scimResourceTypeName:          testScimResourceTypeNameTest,
		primaryCorrelationAttribute:   "cn",
		secondaryCorrelationAttribute: "cn",
	}
	updatedResourceModel := ldapCorrelationAttributePairTestModel{
		id:                            testIdLdapCorrelationAttributePair,
		correlatedLdapDataViewName:    testCorrelatedLdapDataViewName,
		scimResourceTypeName:          testScimResourceTypeNameTest,
		primaryCorrelationAttribute:   "sn",
		secondaryCorrelationAttribute: "sn",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckLdapCorrelationAttributePairDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccLdapCorrelationAttributePairResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedLdapCorrelationAttributePairAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_ldap_correlation_attribute_pair.%s", resourceName), "primary_correlation_attribute", initialResourceModel.primaryCorrelationAttribute),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_ldap_correlation_attribute_pair.%s", resourceName), "secondary_correlation_attribute", initialResourceModel.secondaryCorrelationAttribute),
					resource.TestCheckResourceAttrSet("data.pingdirectory_ldap_correlation_attribute_pairs.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccLdapCorrelationAttributePairResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedLdapCorrelationAttributePairAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccLdapCorrelationAttributePairResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_ldap_correlation_attribute_pair." + resourceName,
				ImportStateId:     updatedResourceModel.scimResourceTypeName + "/" + updatedResourceModel.correlatedLdapDataViewName + "/" + updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.LdapCorrelationAttributePairAPI.DeleteLdapCorrelationAttributePair(ctx, updatedResourceModel.id, updatedResourceModel.correlatedLdapDataViewName, updatedResourceModel.scimResourceTypeName).Execute()
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

func testAccLdapCorrelationAttributePairResource(resourceName string, resourceModel ldapCorrelationAttributePairTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_scim_resource_type" "%[4]s" {
  type        = "ldap-mapping"
  name        = "%[4]s"
  core_schema = pingdirectory_scim_schema.myScimSchema.schema_urn
  enabled     = false
  endpoint    = "myendpoint"
}

resource "pingdirectory_scim_schema" "myScimSchema" {
  schema_urn = "urn:com:example:ldapcorrelationattributepairtest"
}

resource "pingdirectory_correlated_ldap_data_view" "%[3]s" {
  name                            = "%[3]s"
  scim_resource_type_name         = pingdirectory_scim_resource_type.%[4]s.id
  structural_ldap_objectclass     = "ldapObject"
  include_base_dn                 = "cn=com.example"
  primary_correlation_attribute   = "cn"
  secondary_correlation_attribute = "cn"
}

resource "pingdirectory_ldap_correlation_attribute_pair" "%[1]s" {
  name                            = "%[2]s"
  correlated_ldap_data_view_name  = pingdirectory_correlated_ldap_data_view.%[3]s.id
  scim_resource_type_name         = pingdirectory_scim_resource_type.%[4]s.id
  primary_correlation_attribute   = "%[5]s"
  secondary_correlation_attribute = "%[6]s"
}

data "pingdirectory_ldap_correlation_attribute_pair" "%[1]s" {
  name                           = "%[2]s"
  correlated_ldap_data_view_name = "%[3]s"
  scim_resource_type_name        = "%[4]s"
  depends_on = [
    pingdirectory_ldap_correlation_attribute_pair.%[1]s
  ]
}

data "pingdirectory_ldap_correlation_attribute_pairs" "list" {
  correlated_ldap_data_view_name = "%[3]s"
  scim_resource_type_name        = "%[4]s"
  depends_on = [
    pingdirectory_ldap_correlation_attribute_pair.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.correlatedLdapDataViewName,
		resourceModel.scimResourceTypeName,
		resourceModel.primaryCorrelationAttribute,
		resourceModel.secondaryCorrelationAttribute)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedLdapCorrelationAttributePairAttributes(config ldapCorrelationAttributePairTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LdapCorrelationAttributePairAPI.GetLdapCorrelationAttributePair(ctx, config.id, config.correlatedLdapDataViewName, config.scimResourceTypeName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Ldap Correlation Attribute Pair"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "primary-correlation-attribute",
			config.primaryCorrelationAttribute, response.PrimaryCorrelationAttribute)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "secondary-correlation-attribute",
			config.secondaryCorrelationAttribute, response.SecondaryCorrelationAttribute)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckLdapCorrelationAttributePairDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.LdapCorrelationAttributePairAPI.GetLdapCorrelationAttributePair(ctx, testIdLdapCorrelationAttributePair, testCorrelatedLdapDataViewName, testScimResourceTypeNameTest).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Ldap Correlation Attribute Pair", testIdLdapCorrelationAttributePair)
	}
	return nil
}
