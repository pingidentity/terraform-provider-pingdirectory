package config_test

import (
	"fmt"
	"terraform-provider-pingdirectory/internal/acctest"
	"terraform-provider-pingdirectory/internal/provider"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLocation(t *testing.T) {
	resourceName := "TestLocation"
	locationName := "Hoenn"
	locationDescription := "Home of Kyogre"
	updatedDescription := "Home of Groudon"
	updatedName := "Hoennn"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				// Test basic resource
				Config: testAccLocationResource(resourceName, locationName, locationDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_location.%s", resourceName), "name", locationName),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_location.%s", resourceName), "description", locationDescription),
				),
			},
			{
				// Test updating the description
				Config: testAccLocationResource(resourceName, locationName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_location.%s", resourceName), "name", locationName),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_location.%s", resourceName), "description", updatedDescription),
				),
			},
			{
				// Test removing the description
				Config: testAccLocationResourceNoDescription(resourceName, locationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_location.%s", resourceName), "name", locationName),
					resource.TestCheckNoResourceAttr(fmt.Sprintf("pingdirectory_location.%s", resourceName), "description"),
				),
			},
			{
				// Test updating the name
				Config: testAccLocationResource(resourceName, updatedName, locationDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_location.%s", resourceName), "name", updatedName),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingdirectory_location.%s", resourceName), "description", locationDescription),
				),
			},
		},
	})
}

func testAccLocationResource(resourceName, locationName, description string) string {
	return fmt.Sprintf(`
resource "pingdirectory_location" "%[1]s" {
	name = "%[2]s"
	description = "%[3]s"
}`, resourceName, locationName, description)
}

func testAccLocationResourceNoDescription(resourceName, locationName string) string {
	return fmt.Sprintf(`
resource "pingdirectory_location" "%[1]s" {
	name = "%[2]s"
}`, resourceName, locationName)
}
