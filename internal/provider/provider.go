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
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/alerthandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/attributesyntax"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/azureauthenticationmethod"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/backend"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/certificatemapper"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/changesubscriptionhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/cipherstreamprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/conjurauthenticationmethod"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/connectioncriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/connectionhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/datasecurityauditor"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/delegatedadminattribute"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/entrycache"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/extendedoperationhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/externalserver"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/failurelockoutaction"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/gauge"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/gaugedatasource"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/groupimplementation"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/httpservletextension"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/identitymapper"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/idtokenvalidator"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/interserverauthenticationinfo"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/keymanagerprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logfieldbehavior"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logfieldmapping"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logfieldsyntax"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logfilerotationlistener"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logpublisher"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logretentionpolicy"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logrotationpolicy"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/matchingrule"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/monitoringendpoint"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/monitorprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/notificationmanager"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/oauthtokenhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/otpdeliverymechanism"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/passphraseprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/passthroughauthenticationhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/passwordstoragescheme"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/plugin"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/recurringtask"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/requestcriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/restresourcetype"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/resultcriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/saslmechanismhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/scimresourcetype"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/searchentrycriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/searchreferencecriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/serverinstance"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/serverinstancelistener"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/synchronizationprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/tokenclaimvalidation"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/trustmanagerprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/uncachedattributecriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/uncachedentrycriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/vaultauthenticationmethod"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/velocitycontextprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/virtualattribute"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/webapplicationextension"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/workqueue"
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
		accesscontrolhandler.NewAccessControlHandlerResource,
		accountstatusnotificationhandler.NewAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewDefaultAccountStatusNotificationHandlerResource,
		accesstokenvalidator.NewAccessTokenValidatorResource,
		accesstokenvalidator.NewDefaultAccessTokenValidatorResource,
		alerthandler.NewAlertHandlerResource,
		alerthandler.NewDefaultAlertHandlerResource,
		azureauthenticationmethod.NewAzureAuthenticationMethodResource,
		azureauthenticationmethod.NewDefaultAzureAuthenticationMethodResource,
		attributesyntax.NewAttributeSyntaxResource,
		backend.NewBackendResource,
		backend.NewDefaultBackendResource,
		certificatemapper.NewCertificateMapperResource,
		certificatemapper.NewDefaultCertificateMapperResource,
		changesubscriptionhandler.NewChangeSubscriptionHandlerResource,
		changesubscriptionhandler.NewDefaultChangeSubscriptionHandlerResource,
		cipherstreamprovider.NewCipherStreamProviderResource,
		cipherstreamprovider.NewDefaultCipherStreamProviderResource,
		config.NewAlarmManagerResource,
		config.NewChangeSubscriptionResource,
		config.NewDefaultChangeSubscriptionResource,
		config.NewCipherSecretKeyResource,
		config.NewClientConnectionPolicyResource,
		config.NewDefaultClientConnectionPolicyResource,
		config.NewConsentDefinitionResource,
		config.NewDefaultConsentDefinitionResource,
		config.NewConsentDefinitionLocalizationResource,
		config.NewDefaultConsentDefinitionLocalizationResource,
		config.NewConsentServiceResource,
		config.NewCorrelatedLdapDataViewResource,
		config.NewDefaultCorrelatedLdapDataViewResource,
		config.NewConstructedAttributeResource,
		config.NewDefaultConstructedAttributeResource,
		config.NewCryptoManagerResource,
		config.NewCustomLoggedStatsResource,
		config.NewDefaultCustomLoggedStatsResource,
		config.NewDebugTargetResource,
		config.NewDefaultDebugTargetResource,
		config.NewDefaultLocationResource,
		config.NewDefaultDelegatedAdminAttributeCategoryResource,
		config.NewDelegatedAdminAttributeCategoryResource,
		config.NewDelegatedAdminResourceRightsResource,
		config.NewDefaultDelegatedAdminResourceRightsResource,
		config.NewDelegatedAdminRightsResource,
		config.NewDefaultDelegatedAdminRightsResource,
		config.NewDefaultJsonFieldConstraintsResource,
		config.NewJsonFieldConstraintsResource,
		config.NewDnMapResource,
		config.NewDefaultDnMapResource,
		config.NewGlobalConfigurationResource,
		config.NewHttpConfigurationResource,
		config.NewHttpServletCrossOriginPolicyResource,
		config.NewDefaultHttpServletCrossOriginPolicyResource,
		config.NewJsonAttributeConstraintsResource,
		config.NewDefaultJsonAttributeConstraintsResource,
		config.NewDefaultKeyPairResource,
		config.NewKeyPairResource,
		config.NewDefaultLdapCorrelationAttributePairResource,
		config.NewLdapCorrelationAttributePairResource,
		config.NewLdapSdkDebugLoggerResource,
		config.NewLicenseResource,
		config.NewLocalDbIndexResource,
		config.NewDefaultLocalDbIndexResource,
		config.NewDefaultLocalDbVlvIndexResource,
		config.NewLocalDbVlvIndexResource,
		config.NewLocalDbCompositeIndexResource,
		config.NewDefaultLocalDbCompositeIndexResource,
		config.NewLocationResource,
		config.NewMacSecretKeyResource,
		config.NewDefaultObscuredValueResource,
		config.NewObscuredValueResource,
		config.NewPasswordPolicyResource,
		config.NewDefaultPasswordPolicyResource,
		config.NewPluginRootResource,
		config.NewDefaultPrometheusMonitorAttributeMetricResource,
		config.NewPrometheusMonitorAttributeMetricResource,
		config.NewRecurringTaskChainResource,
		config.NewDefaultRecurringTaskChainResource,
		config.NewReplicationAssurancePolicyResource,
		config.NewDefaultReplicationAssurancePolicyResource,
		config.NewReplicationDomainResource,
		config.NewResultCodeMapResource,
		config.NewDefaultResultCodeMapResource,
		config.NewRootDnResource,
		config.NewRootDnUserResource,
		config.NewRootDseBackendResource,
		config.NewDefaultRootDnUserResource,
		config.NewScimAttributeResource,
		config.NewDefaultScimAttributeResource,
		config.NewScimAttributeMappingResource,
		config.NewDefaultScimAttributeMappingResource,
		config.NewScimSchemaResource,
		config.NewDefaultScimSchemaResource,
		config.NewDefaultScimSubattributeResource,
		config.NewScimSubattributeResource,
		config.NewDefaultServerGroupResource,
		config.NewServerGroupResource,
		config.NewSoftDeletePolicyResource,
		config.NewDefaultSoftDeletePolicyResource,
		config.NewTopologyAdminUserResource,
		config.NewDefaultTopologyAdminUserResource,
		config.NewTrustedCertificateResource,
		config.NewDefaultTrustedCertificateResource,
		config.NewVelocityTemplateLoaderResource,
		config.NewDefaultVelocityTemplateLoaderResource,
		conjurauthenticationmethod.NewConjurAuthenticationMethodResource,
		conjurauthenticationmethod.NewDefaultConjurAuthenticationMethodResource,
		connectioncriteria.NewConnectionCriteriaResource,
		connectioncriteria.NewDefaultConnectionCriteriaResource,
		connectionhandler.NewConnectionHandlerResource,
		connectionhandler.NewDefaultConnectionHandlerResource,
		datasecurityauditor.NewDataSecurityAuditorResource,
		datasecurityauditor.NewDefaultDataSecurityAuditorResource,
		delegatedadminattribute.NewDelegatedAdminAttributeResource,
		delegatedadminattribute.NewDefaultDelegatedAdminAttributeResource,
		entrycache.NewEntryCacheResource,
		entrycache.NewDefaultEntryCacheResource,
		extendedoperationhandler.NewExtendedOperationHandlerResource,
		extendedoperationhandler.NewDefaultExtendedOperationHandlerResource,
		externalserver.NewExternalServerResource,
		externalserver.NewDefaultExternalServerResource,
		failurelockoutaction.NewFailureLockoutActionResource,
		failurelockoutaction.NewDefaultFailureLockoutActionResource,
		gauge.NewGaugeResource,
		gauge.NewDefaultGaugeResource,
		gaugedatasource.NewGaugeDataSourceResource,
		gaugedatasource.NewDefaultGaugeDataSourceResource,
		groupimplementation.NewGroupImplementationResource,
		httpservletextension.NewHttpServletExtensionResource,
		httpservletextension.NewDefaultHttpServletExtensionResource,
		identitymapper.NewIdentityMapperResource,
		identitymapper.NewDefaultIdentityMapperResource,
		idtokenvalidator.NewIdTokenValidatorResource,
		idtokenvalidator.NewDefaultIdTokenValidatorResource,
		interserverauthenticationinfo.NewInterServerAuthenticationInfoResource,
		keymanagerprovider.NewKeyManagerProviderResource,
		keymanagerprovider.NewDefaultKeyManagerProviderResource,
		logfieldbehavior.NewLogFieldBehaviorResource,
		logfieldbehavior.NewDefaultLogFieldBehaviorResource,
		logfieldmapping.NewLogFieldMappingResource,
		logfieldmapping.NewDefaultLogFieldMappingResource,
		logfieldsyntax.NewLogFieldSyntaxResource,
		logfilerotationlistener.NewLogFileRotationListenerResource,
		logfilerotationlistener.NewDefaultLogFileRotationListenerResource,
		logpublisher.NewLogPublisherResource,
		logpublisher.NewDefaultLogPublisherResource,
		logretentionpolicy.NewLogRetentionPolicyResource,
		logretentionpolicy.NewDefaultLogRetentionPolicyResource,
		logrotationpolicy.NewLogRotationPolicyResource,
		logrotationpolicy.NewDefaultLogRotationPolicyResource,
		matchingrule.NewMatchingRuleResource,
		monitoringendpoint.NewMonitoringEndpointResource,
		monitoringendpoint.NewDefaultMonitoringEndpointResource,
		monitorprovider.NewMonitorProviderResource,
		monitorprovider.NewDefaultMonitorProviderResource,
		notificationmanager.NewDefaultNotificationManagerResource,
		notificationmanager.NewNotificationManagerResource,
		oauthtokenhandler.NewOauthTokenHandlerResource,
		oauthtokenhandler.NewDefaultOauthTokenHandlerResource,
		otpdeliverymechanism.NewDefaultOtpDeliveryMechanismResource,
		otpdeliverymechanism.NewOtpDeliveryMechanismResource,
		passphraseprovider.NewDefaultPassphraseProviderResource,
		passphraseprovider.NewPassphraseProviderResource,
		passthroughauthenticationhandler.NewPassThroughAuthenticationHandlerResource,
		passthroughauthenticationhandler.NewDefaultPassThroughAuthenticationHandlerResource,
		passwordstoragescheme.NewDefaultPasswordStorageSchemeResource,
		passwordstoragescheme.NewPasswordStorageSchemeResource,
		plugin.NewPluginResource,
		plugin.NewDefaultPluginResource,
		recurringtask.NewRecurringTaskResource,
		recurringtask.NewDefaultRecurringTaskResource,
		requestcriteria.NewRequestCriteriaResource,
		requestcriteria.NewDefaultRequestCriteriaResource,
		restresourcetype.NewRestResourceTypeResource,
		restresourcetype.NewDefaultRestResourceTypeResource,
		resultcriteria.NewResultCriteriaResource,
		resultcriteria.NewDefaultResultCriteriaResource,
		scimresourcetype.NewScimResourceTypeResource,
		scimresourcetype.NewDefaultScimResourceTypeResource,
		saslmechanismhandler.NewSaslMechanismHandlerResource,
		saslmechanismhandler.NewDefaultSaslMechanismHandlerResource,
		searchentrycriteria.NewSearchEntryCriteriaResource,
		searchentrycriteria.NewDefaultSearchEntryCriteriaResource,
		searchreferencecriteria.NewSearchReferenceCriteriaResource,
		searchreferencecriteria.NewDefaultSearchReferenceCriteriaResource,
		serverinstance.NewServerInstanceResource,
		serverinstancelistener.NewServerInstanceListenerResource,
		synchronizationprovider.NewSynchronizationProviderResource,
		tokenclaimvalidation.NewTokenClaimValidationResource,
		tokenclaimvalidation.NewDefaultTokenClaimValidationResource,
		trustmanagerprovider.NewTrustManagerProviderResource,
		trustmanagerprovider.NewDefaultTrustManagerProviderResource,
		uncachedattributecriteria.NewUncachedAttributeCriteriaResource,
		uncachedattributecriteria.NewDefaultUncachedAttributeCriteriaResource,
		uncachedentrycriteria.NewUncachedEntryCriteriaResource,
		uncachedentrycriteria.NewDefaultUncachedEntryCriteriaResource,
		vaultauthenticationmethod.NewDefaultVaultAuthenticationMethodResource,
		vaultauthenticationmethod.NewVaultAuthenticationMethodResource,
		velocitycontextprovider.NewDefaultVelocityContextProviderResource,
		velocitycontextprovider.NewVelocityContextProviderResource,
		virtualattribute.NewVirtualAttributeResource,
		virtualattribute.NewDefaultVirtualAttributeResource,
		webapplicationextension.NewWebApplicationExtensionResource,
		webapplicationextension.NewDefaultWebApplicationExtensionResource,
		workqueue.NewWorkQueueResource,
	}
}
