package correlatedldapdataview_test

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

const testIdCorrelatedLdapDataView = "MyCorrelatedLdapDataViewId"
const testScimResourceTypeNameCorrelated = "MyScimResourceTypeNameCorrelated"

// Attributes to test with. Add optional properties to test here if desired.
type correlatedLdapDataViewTestModel struct {
	id                            string
	scimResourceTypeName          string
	structuralLdapObjectclass     string
	includeBaseDn                 string
	primaryCorrelationAttribute   string
	secondaryCorrelationAttribute string
}

func TestAccCorrelatedLdapDataView(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := correlatedLdapDataViewTestModel{
		id:                            testIdCorrelatedLdapDataView,
		scimResourceTypeName:          testScimResourceTypeNameCorrelated,
		structuralLdapObjectclass:     "ldapObject",
		includeBaseDn:                 "cn=com.company",
		primaryCorrelationAttribute:   "sn",
		secondaryCorrelationAttribute: "sn",
	}
	updatedResourceModel := correlatedLdapDataViewTestModel{
		id:                            testIdCorrelatedLdapDataView,
		scimResourceTypeName:          testScimResourceTypeNameCorrelated,
		structuralLdapObjectclass:     "ldapObject",
		includeBaseDn:                 "cn=com.example",
		primaryCorrelationAttribute:   "cn",
		secondaryCorrelationAttribute: "cn",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckCorrelatedLdapDataViewDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccCorrelatedLdapDataViewResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedCorrelatedLdapDataViewAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_correlated_ldap_data_view.%s", resourceName), "structural_ldap_objectclass", initialResourceModel.structuralLdapObjectclass),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_correlated_ldap_data_view.%s", resourceName), "include_base_dn", initialResourceModel.includeBaseDn),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_correlated_ldap_data_view.%s", resourceName), "primary_correlation_attribute", initialResourceModel.primaryCorrelationAttribute),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_correlated_ldap_data_view.%s", resourceName), "secondary_correlation_attribute", initialResourceModel.secondaryCorrelationAttribute),
					resource.TestCheckResourceAttrSet("data.pingdirectory_correlated_ldap_data_views.list", "ids.0"),
				),
			},
			{
				// Test updating some fields
				Config: testAccCorrelatedLdapDataViewResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedCorrelatedLdapDataViewAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccCorrelatedLdapDataViewResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_correlated_ldap_data_view." + resourceName,
				ImportStateId:     updatedResourceModel.scimResourceTypeName + "/" + updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.CorrelatedLdapDataViewAPI.DeleteCorrelatedLdapDataView(ctx, updatedResourceModel.id, updatedResourceModel.scimResourceTypeName).Execute()
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

func testAccCorrelatedLdapDataViewResource(resourceName string, resourceModel correlatedLdapDataViewTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_scim_resource_type" "%[3]s" {
  type        = "ldap-mapping"
  name        = "%[3]s"
  core_schema = pingdirectory_scim_schema.myScimSchema.schema_urn
  enabled     = false
  endpoint    = "myendpoint"
}

resource "pingdirectory_scim_schema" "myScimSchema" {
  schema_urn = "urn:com:example:correlatedldapdataviewtest"
}

resource "pingdirectory_correlated_ldap_data_view" "%[1]s" {
  name                            = "%[2]s"
  scim_resource_type_name         = pingdirectory_scim_resource_type.%[3]s.id
  structural_ldap_objectclass     = "%[4]s"
  include_base_dn                 = "%[5]s"
  primary_correlation_attribute   = "%[6]s"
  secondary_correlation_attribute = "%[7]s"
}

data "pingdirectory_correlated_ldap_data_view" "%[1]s" {
  name                    = "%[2]s"
  scim_resource_type_name = "%[3]s"
  depends_on = [
    pingdirectory_correlated_ldap_data_view.%[1]s
  ]
}

data "pingdirectory_correlated_ldap_data_views" "list" {
  scim_resource_type_name = "%[3]s"
  depends_on = [
    pingdirectory_correlated_ldap_data_view.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.scimResourceTypeName,
		resourceModel.structuralLdapObjectclass,
		resourceModel.includeBaseDn,
		resourceModel.primaryCorrelationAttribute,
		resourceModel.secondaryCorrelationAttribute)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedCorrelatedLdapDataViewAttributes(config correlatedLdapDataViewTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.CorrelatedLdapDataViewAPI.GetCorrelatedLdapDataView(ctx, config.id, config.scimResourceTypeName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Correlated Ldap Data View"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "structural-ldap-objectclass",
			config.structuralLdapObjectclass, response.StructuralLDAPObjectclass)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "include-base-dn",
			config.includeBaseDn, response.IncludeBaseDN)
		if err != nil {
			return err
		}
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
func testAccCheckCorrelatedLdapDataViewDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.CorrelatedLdapDataViewAPI.GetCorrelatedLdapDataView(ctx, testIdCorrelatedLdapDataView, testScimResourceTypeNameCorrelated).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Correlated Ldap Data View", testIdCorrelatedLdapDataView)
	}
	return nil
}
