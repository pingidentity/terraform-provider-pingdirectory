// Copyright © 2025 Ping Identity Corporation

package rootdsebackend_test

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

// Attributes to test with. Add optional properties to test here if desired.
type rootDseBackendTestModel struct {
	showAllAttributes bool
}

func TestAccRootDseBackend(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := rootDseBackendTestModel{
		showAllAttributes: true,
	}
	updatedResourceModel := rootDseBackendTestModel{
		showAllAttributes: false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccRootDseBackendResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedRootDseBackendAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_root_dse_backend.%s", resourceName), "show_all_attributes", strconv.FormatBool(initialResourceModel.showAllAttributes)),
				),
			},
			{
				// Test updating some fields
				Config: testAccRootDseBackendResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedRootDseBackendAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccRootDseBackendResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_default_root_dse_backend." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRootDseBackendResource(resourceName string, resourceModel rootDseBackendTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_root_dse_backend" "%[1]s" {
  show_all_attributes = %[2]t
}

data "pingdirectory_root_dse_backend" "%[1]s" {
  depends_on = [
    pingdirectory_default_root_dse_backend.%[1]s
  ]
}`, resourceName,
		resourceModel.showAllAttributes)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedRootDseBackendAttributes(config rootDseBackendTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.RootDseBackendAPI.GetRootDseBackend(ctx).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Root Dse Backend"
		err = acctest.TestAttributesMatchBool(resourceType, nil, "show-all-attributes",
			config.showAllAttributes, response.ShowAllAttributes)
		if err != nil {
			return err
		}
		return nil
	}
}
