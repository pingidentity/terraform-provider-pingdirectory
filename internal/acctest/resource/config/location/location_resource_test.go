package location_test

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

const locationName = "Hoenn"
const updatedLocationName = "Hoennn"

func TestAccLocation(t *testing.T) {
	resourceName := "TestLocation"
	locationDescription := "Home of Kyogre"
	updatedDescription := "Home of Groudon"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckLocationDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource
				Config: testAccLocationResource(resourceName, locationName, locationDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedLocationAttributes(locationName, locationDescription),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_location.%s", resourceName), "id", locationName),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_location.%s", resourceName), "description", locationDescription),
					resource.TestCheckResourceAttrSet("data.pingdirectory_locations.list", "ids.0"),
				),
			},
			{
				// Test updating the description
				Config: testAccLocationResource(resourceName, locationName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedLocationAttributes(locationName, updatedDescription),
				),
			},
			{
				// Test removing the description
				Config: testAccLocationResourceNoDescription(resourceName, locationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(fmt.Sprintf("pingdirectory_location.%s", resourceName), "description"),
					testAccCheckExpectedLocationAttributes(locationName, ""),
				),
			},
			{
				// Test updating the name
				Config: testAccLocationResource(resourceName, updatedLocationName, locationDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedLocationAttributes(updatedLocationName, locationDescription),
				),
			},
			{
				// Test importing the default location, which should not have a description attribute
				Config:            testAccLocationResource(resourceName, updatedLocationName, locationDescription),
				ResourceName:      "pingdirectory_location." + resourceName,
				ImportStateId:     updatedLocationName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.LocationAPI.DeleteLocation(ctx, updatedLocationName).Execute()
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

func testAccLocationResource(resourceName, locationName, description string) string {
	return fmt.Sprintf(`
resource "pingdirectory_location" "%[1]s" {
  name        = "%[2]s"
  description = "%[3]s"
}

data "pingdirectory_location" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_location.%[1]s
  ]
}

data "pingdirectory_locations" "list" {
  depends_on = [
    pingdirectory_location.%[1]s
  ]
}`, resourceName, locationName, description)
}

func testAccLocationResourceNoDescription(resourceName, locationName string) string {
	return fmt.Sprintf(`
resource "pingdirectory_location" "%[1]s" {
  name = "%[2]s"
}`, resourceName, locationName)
}

// Test that any locations created by the test are destroyed
func testAccCheckLocationDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	// Check for location names used in this test
	names := []string{locationName, updatedLocationName}
	for _, name := range names {
		_, _, err := testClient.LocationAPI.GetLocation(ctx, name).Execute()
		if err == nil {
			return acctest.ExpectedDestroyError("location", name)
		}
	}
	return nil
}

// Test that the expected location attributes are set on the PingDirectory server
func testAccCheckExpectedLocationAttributes(name, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "location"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		locationResponse, _, err := testClient.LocationAPI.GetLocation(ctx, name).Execute()
		if err != nil {
			return err
		}
		// Verify that description matches expected
		err = acctest.TestAttributesMatchStringPointer(resourceType, &name, "description", description, locationResponse.Description)
		if err != nil {
			return err
		}
		return nil
	}
}
