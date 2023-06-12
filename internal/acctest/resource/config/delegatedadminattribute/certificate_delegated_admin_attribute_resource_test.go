package delegatedadminattribute_test

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

const testCertificateDAAttributeType = "cn"
const testCertificateDARestResourceTypeName = "myParentCertficateRestResource"

// Attributes to test with. Add optional properties to test here if desired.
type certificateDelegatedAdminAttributeTestModel struct {
	restResourceTypeName string
	attributeType        string
	displayName          string
}

func TestAccCertificateDelegatedAdminAttribute(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := certificateDelegatedAdminAttributeTestModel{
		restResourceTypeName: testCertificateDARestResourceTypeName,
		attributeType:        testCertificateDAAttributeType,
		displayName:          "myname",
	}
	updatedResourceModel := certificateDelegatedAdminAttributeTestModel{
		restResourceTypeName: testCertificateDARestResourceTypeName,
		attributeType:        testCertificateDAAttributeType,
		displayName:          "myupdatedname",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckCertificateDelegatedAdminAttributeDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccCertificateDelegatedAdminAttributeResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedCertificateDelegatedAdminAttributeAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccCertificateDelegatedAdminAttributeResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedCertificateDelegatedAdminAttributeAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccCertificateDelegatedAdminAttributeResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_delegated_admin_attribute." + resourceName,
				ImportStateId:     updatedResourceModel.restResourceTypeName + "/" + updatedResourceModel.attributeType,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccCertificateDelegatedAdminAttributeResource(resourceName string, resourceModel certificateDelegatedAdminAttributeTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_rest_resource_type" "%[2]s" {
	type = "generic"
  id                          = "%[2]s"
  enabled                     = false
  resource_endpoint           = "myendpoint"
  structural_ldap_objectclass = "device"
  search_base_dn              = "dc=example,dc=com"
}

resource "pingdirectory_delegated_admin_attribute" "%[1]s" {
	type = "certificate"
  rest_resource_type_name = pingdirectory_rest_resource_type.%[2]s.id
  attribute_type          = "%[3]s"
  display_name            = "%[4]s"
}`, resourceName,
		resourceModel.restResourceTypeName,
		resourceModel.attributeType,
		resourceModel.displayName)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedCertificateDelegatedAdminAttributeAttributes(config certificateDelegatedAdminAttributeTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.DelegatedAdminAttributeApi.GetDelegatedAdminAttribute(ctx, config.attributeType, config.restResourceTypeName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Certificate Delegated Admin Attribute"
		err = acctest.TestAttributesMatchString(resourceType, &config.attributeType, "attribute-type",
			config.attributeType, response.CertificateDelegatedAdminAttributeResponse.AttributeType)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.attributeType, "display-name",
			config.displayName, response.CertificateDelegatedAdminAttributeResponse.DisplayName)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckCertificateDelegatedAdminAttributeDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.DelegatedAdminAttributeApi.GetDelegatedAdminAttribute(ctx, testCertificateDAAttributeType, testCertificateDARestResourceTypeName).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Certificate Delegated Admin Attribute", testCertificateDAAttributeType)
	}
	return nil
}
