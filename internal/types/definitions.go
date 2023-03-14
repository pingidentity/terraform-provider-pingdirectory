package types

import (
	client9100 "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
	client9200 "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
)

// Configuration used by the provider and resources
type ProviderConfiguration struct {
	HttpsHost            string
	Username             string
	Password             string
	PingDirectoryVersion string
}

// Configuration passed to resources
type ResourceConfiguration struct {
	ProviderConfig ProviderConfiguration
	ApiClientV9100 *client9100.APIClient
	ApiClientV9200 *client9200.APIClient
}

// Supported PingDirectory versions
const (
	PingDirectory9100 = "9.1.0.0"
	PingDirectory9200 = "9.2.0.0"
)
