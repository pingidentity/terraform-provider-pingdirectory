package trustedcertificate_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdTrustedCertificate = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type trustedCertificateTestModel struct {
	id          string
	certificate string
}

func TestAccTrustedCertificate(t *testing.T) {
	resourceName := "myresource"

	tempDirString := t.TempDir()
	d1 := []byte(`-----BEGIN CERTIFICATE-----
MIIFkDCCA3gCCQDqVJKvXI7duzANBgkqhkiG9w0BAQsFADCBiTELMAkGA1UEBhMC
dXMxETAPBgNVBAgMCGNvbG9yYWRvMQ8wDQYDVQQHDAZkZW52ZXIxDTALBgNVBAoM
BHBpbmcxDzANBgNVBAsMBmRldm9wczESMBAGA1UEAwwJZGlyZWN0b3J5MSIwIAYJ
KoZIhvcNAQkBFhNleGFtcGxlQGV4YW1wbGUuY29tMB4XDTIzMDUyNjE5MjAyMloX
DTI0MDUyNTE5MjAyMlowgYkxCzAJBgNVBAYTAnVzMREwDwYDVQQIDAhjb2xvcmFk
bzEPMA0GA1UEBwwGZGVudmVyMQ0wCwYDVQQKDARwaW5nMQ8wDQYDVQQLDAZkZXZv
cHMxEjAQBgNVBAMMCWRpcmVjdG9yeTEiMCAGCSqGSIb3DQEJARYTZXhhbXBsZUBl
eGFtcGxlLmNvbTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAOF2vk92
tFeZLwVBhOvfmaNt+IMdve91U9ZUUXgXtiIA3hEb/sckOXx37KKAobZY00uOjiP1
hh4fotg0OUW1C3qQDk4wh6XWJ22dhrmvn33hl+q2L8mSNlR8ICJfSr7YlfKWJRy5
kWeKOchDIf4c0H6G5ZFcIH000JGGpZX6ut6Fndv0PDRPHXnl5H1UGaTDU9nRMCtQ
3qRTZleFEWkRkx52C6EYPAW0rxLKyPGqejkbay7EuAqaBB0/64RbY7rnOWT3Y+oW
IQafXMwwYTMFhgt1Reiztb3iCDVeTiQJ1QbIJeSptTCndIoxRuOIJIWw7UePSOI7
uV47/4JihP60J7xtjPXtBiIwKejj6ZBobF9mLtZaIeYiVfUfppnjh0lyyVgeKWFm
JHRj46CPFIHgRU7yYwqj9bkW2N51iot/VBWFBVXxRGhkBkIYDK1X9qbSX6eFg/4p
l8326hJBtHyBILvHejA/t7A9njfg13HghT9QhYZX8U585duvRkenc793v7qj+9TR
R6FOj5jsgTf8ls1P3SuJvOxInV0w+So+Ee/vkktB0jlhu2ytYFm8mqplwEtX1Kaj
tsHeVFjMDq+ZX1t/T8tsvlSZK5sFi8yMWMeJ/Ce/qbC0kuoQvoa1DE7pbTJXlyTk
dJmMKIDa1EeTgdwAFXSsu/Jl9vE3zNEmVBG1AgMBAAEwDQYJKoZIhvcNAQELBQAD
ggIBAKuniXjjMllrgL07O+uZ8a6f7lT3bIwU9vJs5QmvtiD83OULrrCL1v7uEkrB
OWH1QNftuCJkWSbhAx8IaK9Ws7G7wtuRr/WtxrtQOQGz8daU0FNH9b49AZ0K+K4u
7ceyFTyiIoyb8xhzMXCLpY0VKhkB3JJnxw4On8E1kqT7Jo5TKvcq4QED3Y/pLHmF
tPu7N7OmzFT4btGCBdXV1jZoBoTygkd0HJsARE929B/1hCQVS4snGJDy5cYi6IwE
AjcvpV/R2hv3VXvrOjEL9VW1M5KGPBPHZo5pTMvVuP60IYyhL+coB1o78ODD+0B+
A7C/WBHKYFlc5fnKfcEzm+DG5RcUbeCxcde/yP2S2l7h3oue4Fx4zzKWlaLcw9FT
W6DF/856GCBzBztOH4gaHm8rfy6RtbM4akLMMwoU/i0hkRTT676S5L/zWRLNdvFE
ntOKJEKfSXkG23Ea28LH0XsqQ1eOfnBI1nXEE7iSrbfvP8LiMmKXShquDIgo3ly/
XSFRJt2K9WS3B+CI3vxkM1+J1C1m6Q4CUoy+VEo0yJCmdgckE+Ijwb+AYglnMLmk
O153rbh1O3sXFjeKFSvpi6BM4OBaTDwtlZL+ZtDjvLX5xY278udB140n+XYdJaW7
5ZZbQAG1UD52rb54Z5Si0Z04t5lgP2Qo2o7Wak/Y1yxeUvxs
-----END CERTIFICATE-----`)
	err := os.WriteFile(tempDirString+"cert.pem", d1, 0600)
	if err != nil {
		t.Error(err)
	}
	initialResourceModel := trustedCertificateTestModel{
		id:          testIdTrustedCertificate,
		certificate: tempDirString + "cert.pem",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckTrustedCertificateDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccTrustedCertificateResource(resourceName, initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccTrustedCertificateResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_trusted_certificate." + resourceName,
				ImportStateId:     initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccTrustedCertificateResource(resourceName string, resourceModel trustedCertificateTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_trusted_certificate" "%[1]s" {
  id          = "%[2]s"
  certificate = file("%[3]s")
}

data "pingdirectory_trusted_certificate" "%[1]s" {
  id = "%[2]s"
  depends_on = [
    pingdirectory_trusted_certificate.%[1]s
  ]
}`, resourceName,
		resourceModel.id,
		resourceModel.certificate)
}

// Test that any objects created by the test are destroyed
func testAccCheckTrustedCertificateDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.TrustedCertificateApi.GetTrustedCertificate(ctx, testIdTrustedCertificate).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Trusted Certificate", testIdTrustedCertificate)
	}
	return nil
}
