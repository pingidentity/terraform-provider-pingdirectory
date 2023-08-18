package consentservice_test

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
type consentServiceTestModel struct {
	enabled                    bool
	base_dn                    string
	bind_dn                    string
	unprivileged_consent_scope string
	privileged_consent_scope   string
	search_size_limit          int64
}

func TestAccConsentService(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := consentServiceTestModel{
		enabled:                    true,
		base_dn:                    "ou=consents,dc=example,dc=com",
		bind_dn:                    "cn=consent service account",
		unprivileged_consent_scope: "urn:pingdirectory:consent",
		privileged_consent_scope:   "urn:pingdirectory:consent_admin",
		search_size_limit:          90,
	}
	updatedResourceModel := consentServiceTestModel{
		enabled:                    true,
		base_dn:                    "ou=consents1,dc=example,dc=com",
		bind_dn:                    "cn=consent1 service account",
		unprivileged_consent_scope: "urn:pingdirectory:consent",
		privileged_consent_scope:   "urn:pingdirectory:consent_admin",
		search_size_limit:          50,
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
				Config: testAccConsentServiceResource(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedConsentServiceAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_consent_service.%s", resourceName), "enabled", strconv.FormatBool(initialResourceModel.enabled)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_consent_service.%s", resourceName), "base_dn", initialResourceModel.base_dn),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_consent_service.%s", resourceName), "bind_dn", initialResourceModel.bind_dn),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_consent_service.%s", resourceName), "unprivileged_consent_scope", initialResourceModel.unprivileged_consent_scope),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_consent_service.%s", resourceName), "privileged_consent_scope", initialResourceModel.privileged_consent_scope),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.pingdirectory_consent_service.%s", resourceName), "search_size_limit", strconv.FormatInt(initialResourceModel.search_size_limit, 10)),
				),
			},
			{
				// Test updating some fields
				Config: testAccConsentServiceResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedConsentServiceAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccConsentServiceResource(resourceName, updatedResourceModel),
				ResourceName:            "pingdirectory_default_consent_service." + resourceName,
				ImportStateId:           resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func testAccConsentServiceResource(resourceName string, resourceModel consentServiceTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_consent_service" "%[1]s" {
  enabled                    = %[2]t
  base_dn                    = "%[3]s"
  bind_dn                    = "%[4]s"
  unprivileged_consent_scope = "%[5]s"
  privileged_consent_scope   = "%[6]s"
  search_size_limit          = %[7]d
}

data "pingdirectory_consent_service" "%[1]s" {
  depends_on = [
    pingdirectory_default_consent_service.%[1]s
  ]
}`, resourceName,
		resourceModel.enabled,
		resourceModel.base_dn,
		resourceModel.bind_dn,
		resourceModel.unprivileged_consent_scope,
		resourceModel.privileged_consent_scope,
		resourceModel.search_size_limit)
}

/*
	enabled:                    true,
	base_dn:                    []string{"ou=consents,dc=example,dc=com"},
	bind_dn:                    []string{"cn=consent service account"},
	unprivileged_consent_scope: "urn:pingdirectory:consent",
	privileged_consent_scope:   "urn:pingdirectory:consent_admin",
	search_size_limit:          90,
*/
// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedConsentServiceAttributes(config consentServiceTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ConsentServiceApi.GetConsentService(ctx).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Consent Service"
		err = acctest.TestAttributesMatchBool(resourceType, nil, "enabled",
			config.enabled, response.Enabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "base-dn",
			config.base_dn, response.BaseDN)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "bind-dn",
			config.bind_dn, response.BindDN)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "unprivileged-consent-scope",
			config.unprivileged_consent_scope, response.UnprivilegedConsentScope)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "privileged-consent-scope",
			config.privileged_consent_scope, response.PrivilegedConsentScope)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, nil, "search-size-limit",
			config.search_size_limit, *response.SearchSizeLimit)
		if err != nil {
			return err
		}
		return nil
	}
}
