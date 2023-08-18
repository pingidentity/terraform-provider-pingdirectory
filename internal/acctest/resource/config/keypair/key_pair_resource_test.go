package keypair_test

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

const testIdKeyPair = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type keyPairTestModel struct {
	id        string
	subjectDn string
}

func TestAccKeyPair(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := keyPairTestModel{
		id:        testIdKeyPair,
		subjectDn: "cn=Directory Server,O=Ping Identity Key Pair",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccKeyPairResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedKeyPairAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_key_pair.%s", resourceName), "subject_dn", initialResourceModel.subjectDn),
					resource.TestCheckResourceAttrSet("data.pingdirectory_key_pairs.list", "ids.0"),
				),
			},
			{
				// Test importing the resource
				Config:            testAccKeyPairResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_key_pair." + resourceName,
				ImportStateId:     initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
					"private_key",
				},
			},
			{
				// Test plan after removing config on PD
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.KeyPairApi.DeleteKeyPair(ctx, initialResourceModel.id).Execute()
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

func testAccKeyPairResource(resourceName string, resourceModel keyPairTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_key_pair" "%[1]s" {
  name       = "%[2]s"
  subject_dn = "%[3]s"
}

data "pingdirectory_key_pair" "%[1]s" {
  name = "%[2]s"
  depends_on = [
    pingdirectory_key_pair.%[1]s
  ]
}

data "pingdirectory_key_pairs" "list" {
  depends_on = [
    pingdirectory_key_pair.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.subjectDn)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedKeyPairAttributes(config keyPairTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "Key Pair"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.KeyPairApi.GetKeyPair(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "subject-dn", config.subjectDn, response.SubjectDN)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckKeyPairDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.KeyPairApi.GetKeyPair(ctx, testIdKeyPair).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Key Pair", testIdKeyPair)
	}
	return nil
}
