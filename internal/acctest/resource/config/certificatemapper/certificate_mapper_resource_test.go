package certificatemapper_test

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

const testIdSubjectEqualsDnCertificateMapper = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type subjectEqualsDnCertificateMapperTestModel struct {
	id      string
	enabled bool
}

func TestAccSubjectEqualsDnCertificateMapper(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := subjectEqualsDnCertificateMapperTestModel{
		id:      testIdSubjectEqualsDnCertificateMapper,
		enabled: true,
	}
	updatedResourceModel := subjectEqualsDnCertificateMapperTestModel{
		id:      testIdSubjectEqualsDnCertificateMapper,
		enabled: false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckSubjectEqualsDnCertificateMapperDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSubjectEqualsDnCertificateMapperResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSubjectEqualsDnCertificateMapperAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccSubjectEqualsDnCertificateMapperResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSubjectEqualsDnCertificateMapperAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSubjectEqualsDnCertificateMapperResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_certificate_mapper." + resourceName,
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

func testAccSubjectEqualsDnCertificateMapperResource(resourceName string, resourceModel subjectEqualsDnCertificateMapperTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_certificate_mapper" "%[1]s" {
  type    = "subject-equals-dn"
  id      = "%[2]s"
  enabled = %[3]t
}`, resourceName,
		resourceModel.id,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSubjectEqualsDnCertificateMapperAttributes(config subjectEqualsDnCertificateMapperTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.CertificateMapperApi.GetCertificateMapper(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Subject Equals Dn Certificate Mapper"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.SubjectEqualsDnCertificateMapperResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSubjectEqualsDnCertificateMapperDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.CertificateMapperApi.GetCertificateMapper(ctx, testIdSubjectEqualsDnCertificateMapper).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Subject Equals Dn Certificate Mapper", testIdSubjectEqualsDnCertificateMapper)
	}
	return nil
}
