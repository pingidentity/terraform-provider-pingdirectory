// Copyright © 2025 Ping Identity Corporation

package types

import (
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
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
	ApiClient      *client.APIClient
}
