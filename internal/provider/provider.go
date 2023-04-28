package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client9200 "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/accesscontrolhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/accesstokenvalidator"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/accountstatusnotificationhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/backend"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/certificatemapper"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/connectioncriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/connectionhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/delegatedadminattribute"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/entrycache"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/externalserver"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/gauge"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/gaugedatasource"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/httpservletextension"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/identitymapper"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/idtokenvalidator"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logpublisher"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/monitoringendpoint"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/plugin"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/recurringtask"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/requestcriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/restresourcetype"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/saslmechanismhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/scimresourcetype"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/serverinstance"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/trustmanagerprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/virtualattribute"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/webapplicationextension"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// pingdirectoryProviderModel maps provider schema data to a Go type.
type pingdirectoryProviderModel struct {
	HttpsHost             types.String `tfsdk:"https_host"`
	Username              types.String `tfsdk:"username"`
	Password              types.String `tfsdk:"password"`
	InsecureTrustAllTls   types.Bool   `tfsdk:"insecure_trust_all_tls"`
	CACertificatePEMFiles types.Set    `tfsdk:"ca_certificate_pem_files"`
	ProductVersion        types.String `tfsdk:"product_version"`
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
				Description: "URI for PingDirectory HTTPS port. Default value can be set with the `PINGDIRECTORY_PROVIDER_HTTPS_HOST` environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for PingDirectory admin user. Default value can be set with the `PINGDIRECTORY_PROVIDER_USERNAME` environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for PingDirectory admin user. Default value can be set with the `PINGDIRECTORY_PROVIDER_PASSWORD` environment variable.",
				Sensitive:   true,
				Optional:    true,
			},
			"insecure_trust_all_tls": schema.BoolAttribute{
				Description: "Set to true to trust any certificate when connecting to the PingDirectory server. This is insecure and should not be enabled outside of testing. Default value can be set with the `PINGDIRECTORY_PROVIDER_INSECURE_TRUST_ALL_TLS` environment variable.",
				Optional:    true,
			},
			"ca_certificate_pem_files": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "Paths to files containing PEM-encoded certificates to be trusted as root CAs when connecting to the PingDirectory server over HTTPS. If not set, the host's root CA set will be used. Default value can be set with the `PINGDIRECTORY_PROVIDER_CA_CERTIFICATE_PEM_FILES` environment variable, using commas to delimit multiple PEM files if necessary.",
				Optional:    true,
			},
			"product_version": schema.StringAttribute{
				Description: "Version of the PingDirectory server being configured. Default value can be set with the `PINGDIRECTORY_PROVIDER_PRODUCT_VERSION` environment variable.",
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

	var productVersion string
	var err error
	if !config.ProductVersion.IsUnknown() && !config.ProductVersion.IsNull() {
		productVersion = config.ProductVersion.ValueString()
	} else {
		productVersion = os.Getenv("PINGDIRECTORY_PROVIDER_PRODUCT_VERSION")
	}

	if productVersion == "" {
		resp.Diagnostics.AddError(
			"Unable to find PingDirectory version",
			"product_version cannot be an empty string. Either set it in the configuration or use the PINGDIRECTORY_PROVIDER_PRODUCT_VERSION environment variable.",
		)
	} else {
		// Validate the PingDirectory version
		productVersion, err = version.Parse(productVersion)
		if err != nil {
			resp.Diagnostics.AddError("Failed to parse PingDirectory version", err.Error())
		}
	}

	// Optional attributes
	var insecureTrustAllTls bool
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
		HttpsHost:      httpsHost,
		Username:       username,
		Password:       password,
		ProductVersion: productVersion,
	}
	resourceConfig.ProviderConfig = providerConfig
	//#nosec G402
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureTrustAllTls,
			RootCAs:            caCertPool,
		},
	}
	httpClient := &http.Client{Transport: tr}
	// Always create a client for the most recent version, since it is
	// the default used by resources that are compatible with multiple versions
	clientConfig9200 := client9200.NewConfiguration()
	clientConfig9200.Servers = client9200.ServerConfigurations{
		{
			URL: httpsHost + "/config",
		},
	}
	clientConfig9200.HTTPClient = httpClient
	resourceConfig.ApiClientV9200 = client9200.NewAPIClient(clientConfig9200)

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
		accountstatusnotificationhandler.NewDefaultAdminAlertAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewErrorLogAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewDefaultErrorLogAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewGroovyScriptedAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewDefaultGroovyScriptedAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewMultiPartEmailAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewDefaultMultiPartEmailAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewSmtpAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewDefaultSmtpAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewThirdPartyAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewDefaultThirdPartyAccountStatusNotificationHandlerResource,
		accesstokenvalidator.NewDefaultJwtAccessTokenValidatorResource,
		accesstokenvalidator.NewDefaultMockAccessTokenValidatorResource,
		accesstokenvalidator.NewDefaultPingFederateAccessTokenValidatorResource,
		accesstokenvalidator.NewDefaultThirdPartyAccessTokenValidatorResource,
		accesstokenvalidator.NewJwtAccessTokenValidatorResource,
		accesstokenvalidator.NewMockAccessTokenValidatorResource,
		accesstokenvalidator.NewPingFederateAccessTokenValidatorResource,
		accesstokenvalidator.NewThirdPartyAccessTokenValidatorResource,
		backend.NewAlarmBackendResource,
		backend.NewAlertBackendResource,
		backend.NewBackupBackendResource,
		backend.NewChangelogBackendResource,
		backend.NewConfigFileHandlerBackendResource,
		backend.NewCustomBackendResource,
		backend.NewDefaultLocalDbBackendResource,
		backend.NewEncryptionSettingsBackendResource,
		backend.NewLdifBackendResource,
		backend.NewLocalDbBackendResource,
		backend.NewMetricsBackendResource,
		backend.NewMonitorBackendResource,
		backend.NewSchemaBackendResource,
		backend.NewTaskBackendResource,
		backend.NewTrustStoreBackendResource,
		certificatemapper.NewDefaultFingerprintCertificateMapperResource,
		certificatemapper.NewDefaultGroovyScriptedCertificateMapperResource,
		certificatemapper.NewDefaultSubjectAttributeToUserAttributeCertificateMapperResource,
		certificatemapper.NewDefaultSubjectDnToUserAttributeCertificateMapperResource,
		certificatemapper.NewDefaultSubjectEqualsDnCertificateMapperResource,
		certificatemapper.NewDefaultThirdPartyCertificateMapperResource,
		certificatemapper.NewFingerprintCertificateMapperResource,
		certificatemapper.NewGroovyScriptedCertificateMapperResource,
		certificatemapper.NewSubjectAttributeToUserAttributeCertificateMapperResource,
		certificatemapper.NewSubjectDnToUserAttributeCertificateMapperResource,
		certificatemapper.NewSubjectEqualsDnCertificateMapperResource,
		certificatemapper.NewThirdPartyCertificateMapperResource,
		config.NewConsentDefinitionResource,
		config.NewDefaultConsentDefinitionResource,
		config.NewConsentDefinitionLocalizationResource,
		config.NewDefaultConsentDefinitionLocalizationResource,
		config.NewConsentServiceResource,
		config.NewDebugTargetResource,
		config.NewDefaultDebugTargetResource,
		config.NewDefaultLocationResource,
		config.NewDelegatedAdminResourceRightsResource,
		config.NewDefaultDelegatedAdminResourceRightsResource,
		config.NewDelegatedAdminRightsResource,
		config.NewDefaultDelegatedAdminRightsResource,
		config.NewDnMapResource,
		config.NewDefaultDnMapResource,
		config.NewGlobalConfigurationResource,
		config.NewHttpConfigurationResource,
		config.NewHttpServletCrossOriginPolicyResource,
		config.NewDefaultHttpServletCrossOriginPolicyResource,
		config.NewLocalDbIndexResource,
		config.NewDefaultLocalDbIndexResource,
		config.NewLocationResource,
		config.NewDefaultRecurringTaskChainResource,
		config.NewRecurringTaskChainResource,
		config.NewRootDnResource,
		config.NewRootDnUserResource,
		config.NewDefaultRootDnUserResource,
		config.NewScimAttributeResource,
		config.NewDefaultScimAttributeResource,
		config.NewScimAttributeMappingResource,
		config.NewDefaultScimAttributeMappingResource,
		config.NewScimSchemaResource,
		config.NewDefaultScimSchemaResource,
		config.NewTopologyAdminUserResource,
		config.NewDefaultTopologyAdminUserResource,
		connectioncriteria.NewAggregateConnectionCriteriaResource,
		connectioncriteria.NewDefaultAggregateConnectionCriteriaResource,
		connectionhandler.NewHttpConnectionHandlerResource,
		connectionhandler.NewDefaultHttpConnectionHandlerResource,
		connectionhandler.NewJmxConnectionHandlerResource,
		connectionhandler.NewDefaultJmxConnectionHandlerResource,
		connectionhandler.NewLdapConnectionHandlerResource,
		connectionhandler.NewDefaultLdapConnectionHandlerResource,
		connectionhandler.NewLdifConnectionHandlerResource,
		connectionhandler.NewDefaultLdifConnectionHandlerResource,
		connectioncriteria.NewSimpleConnectionCriteriaResource,
		connectioncriteria.NewDefaultSimpleConnectionCriteriaResource,
		connectioncriteria.NewThirdPartyConnectionCriteriaResource,
		connectioncriteria.NewDefaultThirdPartyConnectionCriteriaResource,
		delegatedadminattribute.NewCertificateDelegatedAdminAttributeResource,
		delegatedadminattribute.NewDefaultCertificateDelegatedAdminAttributeResource,
		delegatedadminattribute.NewDefaultGenericDelegatedAdminAttributeResource,
		delegatedadminattribute.NewDefaultPhotoDelegatedAdminAttributeResource,
		delegatedadminattribute.NewGenericDelegatedAdminAttributeResource,
		delegatedadminattribute.NewPhotoDelegatedAdminAttributeResource,
		entrycache.NewFifoEntryCacheResource,
		entrycache.NewDefaultFifoEntryCacheResource,
		externalserver.NewActiveDirectoryExternalServerResource,
		externalserver.NewAmazonAwsExternalServerResource,
		externalserver.NewConjurExternalServerResource,
		externalserver.NewDefaultActiveDirectoryExternalServerResource,
		externalserver.NewDefaultAmazonAwsExternalServerResource,
		externalserver.NewDefaultConjurExternalServerResource,
		externalserver.NewDefaultHttpExternalServerResource,
		externalserver.NewDefaultHttpProxyExternalServerResource,
		externalserver.NewDefaultJdbcExternalServerResource,
		externalserver.NewDefaultLdapExternalServerResource,
		externalserver.NewDefaultNokiaDsExternalServerResource,
		externalserver.NewDefaultNokiaProxyServerExternalServerResource,
		externalserver.NewDefaultOpendjExternalServerResource,
		externalserver.NewDefaultOracleUnifiedDirectoryExternalServerResource,
		externalserver.NewDefaultPingIdentityDsExternalServerResource,
		externalserver.NewDefaultPingIdentityProxyServerExternalServerResource,
		externalserver.NewDefaultPingOneHttpExternalServerResource,
		externalserver.NewDefaultSmtpExternalServerResource,
		externalserver.NewDefaultSyslogExternalServerResource,
		externalserver.NewDefaultVaultExternalServerResource,
		externalserver.NewHttpExternalServerResource,
		externalserver.NewHttpProxyExternalServerResource,
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
		gauge.NewDefaultIndicatorGaugeResource,
		gauge.NewDefaultNumericGaugeResource,
		gauge.NewIndicatorGaugeResource,
		gauge.NewNumericGaugeResource,
		gaugedatasource.NewIndicatorGaugeDataSourceResource,
		gaugedatasource.NewDefaultIndicatorGaugeDataSourceResource,
		gaugedatasource.NewNumericGaugeDataSourceResource,
		gaugedatasource.NewDefaultNumericGaugeDataSourceResource,
		httpservletextension.NewAvailabilityStateHttpServletExtensionResource,
		httpservletextension.NewConfigHttpServletExtensionResource,
		httpservletextension.NewConsentHttpServletExtensionResource,
		httpservletextension.NewDefaultAvailabilityStateHttpServletExtensionResource,
		httpservletextension.NewDefaultFileServerHttpServletExtensionResource,
		httpservletextension.NewDefaultGroovyScriptedHttpServletExtensionResource,
		httpservletextension.NewDefaultLdapMappedScimHttpServletExtensionResource,
		httpservletextension.NewDefaultPrometheusMonitoringHttpServletExtensionResource,
		httpservletextension.NewDefaultQuickstartHttpServletExtensionResource,
		httpservletextension.NewDefaultThirdPartyHttpServletExtensionResource,
		httpservletextension.NewDelegatedAdminHttpServletExtensionResource,
		httpservletextension.NewDirectoryRestApiHttpServletExtensionResource,
		httpservletextension.NewFileServerHttpServletExtensionResource,
		httpservletextension.NewGroovyScriptedHttpServletExtensionResource,
		httpservletextension.NewLdapMappedScimHttpServletExtensionResource,
		httpservletextension.NewPrometheusMonitoringHttpServletExtensionResource,
		httpservletextension.NewQuickstartHttpServletExtensionResource,
		httpservletextension.NewScim2HttpServletExtensionResource,
		httpservletextension.NewThirdPartyHttpServletExtensionResource,
		httpservletextension.NewVelocityHttpServletExtensionResource,
		identitymapper.NewAggregateIdentityMapperResource,
		identitymapper.NewDefaultAggregateIdentityMapperResource,
		identitymapper.NewDefaultExactMatchIdentityMapperResource,
		identitymapper.NewDefaultGroovyScriptedIdentityMapperResource,
		identitymapper.NewDefaultRegularExpressionIdentityMapperResource,
		identitymapper.NewDefaultThirdPartyIdentityMapperResource,
		identitymapper.NewExactMatchIdentityMapperResource,
		identitymapper.NewGroovyScriptedIdentityMapperResource,
		identitymapper.NewRegularExpressionIdentityMapperResource,
		identitymapper.NewThirdPartyIdentityMapperResource,
		idtokenvalidator.NewDefaultOpenidConnectIdTokenValidatorResource,
		idtokenvalidator.NewDefaultPingOneIdTokenValidatorResource,
		idtokenvalidator.NewOpenidConnectIdTokenValidatorResource,
		idtokenvalidator.NewPingOneIdTokenValidatorResource,
		logpublisher.NewAdminAlertAccessLogPublisherResource,
		logpublisher.NewCommonLogFileHttpOperationLogPublisherResource,
		logpublisher.NewConsoleJsonAccessLogPublisherResource,
		logpublisher.NewConsoleJsonAuditLogPublisherResource,
		logpublisher.NewConsoleJsonErrorLogPublisherResource,
		logpublisher.NewConsoleJsonHttpOperationLogPublisherResource,
		logpublisher.NewDebugAccessLogPublisherResource,
		logpublisher.NewDefaultAdminAlertAccessLogPublisherResource,
		logpublisher.NewDefaultCommonLogFileHttpOperationLogPublisherResource,
		logpublisher.NewDefaultConsoleJsonAuditLogPublisherResource,
		logpublisher.NewDefaultConsoleJsonHttpOperationLogPublisherResource,
		logpublisher.NewDefaultDebugAccessLogPublisherResource,
		logpublisher.NewDefaultDetailedHttpOperationLogPublisherResource,
		logpublisher.NewDefaultFileBasedAccessLogPublisherResource,
		logpublisher.NewDefaultFileBasedAuditLogPublisherResource,
		logpublisher.NewDefaultFileBasedDebugLogPublisherResource,
		logpublisher.NewDefaultFileBasedErrorLogPublisherResource,
		logpublisher.NewDefaultFileBasedJsonAuditLogPublisherResource,
		logpublisher.NewDefaultFileBasedJsonHttpOperationLogPublisherResource,
		logpublisher.NewDefaultFileBasedTraceLogPublisherResource,
		logpublisher.NewDefaultGroovyScriptedAccessLogPublisherResource,
		logpublisher.NewDefaultGroovyScriptedErrorLogPublisherResource,
		logpublisher.NewDefaultGroovyScriptedFileBasedAccessLogPublisherResource,
		logpublisher.NewDefaultGroovyScriptedFileBasedErrorLogPublisherResource,
		logpublisher.NewDefaultGroovyScriptedHttpOperationLogPublisherResource,
		logpublisher.NewDefaultJdbcBasedAccessLogPublisherResource,
		logpublisher.NewDefaultJdbcBasedErrorLogPublisherResource,
		logpublisher.NewDefaultJsonAccessLogPublisherResource,
		logpublisher.NewDefaultJsonErrorLogPublisherResource,
		logpublisher.NewDefaultOperationTimingAccessLogPublisherResource,
		logpublisher.NewDefaultSyslogBasedAccessLogPublisherResource,
		logpublisher.NewDefaultSyslogBasedErrorLogPublisherResource,
		logpublisher.NewDefaultSyslogJsonAccessLogPublisherResource,
		logpublisher.NewDefaultSyslogJsonAuditLogPublisherResource,
		logpublisher.NewDefaultSyslogJsonErrorLogPublisherResource,
		logpublisher.NewDefaultSyslogJsonHttpOperationLogPublisherResource,
		logpublisher.NewDefaultSyslogTextAccessLogPublisherResource,
		logpublisher.NewDefaultSyslogTextErrorLogPublisherResource,
		logpublisher.NewDefaultThirdPartyAccessLogPublisherResource,
		logpublisher.NewDefaultThirdPartyErrorLogPublisherResource,
		logpublisher.NewDefaultThirdPartyFileBasedAccessLogPublisherResource,
		logpublisher.NewDefaultThirdPartyFileBasedErrorLogPublisherResource,
		logpublisher.NewDefaultThirdPartyHttpOperationLogPublisherResource,
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
		monitoringendpoint.NewDefaultStatsdMonitoringEndpointResource,
		monitoringendpoint.NewStatsdMonitoringEndpointResource,
		plugin.NewAttributeMapperPluginResource,
		plugin.NewChangeSubscriptionNotificationPluginResource,
		plugin.NewChangelogPasswordEncryptionPluginResource,
		plugin.NewCleanUpExpiredPingfederatePersistentAccessGrantsPluginResource,
		plugin.NewCleanUpExpiredPingfederatePersistentSessionsPluginResource,
		plugin.NewCleanUpInactivePingfederatePersistentSessionsPluginResource,
		plugin.NewComposedAttributePluginResource,
		plugin.NewCustomPluginResource,
		plugin.NewDefaultAttributeMapperPluginResource,
		plugin.NewDefaultCleanUpExpiredPingfederatePersistentAccessGrantsPluginResource,
		plugin.NewDefaultCleanUpExpiredPingfederatePersistentSessionsPluginResource,
		plugin.NewDefaultCleanUpInactivePingfederatePersistentSessionsPluginResource,
		plugin.NewDefaultComposedAttributePluginResource,
		plugin.NewDefaultDelayPluginResource,
		plugin.NewDefaultDnMapperPluginResource,
		plugin.NewDefaultGroovyScriptedPluginResource,
		plugin.NewDefaultInternalSearchRatePluginResource,
		plugin.NewDefaultModifiablePasswordPolicyStatePluginResource,
		plugin.NewDefaultPassThroughAuthenticationPluginResource,
		plugin.NewDefaultPeriodicGcPluginResource,
		plugin.NewDefaultPeriodicStatsLoggerPluginResource,
		plugin.NewDefaultPingOnePassThroughAuthenticationPluginResource,
		plugin.NewDefaultPluggablePassThroughAuthenticationPluginResource,
		plugin.NewDefaultPurgeExpiredDataPluginResource,
		plugin.NewDefaultReferentialIntegrityPluginResource,
		plugin.NewDefaultReferralOnUpdatePluginResource,
		plugin.NewDefaultSearchShutdownPluginResource,
		plugin.NewDefaultSevenBitCleanPluginResource,
		plugin.NewDefaultSimpleToExternalBindPluginResource,
		plugin.NewDefaultSnmpSubagentPluginResource,
		plugin.NewDefaultSubOperationTimingPluginResource,
		plugin.NewDefaultThirdPartyPluginResource,
		plugin.NewDefaultUniqueAttributePluginResource,
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
		recurringtask.NewAuditDataSecurityRecurringTaskResource,
		recurringtask.NewBackupRecurringTaskResource,
		recurringtask.NewCollectSupportDataRecurringTaskResource,
		recurringtask.NewDefaultAuditDataSecurityRecurringTaskResource,
		recurringtask.NewDefaultBackupRecurringTaskResource,
		recurringtask.NewDefaultCollectSupportDataRecurringTaskResource,
		recurringtask.NewDefaultDelayRecurringTaskResource,
		recurringtask.NewDefaultEnterLockdownModeRecurringTaskResource,
		recurringtask.NewDefaultExecRecurringTaskResource,
		recurringtask.NewDefaultFileRetentionRecurringTaskResource,
		recurringtask.NewDefaultGenerateServerProfileRecurringTaskResource,
		recurringtask.NewDefaultLdifExportRecurringTaskResource,
		recurringtask.NewDefaultLeaveLockdownModeRecurringTaskResource,
		recurringtask.NewDefaultStaticallyDefinedRecurringTaskResource,
		recurringtask.NewDefaultThirdPartyRecurringTaskResource,
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
		requestcriteria.NewDefaultAggregateRequestCriteriaResource,
		requestcriteria.NewDefaultRootDseRequestCriteriaResource,
		requestcriteria.NewDefaultSimpleRequestCriteriaResource,
		requestcriteria.NewDefaultThirdPartyRequestCriteriaResource,
		requestcriteria.NewRootDseRequestCriteriaResource,
		requestcriteria.NewSimpleRequestCriteriaResource,
		requestcriteria.NewThirdPartyRequestCriteriaResource,
		restresourcetype.NewDefaultGenericRestResourceTypeResource,
		restresourcetype.NewDefaultGroupRestResourceTypeResource,
		restresourcetype.NewDefaultUserRestResourceTypeResource,
		restresourcetype.NewGenericRestResourceTypeResource,
		restresourcetype.NewGroupRestResourceTypeResource,
		restresourcetype.NewUserRestResourceTypeResource,
		scimresourcetype.NewDefaultLdapMappingScimResourceTypeResource,
		scimresourcetype.NewDefaultLdapPassThroughScimResourceTypeResource,
		scimresourcetype.NewLdapMappingScimResourceTypeResource,
		scimresourcetype.NewLdapPassThroughScimResourceTypeResource,
		saslmechanismhandler.NewAnonymousSaslMechanismHandlerResource,
		saslmechanismhandler.NewCramMd5SaslMechanismHandlerResource,
		saslmechanismhandler.NewDefaultOauthBearerSaslMechanismHandlerResource,
		saslmechanismhandler.NewDefaultThirdPartySaslMechanismHandlerResource,
		saslmechanismhandler.NewDefaultUnboundidDeliveredOtpSaslMechanismHandlerResource,
		saslmechanismhandler.NewDefaultUnboundidMsChapV2SaslMechanismHandlerResource,
		saslmechanismhandler.NewDigestMd5SaslMechanismHandlerResource,
		saslmechanismhandler.NewExternalSaslMechanismHandlerResource,
		saslmechanismhandler.NewGssapiSaslMechanismHandlerResource,
		saslmechanismhandler.NewOauthBearerSaslMechanismHandlerResource,
		saslmechanismhandler.NewPlainSaslMechanismHandlerResource,
		saslmechanismhandler.NewThirdPartySaslMechanismHandlerResource,
		saslmechanismhandler.NewUnboundidCertificatePlusPasswordSaslMechanismHandlerResource,
		saslmechanismhandler.NewUnboundidDeliveredOtpSaslMechanismHandlerResource,
		saslmechanismhandler.NewUnboundidExternalAuthSaslMechanismHandlerResource,
		saslmechanismhandler.NewUnboundidMsChapV2SaslMechanismHandlerResource,
		saslmechanismhandler.NewUnboundidTotpSaslMechanismHandlerResource,
		saslmechanismhandler.NewUnboundidYubikeyOtpSaslMechanismHandlerResource,
		serverinstance.NewAuthorizeServerInstanceResource,
		serverinstance.NewDirectoryServerInstanceResource,
		serverinstance.NewProxyServerInstanceResource,
		serverinstance.NewSyncServerInstanceResource,
		trustmanagerprovider.NewBlindTrustManagerProviderResource,
		trustmanagerprovider.NewDefaultBlindTrustManagerProviderResource,
		trustmanagerprovider.NewDefaultFileBasedTrustManagerProviderResource,
		trustmanagerprovider.NewDefaultJvmDefaultTrustManagerProviderResource,
		trustmanagerprovider.NewDefaultThirdPartyTrustManagerProviderResource,
		trustmanagerprovider.NewFileBasedTrustManagerProviderResource,
		trustmanagerprovider.NewJvmDefaultTrustManagerProviderResource,
		trustmanagerprovider.NewThirdPartyTrustManagerProviderResource,
		virtualattribute.NewConstructedVirtualAttributeResource,
		virtualattribute.NewCurrentTimeVirtualAttributeResource,
		virtualattribute.NewCustomVirtualAttributeResource,
		virtualattribute.NewDefaultConstructedVirtualAttributeResource,
		virtualattribute.NewDefaultDnJoinVirtualAttributeResource,
		virtualattribute.NewDefaultEntryDnVirtualAttributeResource,
		virtualattribute.NewDefaultEqualityJoinVirtualAttributeResource,
		virtualattribute.NewDefaultGroovyScriptedVirtualAttributeResource,
		virtualattribute.NewDefaultIdentifyReferencesVirtualAttributeResource,
		virtualattribute.NewDefaultIsMemberOfVirtualAttributeResource,
		virtualattribute.NewDefaultMemberVirtualAttributeResource,
		virtualattribute.NewDefaultMirrorVirtualAttributeResource,
		virtualattribute.NewDefaultPasswordPolicyStateJsonVirtualAttributeResource,
		virtualattribute.NewDefaultReverseDnJoinVirtualAttributeResource,
		virtualattribute.NewDefaultThirdPartyVirtualAttributeResource,
		virtualattribute.NewDefaultUserDefinedVirtualAttributeResource,
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
		webapplicationextension.NewConsoleWebApplicationExtensionResource,
		webapplicationextension.NewDefaultGenericWebApplicationExtensionResource,
		webapplicationextension.NewGenericWebApplicationExtensionResource,
	}
}
