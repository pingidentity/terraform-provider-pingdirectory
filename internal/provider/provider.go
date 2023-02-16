package provider

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/accesscontrolhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/backend"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/connectionhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/gauge"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/recurringtask"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/requestcriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/restresourcetype"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/serverinstance"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/trustmanagerprovider"

	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// pingdirectoryProviderModel maps provider schema data to a Go type.
type pingdirectoryProviderModel struct {
	HttpsHost types.String `tfsdk:"https_host"`
	Username  types.String `tfsdk:"username"`
	Password  types.String `tfsdk:"password"`
}

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &pingdirectoryProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &pingdirectoryProvider{}
}

// pingdirectoryProvider is the provider implementation.
type pingdirectoryProvider struct{}

// Metadata returns the provider type name.
func (p *pingdirectoryProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pingdirectory"
}

// GetSchema defines the provider-level schema for configuration data.
func (p *pingdirectoryProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "PingDirectory POC Provider.",
		Attributes: map[string]schema.Attribute{
			"https_host": schema.StringAttribute{
				Description: "URI for PingDirectory HTTPS port.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for PingDirectory admin user.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for PingDirectory admin user.",
				Sensitive:   true,
				Optional:    true,
			},
		},
	}
}

// Configure prepares a PingDirectory LDAP client
func (p *pingdirectoryProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring PingDirectory client")

	// Retrieve provider data from configuration
	var config pingdirectoryProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide a https host to the provider
	var httpsHost string
	if config.HttpsHost.IsUnknown() {
		// Cannot connect to PingDirectory with an unknown value
		resp.Diagnostics.AddError(
			"Unable to connect to the PingDirectory instance",
			"Cannot use unknown value as https_host",
		)
	} else {
		if config.HttpsHost.IsNull() {
			httpsHost = os.Getenv("PINGDIRECTORY_PROVIDER_HTTPS_HOST")
		} else {
			httpsHost = config.HttpsHost.ValueString()
		}
		if httpsHost == "" {
			resp.Diagnostics.AddError(
				"Unable to find https_host",
				"https_host cannot be an empty string. Either set it in the configuration or use the PINGDIRECTORY_PROVIDER_HTTPS_HOST environment variable.",
			)
		}
	}

	// User must provide a username to the provider
	var username string
	if config.Username.IsUnknown() {
		// Cannot connect to PingDirectory with an unknown value
		resp.Diagnostics.AddError(
			"Unable to connect to the PingDirectory instance",
			"Cannot use unknown value as username",
		)
	} else {
		if config.Username.IsNull() {
			username = os.Getenv("PINGDIRECTORY_PROVIDER_USERNAME")
		} else {
			username = config.Username.ValueString()
		}
		if username == "" {
			resp.Diagnostics.AddError(
				"Unable to find username",
				"username cannot be an empty string. Either set it in the configuration or use the PINGDIRECTORY_PROVIDER_USERNAME environment variable.",
			)
		}
	}

	// User must provide a username to the provider
	var password string
	if config.Password.IsUnknown() {
		// Cannot connect to PingDirectory with an unknown value
		resp.Diagnostics.AddError(
			"Unable to connect to the PingDirectory instance",
			"Cannot use unknown value as password",
		)
	} else {
		if config.Password.IsNull() {
			password = os.Getenv("PINGDIRECTORY_PROVIDER_PASSWORD")
		} else {
			password = config.Password.ValueString()
		}
		if password == "" {
			resp.Diagnostics.AddError(
				"Unable to find password",
				"password cannot be an empty string. Either set it in the configuration or use the PINGDIRECTORY_PROVIDER_PASSWORD environment variable.",
			)
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Make the PingDirectory config and API client info available during DataSource and Resource
	// type Configure methods.
	var resourceConfig internaltypes.ResourceConfiguration
	providerConfig := internaltypes.ProviderConfiguration{
		HttpsHost: httpsHost,
		Username:  username,
		Password:  password,
	}
	resourceConfig.ProviderConfig = providerConfig
	clientConfig := client.NewConfiguration()
	clientConfig.Servers = client.ServerConfigurations{
		{
			URL: httpsHost + "/config",
		},
	}
	//TODO THIS IS NOT SAFE!! Eventually need to add way to trust a specific cert/signer here rather than just trusting everything
	//https://stackoverflow.com/questions/12122159/how-to-do-a-https-request-with-bad-certificate
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	clientConfig.HTTPClient = httpClient
	resourceConfig.ApiClient = client.NewAPIClient(clientConfig)
	resp.ResourceData = resourceConfig

	tflog.Info(ctx, "Configured PingDirectory client", map[string]interface{}{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *pingdirectoryProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
// Maintain alphabetical order for ease of management
func (p *pingdirectoryProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		accesscontrolhandler.NewDseeCompatAccessControlHandlerResource,
		backend.NewAlarmBackendResource,
		backend.NewAlertBackendResource,
		backend.NewBackupBackendResource,
		backend.NewChangelogBackendResource,
		backend.NewConfigFileHandlerBackendResource,
		backend.NewEncryptionSettingsBackendResource,
		backend.NewLocalDbBackendResource,
		backend.NewMetricsBackendResource,
		backend.NewMonitorBackendResource,
		backend.NewSchemaBackendResource,
		backend.NewTaskBackendResource,
		backend.NewTrustStoreBackendResource,
		config.NewConsentDefinitionResource,
		config.NewGlobalConfigurationResource,
		config.NewLocationResource,
		config.NewRootDnResource,
		config.NewTopologyAdminUserResource,
		connectionhandler.NewHttpConnectionHandlerResource,
		connectionhandler.NewJmxConnectionHandlerResource,
		connectionhandler.NewLdapConnectionHandlerResource,
		connectionhandler.NewLdifConnectionHandlerResource,
		gauge.NewIndicatorGaugeResource,
		gauge.NewNumericGaugeResource,
		recurringtask.NewBackupRecurringTaskResource,
		recurringtask.NewCollectSupportDataRecurringTaskResource,
		recurringtask.NewDelayRecurringTaskResource,
		recurringtask.NewEnterLockdownModeRecurringTaskResource,
		recurringtask.NewExecRecurringTaskResource,
		recurringtask.NewFileRetentionRecurringTaskResource,
		recurringtask.NewGenerateServerProfileRecurringTaskResource,
		recurringtask.NewLdifExportRecurringTaskResource,
		recurringtask.NewLeaveLockdownModeRecurringTaskResource,
		recurringtask.NewStaticallyDefinedRecurringTaskResource,
		recurringtask.NewThirdPartyRecurringTaskResource,
		requestcriteria.NewAggregateRequestCriteriaResource,
		requestcriteria.NewRootDseRequestCriteriaResource,
		requestcriteria.NewSimpleRequestCriteriaResource,
		requestcriteria.NewThirdPartyRequestCriteriaResource,
		restresourcetype.NewGenericRestResourceTypeResource,
		restresourcetype.NewGroupRestResourceTypeResource,
		restresourcetype.NewUserRestResourceTypeResource,
		serverinstance.NewAuthorizeServerInstanceResource,
		serverinstance.NewDirectoryServerInstanceResource,
		serverinstance.NewProxyServerInstanceResource,
		serverinstance.NewSyncServerInstanceResource,
		trustmanagerprovider.NewBlindTrustManagerProviderResource,
		trustmanagerprovider.NewFileBasedTrustManagerProviderResource,
		trustmanagerprovider.NewJvmDefaultTrustManagerProviderResource,
		trustmanagerprovider.NewThirdPartyTrustManagerProviderResource,
	}
}
