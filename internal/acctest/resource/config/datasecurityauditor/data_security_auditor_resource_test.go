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
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckExpiredPasswordDataSecurityAuditorDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccExpiredPasswordDataSecurityAuditorResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_data_security_auditor.%s", resourceName), "type", "expired-password"),
					resource.TestCheckResourceAttrSet("data.pingdirectory_data_security_auditors.list", "objects.0.id"),
				),
			},
			{
				// Test importing the resource
				Config:            testAccExpiredPasswordDataSecurityAuditorResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_data_security_auditor." + resourceName,
				ImportStateId:     initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.DataSecurityAuditorAPI.DeleteDataSecurityAuditor(ctx, initialResourceModel.id).Execute()
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

func testAccExpiredPasswordDataSecurityAuditorResource(resourceName string, resourceModel expiredPasswordDataSecurityAuditorTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_data_security_auditor" "%[1]s" {
  type = "expired-password"
  name = "%[2]s"
}

data "pingdirectory_data_security_auditor" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_data_security_auditor.%[1]s
  ]
}

data "pingdirectory_data_security_auditors" "list" {
  depends_on = [
    pingdirectory_data_security_auditor.%[1]s
  ]
}`, resourceName,
		resourceModel.id)
}

// Test that any objects created by the test are destroyed
func testAccCheckExpiredPasswordDataSecurityAuditorDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.DataSecurityAuditorAPI.GetDataSecurityAuditor(ctx, testIdExpiredPasswordDataSecurityAuditor).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Expired Password Data Security Auditor", testIdExpiredPasswordDataSecurityAuditor)
	}
	return nil
}
