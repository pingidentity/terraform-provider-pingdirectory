package workqueue_test

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

//const testIdHighThroughputWorkQueue = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type highThroughputWorkQueueTestModel struct {
	num_worker_threads      int64
	max_work_queue_capacity int64
}

func TestAccHighThroughputWorkQueue(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := highThroughputWorkQueueTestModel{
		num_worker_threads:      3,
		max_work_queue_capacity: 800,
	}
	// set back to default values for other tests
	updatedResourceModel := highThroughputWorkQueueTestModel{
		num_worker_threads:      0,
		max_work_queue_capacity: 100,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccHighThroughputWorkQueueResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedHighThroughputWorkQueueAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccHighThroughputWorkQueueResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedHighThroughputWorkQueueAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccHighThroughputWorkQueueResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_default_high_throughput_work_queue." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccHighThroughputWorkQueueResource(resourceName string, resourceModel highThroughputWorkQueueTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_high_throughput_work_queue" "%[1]s" {
  num_worker_threads      = %[2]d
  max_work_queue_capacity = %[3]d
}`, resourceName,
		resourceModel.num_worker_threads,
		resourceModel.max_work_queue_capacity)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedHighThroughputWorkQueueAttributes(config highThroughputWorkQueueTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "workqueue"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.WorkQueueApi.GetWorkQueue(ctx).Execute()
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, nil, "num-worker-threads", config.num_worker_threads, *response.NumWorkerThreads)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, nil, "max-work-queue-capacity", config.max_work_queue_capacity, *response.MaxWorkQueueCapacity)
		if err != nil {
			return err
		}
		return nil
	}
}
