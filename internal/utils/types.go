package utils

import client "github.com/pingidentity/pingdata-config-api-go-client"

// Configuration used by the provider and resources
type ProviderConfiguration struct {
	LdapHost  string
	HttpsHost string
	Username  string
	Password  string
	//DefaultUserPassword string
}

// Configuration passed to resources
type ResourceConfiguration struct {
	ProviderConfig ProviderConfiguration
	ApiClient      *client.APIClient
}
