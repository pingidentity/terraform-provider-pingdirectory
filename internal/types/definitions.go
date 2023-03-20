package types

import (
	client9200 "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
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
	ApiClientV9200 *client9200.APIClient
}
