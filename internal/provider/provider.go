package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/accesscontrolhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/accesstokenvalidator"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/accountstatusnotificationhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/backend"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/connectioncriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/connectionhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/delegatedadminattribute"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/externalserver"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/gauge"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/httpservletextension"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/identitymapper"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logpublisher"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/plugin"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/recurringtask"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/requestcriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/restresourcetype"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/serverinstance"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/trustmanagerprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/virtualattribute"

	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// pingdirectoryProviderModel maps provider schema data to a Go type.
type pingdirectoryProviderModel struct {
	HttpsHost             types.String `tfsdk:"https_host"`
	Username              types.String `tfsdk:"username"`
	Password              types.String `tfsdk:"password"`
	InsecureTrustAllTls   types.Bool   `tfsdk:"insecure_trust_all_tls"`
	CACertificatePEMFiles types.Set    `tfsdk:"ca_certificate_pem_files"`
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
		Description: "PingDirectory Terraform Provider.",
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
			"insecure_trust_all_tls": schema.BoolAttribute{
				Description: "Set to true to trust any certificate when connecting to the PingDirectory server. This is insecure and should not be enabled outside of testing.",
				Optional:    true,
			},
			"ca_certificate_pem_files": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "Paths to files containing PEM-encoded certificates to be trusted as root CAs when connecting to the PingDirectory server over HTTPS. If not set, the host's root CA set will be used.",
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

	// Optional attributes
	var insecureTrustAllTls bool
	var err error
	if !config.InsecureTrustAllTls.IsUnknown() && !config.InsecureTrustAllTls.IsNull() {
		insecureTrustAllTls = config.InsecureTrustAllTls.ValueBool()
	} else {
		insecureTrustAllTls, err = strconv.ParseBool(os.Getenv("PINGDIRECTORY_PROVIDER_INSECURE_TRUST_ALL_TLS"))
		if err != nil {
			insecureTrustAllTls = false
			tflog.Info(ctx, "Failed to parse boolean from 'PINGDIRECTORY_PROVIDER_INSECURE_TRUST_ALL_TLS' environment variable, defaulting 'insecure_trust_all_tls' to false")
		}
	}

	var caCertPemFiles []string
	if !config.CACertificatePEMFiles.IsUnknown() && !config.CACertificatePEMFiles.IsNull() {
		config.CACertificatePEMFiles.ElementsAs(ctx, &caCertPemFiles, false)
	} else {
		pemFilesEnvVar := os.Getenv("PINGDIRECTORY_PROVIDER_CA_CERTIFICATE_PEM_FILES")
		if len(pemFilesEnvVar) == 0 {
			tflog.Info(ctx, "Did not find any certificate paths specified via the 'PINGDIRECTORY_PROVIDER_CA_CERTIFICATE_PEM_FILES' environment variable, using the host's root CA set")
		} else {
			caCertPemFiles = strings.Split(pemFilesEnvVar, ",")
		}
	}

	var caCertPool *x509.CertPool
	if len(caCertPemFiles) == 0 {
		tflog.Info(ctx, "No CA certs specified, using the host's root CA set")
		caCertPool = nil
	} else {
		caCertPool = x509.NewCertPool()
		for _, pemFilename := range caCertPemFiles {
			// Load CA cert
			caCert, err := os.ReadFile(pemFilename)
			if err != nil {
				resp.Diagnostics.AddError("Failed to read CA PEM certificate file: "+pemFilename, err.Error())
			}
			tflog.Info(ctx, "Adding CA cert from file: "+pemFilename)
			if !caCertPool.AppendCertsFromPEM(caCert) {
				resp.Diagnostics.AddWarning("Failed to parse certificate", "Failed to parse CA PEM certificate from file: "+pemFilename)
			}
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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureTrustAllTls,
			RootCAs:            caCertPool,
		},
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
		accountstatusnotificationhandler.NewAdminAlertAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewErrorLogAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewGroovyScriptedAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewMultiPartEmailAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewSmtpAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewThirdPartyAccountStatusNotificationHandlerResource,
		accesstokenvalidator.NewJwtAccessTokenValidatorResource,
		accesstokenvalidator.NewMockAccessTokenValidatorResource,
		accesstokenvalidator.NewPingFederateAccessTokenValidatorResource,
		accesstokenvalidator.NewThirdPartyAccessTokenValidatorResource,
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
		config.NewConsentDefinitionLocalizationResource,
		config.NewConsentServiceResource,
		config.NewDebugTargetResource,
		config.NewDelegatedAdminResourceRightsResource,
		config.NewDelegatedAdminRightsResource,
		config.NewGlobalConfigurationResource,
		config.NewHttpServletCrossOriginPolicyResource,
		config.NewLocationResource,
		config.NewRootDnResource,
		config.NewRootDnUserResource,
		config.NewTopologyAdminUserResource,
		connectioncriteria.NewAggregateConnectionCriteriaResource,
		connectioncriteria.NewSimpleConnectionCriteriaResource,
		connectioncriteria.NewThirdPartyConnectionCriteriaResource,
		connectionhandler.NewHttpConnectionHandlerResource,
		connectionhandler.NewJmxConnectionHandlerResource,
		connectionhandler.NewLdapConnectionHandlerResource,
		connectionhandler.NewLdifConnectionHandlerResource,
		delegatedadminattribute.NewCertificateDelegatedAdminAttributeResource,
		delegatedadminattribute.NewGenericDelegatedAdminAttributeResource,
		delegatedadminattribute.NewPhotoDelegatedAdminAttributeResource,
		externalserver.NewActiveDirectoryExternalServerResource,
		externalserver.NewAmazonAwsExternalServerResource,
		externalserver.NewConjurExternalServerResource,
		externalserver.NewHttpExternalServerResource,
		externalserver.NewJdbcExternalServerResource,
		externalserver.NewLdapExternalServerResource,
		externalserver.NewNokiaDsExternalServerResource,
		externalserver.NewNokiaProxyServerExternalServerResource,
		externalserver.NewOpendjExternalServerResource,
		externalserver.NewOracleUnifiedDirectoryExternalServerResource,
		externalserver.NewPingIdentityDsExternalServerResource,
		externalserver.NewPingIdentityProxyServerExternalServerResource,
		externalserver.NewPingOneHttpExternalServerResource,
		externalserver.NewSmtpExternalServerResource,
		externalserver.NewSyslogExternalServerResource,
		externalserver.NewVaultExternalServerResource,
		gauge.NewIndicatorGaugeResource,
		gauge.NewNumericGaugeResource,
		httpservletextension.NewAvailabilityStateHttpServletExtensionResource,
		httpservletextension.NewConfigHttpServletExtensionResource,
		httpservletextension.NewConsentHttpServletExtensionResource,
		httpservletextension.NewDelegatedAdminHttpServletExtensionResource,
		httpservletextension.NewDirectoryRestApiHttpServletExtensionResource,
		httpservletextension.NewFileServerHttpServletExtensionResource,
		httpservletextension.NewGroovyScriptedHttpServletExtensionResource,
		httpservletextension.NewLdapMappedScimHttpServletExtensionResource,
		httpservletextension.NewQuickstartHttpServletExtensionResource,
		httpservletextension.NewScim2HttpServletExtensionResource,
		httpservletextension.NewThirdPartyHttpServletExtensionResource,
		httpservletextension.NewVelocityHttpServletExtensionResource,
		identitymapper.NewAggregateIdentityMapperResource,
		identitymapper.NewExactMatchIdentityMapperResource,
		identitymapper.NewGroovyScriptedIdentityMapperResource,
		identitymapper.NewRegularExpressionIdentityMapperResource,
		identitymapper.NewThirdPartyIdentityMapperResource,
		logpublisher.NewAdminAlertAccessLogPublisherResource,
		logpublisher.NewCommonLogFileHttpOperationLogPublisherResource,
		logpublisher.NewConsoleJsonAccessLogPublisherResource,
		logpublisher.NewConsoleJsonAuditLogPublisherResource,
		logpublisher.NewConsoleJsonErrorLogPublisherResource,
		logpublisher.NewConsoleJsonHttpOperationLogPublisherResource,
		logpublisher.NewDebugAccessLogPublisherResource,
		logpublisher.NewDetailedHttpOperationLogPublisherResource,
		logpublisher.NewFileBasedAccessLogPublisherResource,
		logpublisher.NewFileBasedAuditLogPublisherResource,
		logpublisher.NewFileBasedDebugLogPublisherResource,
		logpublisher.NewFileBasedErrorLogPublisherResource,
		logpublisher.NewFileBasedJsonAuditLogPublisherResource,
		logpublisher.NewFileBasedJsonHttpOperationLogPublisherResource,
		logpublisher.NewFileBasedTraceLogPublisherResource,
		logpublisher.NewGroovyScriptedAccessLogPublisherResource,
		logpublisher.NewGroovyScriptedErrorLogPublisherResource,
		logpublisher.NewGroovyScriptedFileBasedAccessLogPublisherResource,
		logpublisher.NewGroovyScriptedFileBasedErrorLogPublisherResource,
		logpublisher.NewGroovyScriptedHttpOperationLogPublisherResource,
		logpublisher.NewJdbcBasedAccessLogPublisherResource,
		logpublisher.NewJdbcBasedErrorLogPublisherResource,
		logpublisher.NewJsonAccessLogPublisherResource,
		logpublisher.NewJsonErrorLogPublisherResource,
		logpublisher.NewOperationTimingAccessLogPublisherResource,
		logpublisher.NewSyslogBasedAccessLogPublisherResource,
		logpublisher.NewSyslogBasedErrorLogPublisherResource,
		logpublisher.NewSyslogJsonAccessLogPublisherResource,
		logpublisher.NewSyslogJsonAuditLogPublisherResource,
		logpublisher.NewSyslogJsonErrorLogPublisherResource,
		logpublisher.NewSyslogJsonHttpOperationLogPublisherResource,
		logpublisher.NewSyslogTextAccessLogPublisherResource,
		logpublisher.NewSyslogTextErrorLogPublisherResource,
		logpublisher.NewThirdPartyAccessLogPublisherResource,
		logpublisher.NewThirdPartyErrorLogPublisherResource,
		logpublisher.NewThirdPartyFileBasedAccessLogPublisherResource,
		logpublisher.NewThirdPartyFileBasedErrorLogPublisherResource,
		logpublisher.NewThirdPartyHttpOperationLogPublisherResource,
		plugin.NewAttributeMapperPluginResource,
		plugin.NewChangeSubscriptionNotificationPluginResource,
		plugin.NewChangelogPasswordEncryptionPluginResource,
		plugin.NewCleanUpExpiredPingfederatePersistentAccessGrantsPluginResource,
		plugin.NewCleanUpExpiredPingfederatePersistentSessionsPluginResource,
		plugin.NewCleanUpInactivePingfederatePersistentSessionsPluginResource,
		plugin.NewComposedAttributePluginResource,
		plugin.NewDelayPluginResource,
		plugin.NewDnMapperPluginResource,
		plugin.NewEncryptAttributeValuesPluginResource,
		plugin.NewGroovyScriptedPluginResource,
		plugin.NewInternalSearchRatePluginResource,
		plugin.NewLastAccessTimePluginResource,
		plugin.NewLastModPluginResource,
		plugin.NewLdapResultCodeTrackerPluginResource,
		plugin.NewModifiablePasswordPolicyStatePluginResource,
		plugin.NewMonitorHistoryPluginResource,
		plugin.NewPassThroughAuthenticationPluginResource,
		plugin.NewPasswordPolicyImportPluginResource,
		plugin.NewPeriodicGcPluginResource,
		plugin.NewPeriodicStatsLoggerPluginResource,
		plugin.NewPingOnePassThroughAuthenticationPluginResource,
		plugin.NewPluggablePassThroughAuthenticationPluginResource,
		plugin.NewProcessingTimeHistogramPluginResource,
		plugin.NewProfilerPluginResource,
		plugin.NewPurgeExpiredDataPluginResource,
		plugin.NewReferentialIntegrityPluginResource,
		plugin.NewReferralOnUpdatePluginResource,
		plugin.NewSearchShutdownPluginResource,
		plugin.NewSevenBitCleanPluginResource,
		plugin.NewSimpleToExternalBindPluginResource,
		plugin.NewSnmpSubagentPluginResource,
		plugin.NewStatsCollectorPluginResource,
		plugin.NewSubOperationTimingPluginResource,
		plugin.NewThirdPartyPluginResource,
		plugin.NewUniqueAttributePluginResource,
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
		virtualattribute.NewConstructedVirtualAttributeResource,
		virtualattribute.NewCurrentTimeVirtualAttributeResource,
		virtualattribute.NewDnJoinVirtualAttributeResource,
		virtualattribute.NewEntryChecksumVirtualAttributeResource,
		virtualattribute.NewEntryDnVirtualAttributeResource,
		virtualattribute.NewEqualityJoinVirtualAttributeResource,
		virtualattribute.NewGroovyScriptedVirtualAttributeResource,
		virtualattribute.NewHasSubordinatesVirtualAttributeResource,
		virtualattribute.NewIdentifyReferencesVirtualAttributeResource,
		virtualattribute.NewInstanceNameVirtualAttributeResource,
		virtualattribute.NewIsMemberOfVirtualAttributeResource,
		virtualattribute.NewMemberOfServerGroupVirtualAttributeResource,
		virtualattribute.NewMemberVirtualAttributeResource,
		virtualattribute.NewMirrorVirtualAttributeResource,
		virtualattribute.NewNumSubordinatesVirtualAttributeResource,
		virtualattribute.NewPasswordPolicyStateJsonVirtualAttributeResource,
		virtualattribute.NewReplicationStateDetailVirtualAttributeResource,
		virtualattribute.NewReverseDnJoinVirtualAttributeResource,
		virtualattribute.NewShortUniqueIdVirtualAttributeResource,
		virtualattribute.NewSubschemaSubentryVirtualAttributeResource,
		virtualattribute.NewThirdPartyVirtualAttributeResource,
		virtualattribute.NewUserDefinedVirtualAttributeResource,
	}
}
