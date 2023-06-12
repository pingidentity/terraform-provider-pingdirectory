package datasecurityauditor_test

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

const testIdExpiredPasswordDataSecurityAuditor = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type expiredPasswordDataSecurityAuditorTestModel struct {
	id string
}

func TestAccExpiredPasswordDataSecurityAuditor(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := expiredPasswordDataSecurityAuditorTestModel{
		id: testIdExpiredPasswordDataSecurityAuditor,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckExpiredPasswordDataSecurityAuditorDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccExpiredPasswordDataSecurityAuditorResource(resourceName, initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccExpiredPasswordDataSecurityAuditorResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_data_security_auditor." + resourceName,
				ImportStateId:     initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccExpiredPasswordDataSecurityAuditorResource(resourceName string, resourceModel expiredPasswordDataSecurityAuditorTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_data_security_auditor" "%[1]s" {
	type = "expired-password"
  id = "%[2]s"
}`, resourceName,
		resourceModel.id)
}

// Test that any objects created by the test are destroyed
func testAccCheckExpiredPasswordDataSecurityAuditorDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.DataSecurityAuditorApi.GetDataSecurityAuditor(ctx, testIdExpiredPasswordDataSecurityAuditor).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Expired Password Data Security Auditor", testIdExpiredPasswordDataSecurityAuditor)
	}
	return nil
}
