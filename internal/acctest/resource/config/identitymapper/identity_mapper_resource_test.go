package identitymapper_test

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

const testIdExactMatchIdentityMapper = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type exactMatchIdentityMapperTestModel struct {
	id             string
	matchAttribute []string
	enabled        bool
}

func TestAccExactMatchIdentityMapper(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := exactMatchIdentityMapperTestModel{
		id:             testIdExactMatchIdentityMapper,
		matchAttribute: []string{"uid"},
		enabled:        true,
	}
	updatedResourceModel := exactMatchIdentityMapperTestModel{
		id:             testIdExactMatchIdentityMapper,
		matchAttribute: []string{"uid"},
		enabled:        false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckExactMatchIdentityMapperDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccExactMatchIdentityMapperResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedExactMatchIdentityMapperAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccExactMatchIdentityMapperResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedExactMatchIdentityMapperAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccExactMatchIdentityMapperResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_identity_mapper." + resourceName,
				ImportStateId:           updatedResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccExactMatchIdentityMapperResource(resourceName string, resourceModel exactMatchIdentityMapperTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_identity_mapper" "%[1]s" {
  type            = "exact-match"
  id              = "%[2]s"
  match_attribute = %[3]s
  enabled         = %[4]t
}`, resourceName, resourceModel.id,
		acctest.StringSliceToTerraformString(resourceModel.matchAttribute),
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedExactMatchIdentityMapperAttributes(config exactMatchIdentityMapperTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.IdentityMapperApi.GetIdentityMapper(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Exact Match Identity Mapper"
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "match-attribute",
			config.matchAttribute, response.ExactMatchIdentityMapperResponse.MatchAttribute)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.ExactMatchIdentityMapperResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckExactMatchIdentityMapperDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.IdentityMapperApi.GetIdentityMapper(ctx, testIdExactMatchIdentityMapper).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Exact Match Identity Mapper", testIdExactMatchIdentityMapper)
	}
	return nil
}
