package config_test

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

const testIdDnMap = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type dnMapTestModel struct {
	id            string
	fromDnPattern string
	toDnPattern   string
}

func TestAccDnMap(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := dnMapTestModel{
		id:            testIdDnMap,
		fromDnPattern: "*,**,dc=com",
		toDnPattern:   "uid={givenname:/^(.)(.*)/$1/s}{sn:/^(.)(.*)/$1/s}{eid},{2},o=example",
	}
	updatedResourceModel := dnMapTestModel{
		id:            testIdDnMap,
		fromDnPattern: "*,**,dc=edu",
		toDnPattern:   "uid={givenname:/^(.)(.*)/$1/s}{sn:/^(.)(.*)/$1/s}{eid},{2},o=mycorp",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckDnMapDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDnMapResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedDnMapAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccDnMapResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedDnMapAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccDnMapResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_dn_map." + resourceName,
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

func testAccDnMapResource(resourceName string, resourceModel dnMapTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_dn_map" "%[1]s" {
  id              = "%[2]s"
  from_dn_pattern = "%[3]s"
  to_dn_pattern   = "%[4]s"
}`, resourceName,
		resourceModel.id,
		resourceModel.fromDnPattern,
		resourceModel.toDnPattern)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedDnMapAttributes(config dnMapTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.DnMapApi.GetDnMap(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Dn Map"
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "from-dn-pattern",
			config.fromDnPattern, response.FromDNPattern)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "to-dn-pattern",
			config.toDnPattern, response.ToDNPattern)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckDnMapDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.DnMapApi.GetDnMap(ctx, testIdDnMap).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Dn Map", testIdDnMap)
	}
	return nil
}
