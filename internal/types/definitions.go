package types

import (
	client9300 "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
)

// Configuration used by the provider and resources
type ProviderConfiguration struct {
	HttpsHost      string
	Username       string
	Password       string
	ProductVersion string
}

// Configuration passed to resources
type ResourceConfiguration struct {
	ProviderConfig ProviderConfiguration
	ApiClientV9300 *client9300.APIClient
}
