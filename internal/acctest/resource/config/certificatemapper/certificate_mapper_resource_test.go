package certificatemapper_test

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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedSubjectEqualsDnCertificateMapperAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_certificate_mapper.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_certificate_mappers.list", "objects.0.id"),
				),
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
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.CertificateMapperApi.DeleteCertificateMapper(ctx, updatedResourceModel.id).Execute()
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

func testAccSubjectEqualsDnCertificateMapperResource(resourceName string, resourceModel subjectEqualsDnCertificateMapperTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_certificate_mapper" "%[1]s" {
  type    = "subject-equals-dn"
  name    = "%[2]s"
  enabled = %[3]t
}

data "pingdirectory_certificate_mapper" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_certificate_mapper.%[1]s
  ]
}

data "pingdirectory_certificate_mappers" "list" {
  depends_on = [
    pingdirectory_certificate_mapper.%[1]s
  ]
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
