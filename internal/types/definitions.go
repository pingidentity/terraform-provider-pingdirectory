package types

import client "github.com/pingidentity/pingdirectory-go-client/v9100"

// Configuration used by the provider and resources
type ProviderConfiguration struct {
	HttpsHost string
	Username  string
	Password  string
}

// Configuration passed to resources
type ResourceConfiguration struct {
	ProviderConfig ProviderConfiguration
	ApiClient      *client.APIClient
}
