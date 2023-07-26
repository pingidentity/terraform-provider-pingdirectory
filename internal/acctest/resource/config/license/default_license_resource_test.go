package license_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdLicense = "MyId"

func TestAccLicense(t *testing.T) {
	resourceName := "myresource"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccLicenseResource(resourceName),
			},
			{
				// Test importing the resource
				Config:            testAccLicenseResource(resourceName),
				ResourceName:      "pingdirectory_default_license." + resourceName,
				ImportStateId:     testIdLicense,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccLicenseResource(resourceName string) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_license" "%[1]s" {
}

data "pingdirectory_license" "%[1]s" {
  depends_on = [
    pingdirectory_default_license.%[1]s
  ]
}`, resourceName)
}
