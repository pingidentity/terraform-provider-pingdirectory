package delegatedadminattributecategory_test

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

const testIdDelegatedAdminAttributeCategory = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type delegatedAdminAttributeCategoryTestModel struct {
	displayName       string
	displayOrderIndex int64
}

func TestAccDelegatedAdminAttributeCategory(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := delegatedAdminAttributeCategoryTestModel{
		displayName:       testIdDelegatedAdminAttributeCategory,
		displayOrderIndex: 0,
	}
	updatedResourceModel := delegatedAdminAttributeCategoryTestModel{
		displayName:       testIdDelegatedAdminAttributeCategory,
		displayOrderIndex: 2,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckDelegatedAdminAttributeCategoryDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDelegatedAdminAttributeCategoryResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedDelegatedAdminAttributeCategoryAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccDelegatedAdminAttributeCategoryResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedDelegatedAdminAttributeCategoryAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccDelegatedAdminAttributeCategoryResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_delegated_admin_attribute_category." + resourceName,
				ImportStateId:     updatedResourceModel.displayName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccDelegatedAdminAttributeCategoryResource(resourceName string, resourceModel delegatedAdminAttributeCategoryTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_delegated_admin_attribute_category" "%[1]s" {
  display_name        = "%[2]s"
  display_order_index = %[3]d
}`, resourceName,
		resourceModel.displayName,
		resourceModel.displayOrderIndex)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDelegatedAdminAttributeCategoryAttributes(config delegatedAdminAttributeCategoryTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.DelegatedAdminAttributeCategoryApi.GetDelegatedAdminAttributeCategory(ctx, config.displayName).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Delegated Admin Attribute Category"
		err = acctest.TestAttributesMatchString(resourceType, &config.displayName, "display-name",
			config.displayName, response.DisplayName)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.displayName, "display-order-index",
			config.displayOrderIndex, response.DisplayOrderIndex)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckDelegatedAdminAttributeCategoryDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.DelegatedAdminAttributeCategoryApi.GetDelegatedAdminAttributeCategory(ctx, testIdDelegatedAdminAttributeCategory).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Delegated Admin Attribute Category", testIdDelegatedAdminAttributeCategory)
	}
	return nil
}
