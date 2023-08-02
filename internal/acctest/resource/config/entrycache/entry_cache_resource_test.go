package entrycache_test

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

const testIdFifoEntryCache = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type fifoEntryCacheTestModel struct {
	id         string
	enabled    bool
	cacheLevel int64
}

func TestAccFifoEntryCache(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := fifoEntryCacheTestModel{
		id:         testIdFifoEntryCache,
		enabled:    true,
		cacheLevel: 1,
	}
	updatedResourceModel := fifoEntryCacheTestModel{
		id:         testIdFifoEntryCache,
		enabled:    false,
		cacheLevel: 2,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckFifoEntryCacheDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccFifoEntryCacheResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedFifoEntryCacheAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_entry_cache.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_entry_cache.%s", resourceName), "cache_level", strconv.FormatInt(initialResourceModel.cacheLevel, 10)),
					resource.TestCheckResourceAttrSet("data.pingdirectory_entry_caches.list", "objects.0.id"),
				),
			},
			{
				// Test updating some fields
				Config: testAccFifoEntryCacheResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedFifoEntryCacheAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccFifoEntryCacheResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_entry_cache." + resourceName,
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

func testAccFifoEntryCacheResource(resourceName string, resourceModel fifoEntryCacheTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_entry_cache" "%[1]s" {
  id          = "%[2]s"
  enabled     = %[3]t
  cache_level = %[4]d
}

data "pingdirectory_entry_cache" "%[1]s" {
  id = "%[2]s"
  depends_on = [
    pingdirectory_entry_cache.%[1]s
  ]
}

data "pingdirectory_entry_caches" "list" {
  depends_on = [
    pingdirectory_entry_cache.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.enabled,
		resourceModel.cacheLevel)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedFifoEntryCacheAttributes(config fifoEntryCacheTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.EntryCacheApi.GetEntryCache(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Fifo Entry Cache"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "cache-level",
			config.cacheLevel, response.CacheLevel)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckFifoEntryCacheDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.EntryCacheApi.GetEntryCache(ctx, testIdFifoEntryCache).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Fifo Entry Cache", testIdFifoEntryCache)
	}
	return nil
}
