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

const testIdHttpServletCrossOriginPolicy = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type httpServletCrossOriginPolicyTestModel struct {
	id                 string
	corsAllowedHeaders []string
}

func TestAccHttpServletCrossOriginPolicy(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := httpServletCrossOriginPolicyTestModel{
		id:                 testIdHttpServletCrossOriginPolicy,
		corsAllowedHeaders: []string{"Accept, Access-Control-Request-Headers"},
	}
	updatedResourceModel := httpServletCrossOriginPolicyTestModel{
		id:                 testIdHttpServletCrossOriginPolicy,
		corsAllowedHeaders: []string{"Accept"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckHttpServletCrossOriginPolicyDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccHttpServletCrossOriginPolicyResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedHttpServletCrossOriginPolicyAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccHttpServletCrossOriginPolicyResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedHttpServletCrossOriginPolicyAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccHttpServletCrossOriginPolicyResource(resourceName, initialResourceModel),
				ResourceName:            "pingdirectory_http_servlet_cross_origin_policy." + resourceName,
				ImportStateId:           initialResourceModel.id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccHttpServletCrossOriginPolicyResource(resourceName string, resourceModel httpServletCrossOriginPolicyTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_http_servlet_cross_origin_policy" "%[1]s" {
  id                   = "%[2]s"
  cors_allowed_headers = %[3]s
}`, resourceName,
		resourceModel.id,
		acctest.StringSliceToTerraformString(resourceModel.corsAllowedHeaders))
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedHttpServletCrossOriginPolicyAttributes(config httpServletCrossOriginPolicyTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "http servlet cross origin policy"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.HttpServletCrossOriginPolicyApi.GetHttpServletCrossOriginPolicy(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "cors-allowed-headers",
			config.corsAllowedHeaders, response.CorsAllowedHeaders)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckHttpServletCrossOriginPolicyDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.HttpServletCrossOriginPolicyApi.GetHttpServletCrossOriginPolicy(ctx, testIdHttpServletCrossOriginPolicy).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Http Servlet Cross Origin Policy", testIdHttpServletCrossOriginPolicy)
	}
	return nil
}
