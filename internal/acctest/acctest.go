// Copyright © 2025 Ping Identity Corporation

package acctest

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Verify that any required environment variables are set before the test begins
func ConfigurationPreCheck(t *testing.T) {
	envVars := []string{
		"PINGDIRECTORY_PROVIDER_HTTPS_HOST",
		"PINGDIRECTORY_PROVIDER_USERNAME",
		"PINGDIRECTORY_PROVIDER_PASSWORD",
		"PINGDIRECTORY_PROVIDER_INSECURE_TRUST_ALL_TLS",
		"PINGDIRECTORY_PROVIDER_PRODUCT_VERSION",
	}

	errorFound := false
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			t.Errorf("The '%s' environment variable must be set to run acceptance tests", envVar)
			errorFound = true
		}
	}

	if errorFound {
		t.FailNow()
	}
}

func TestClient() *client.APIClient {
	httpsHost := os.Getenv("PINGDIRECTORY_PROVIDER_HTTPS_HOST")
	clientConfig := client.NewConfiguration()
	clientConfig.Servers = client.ServerConfigurations{
		{
			URL: httpsHost + "/config",
		},
	}
	// Trusting all for the acceptance tests, since they run on localhost
	// May want to incorporate actual trust here in the future.
	//#nosec G402
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	clientConfig.HTTPClient = httpClient
	return client.NewAPIClient(clientConfig)
}

func TestBasicAuthContext() context.Context {
	ctx := context.Background()
	return config.BasicAuthContext(ctx, os.Getenv("PINGDIRECTORY_PROVIDER_USERNAME"), os.Getenv("PINGDIRECTORY_PROVIDER_PASSWORD"))
}

// Convert a string slice to the format used in Terraform files
func StringSliceToTerraformString(values []string) string {
	var builder strings.Builder
	builder.WriteString("[")
	for i, str := range values {
		builder.WriteString(fmt.Sprintf("\"%s\"", str))
		if i < len(values)-1 {
			builder.WriteString(",")
		}
	}
	builder.WriteString("]")
	return builder.String()
}

// Utility methods for testing whether attributes match the expected values

// Test if string attributes match
func TestAttributesMatchString(resourceType string, resourceName *string, attributeName, expected, found string) error {
	if expected != found {
		return mismatchedAttributeError(resourceType, resourceName, attributeName, expected, found)
	}
	return nil
}

// Test if expected string matches found string pointer
func TestAttributesMatchStringPointer(resourceType string, resourceName *string, attributeName, expected string, found *string) error {
	if found == nil && expected != "" {
		// Expect empty string to match nil pointer
		return missingAttributeError(resourceType, resourceName, attributeName, expected)
	}
	if found != nil {
		return TestAttributesMatchString(resourceType, resourceName, attributeName, expected, *found)
	}
	return nil
}

// Test if boolean attributes match
func TestAttributesMatchBool(resourceType string, resourceName *string, attributeName string, expected, found bool) error {
	if expected != found {
		return mismatchedAttributeError(resourceType, resourceName, attributeName, strconv.FormatBool(expected), strconv.FormatBool(found))
	}
	return nil
}

// Test if int attributes match
func TestAttributesMatchInt(resourceType string, resourceName *string, attributeName string, expected, found int64) error {
	if expected != found {
		return mismatchedAttributeError(resourceType, resourceName, attributeName, strconv.FormatInt(expected, 10), strconv.FormatInt(found, 10))
	}
	return nil
}

// Test if string slice attributes match
func TestAttributesMatchStringSlice(resourceType string, resourceName *string, attributeName string, expected, found []string) error {
	if !types.SetsEqual(expected, found) {
		return mismatchedAttributeError(resourceType, resourceName, attributeName, StringSliceToTerraformString(expected), StringSliceToTerraformString(found))
	}
	return nil
}

func ExpectedDestroyError(resourceType, resourceName string) error {
	return fmt.Errorf("%s '%s' still exists after tests. Expected it to be destroyed.", resourceType, resourceName)
}

func mismatchedAttributeError(resourceType string, resourceName *string, attributeName, expected, found string) error {
	if resourceName == nil {
		return mismatchedAttributeErrorSingletonResource(resourceType, attributeName, expected, found)
	}
	return fmt.Errorf("mismatched %s attribute for %s '%s'. expected '%s', found '%s'", attributeName, resourceType, *resourceName, expected, found)
}

func mismatchedAttributeErrorSingletonResource(resourceType, attributeName, expected, found string) error {
	return fmt.Errorf("mismatched %s attribute for %s. expected '%s', found '%s'", attributeName, resourceType, expected, found)
}

func missingAttributeError(resourceType string, resourceName *string, attributeName, expected string) error {
	if resourceName == nil {
		return missingAttributeErrorSingletonResource(resourceType, attributeName, expected)
	}
	return fmt.Errorf("missing %s attribute for %s '%s'. expected '%s'", attributeName, resourceType, *resourceName, expected)
}

func missingAttributeErrorSingletonResource(resourceType, attributeName, expected string) error {
	return fmt.Errorf("missing %s attribute for %s. expected '%s'", attributeName, resourceType, expected)
}
