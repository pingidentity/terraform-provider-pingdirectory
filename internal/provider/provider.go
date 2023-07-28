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
	client9300 "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/accesscontrolhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/accesstokenvalidator"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/accountstatusnotificationhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/alarmmanager"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/alerthandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/attributesyntax"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/azureauthenticationmethod"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/backend"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/certificatemapper"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/changesubscription"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/changesubscriptionhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/ciphersecretkey"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/cipherstreamprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/clientconnectionpolicy"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/conjurauthenticationmethod"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/connectioncriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/connectionhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/consentdefinition"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/consentdefinitionlocalization"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/consentservice"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/constructedattribute"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/correlatedldapdataview"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/cryptomanager"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/customloggedstats"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/datasecurityauditor"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/debugtarget"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/delegatedadminattribute"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/delegatedadminattributecategory"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/delegatedadmincorrelatedrestresource"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/delegatedadminresourcerights"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/delegatedadminrights"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/dnmap"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/entrycache"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/extendedoperationhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/externalserver"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/failurelockoutaction"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/gauge"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/gaugedatasource"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/globalconfiguration"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/groupimplementation"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/httpconfiguration"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/httpservletcrossoriginpolicy"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/httpservletextension"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/identitymapper"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/idtokenvalidator"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/interserverauthenticationinfo"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/jsonattributeconstraints"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/jsonfieldconstraints"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/keymanagerprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/keypair"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/ldapcorrelationattributepair"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/ldapsdkdebuglogger"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/license"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/localdbcompositeindex"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/localdbindex"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/localdbvlvindex"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/location"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logfieldbehavior"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logfieldmapping"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logfieldsyntax"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logfilerotationlistener"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logpublisher"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logretentionpolicy"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/logrotationpolicy"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/macsecretkey"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/matchingrule"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/monitoringendpoint"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/monitorprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/notificationmanager"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/oauthtokenhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/obscuredvalue"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/otpdeliverymechanism"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/passphraseprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/passthroughauthenticationhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/passwordgenerator"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/passwordpolicy"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/passwordstoragescheme"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/passwordvalidator"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/plugin"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/pluginroot"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/prometheusmonitorattributemetric"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/recurringtask"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/recurringtaskchain"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/replicationassurancepolicy"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/replicationdomain"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/replicationserver"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/requestcriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/restresourcetype"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/resultcodemap"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/resultcriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/rootdn"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/rootdnuser"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/rootdsebackend"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/saslmechanismhandler"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/scimattribute"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/scimattributemapping"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/scimresourcetype"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/scimschema"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/scimsubattribute"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/searchentrycriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/searchreferencecriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/sensitiveattribute"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/servergroup"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/serverinstance"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/serverinstancelistener"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/softdeletepolicy"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/synchronizationprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/tokenclaimvalidation"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/topologyadminuser"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/trustedcertificate"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/trustmanagerprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/uncachedattributecriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/uncachedentrycriteria"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/vaultauthenticationmethod"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/velocitycontextprovider"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config/velocitytemplateloader"
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
	clientConfig9200 := client9300.NewConfiguration()
	clientConfig9200.Servers = client9300.ServerConfigurations{
		{
			URL: httpsHost + "/config",
		},
	}
	clientConfig9200.HTTPClient = httpClient
	resourceConfig.ApiClientV9300 = client9300.NewAPIClient(clientConfig9200)

	resp.ResourceData = resourceConfig
	resp.DataSourceData = resourceConfig
	tflog.Info(ctx, "Configured PingDirectory client", map[string]interface{}{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *pingdirectoryProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		accesscontrolhandler.NewAccessControlHandlerDataSource,
		accesstokenvalidator.NewAccessTokenValidatorDataSource,
		accountstatusnotificationhandler.NewAccountStatusNotificationHandlerDataSource,
		alarmmanager.NewAlarmManagerDataSource,
		alerthandler.NewAlertHandlerDataSource,
		attributesyntax.NewAttributeSyntaxDataSource,
		azureauthenticationmethod.NewAzureAuthenticationMethodDataSource,
		backend.NewBackendDataSource,
		certificatemapper.NewCertificateMapperDataSource,
		changesubscription.NewChangeSubscriptionDataSource,
		changesubscriptionhandler.NewChangeSubscriptionHandlerDataSource,
		ciphersecretkey.NewCipherSecretKeyDataSource,
		cipherstreamprovider.NewCipherStreamProviderDataSource,
		clientconnectionpolicy.NewClientConnectionPolicyDataSource,
		conjurauthenticationmethod.NewConjurAuthenticationMethodDataSource,
		connectioncriteria.NewConnectionCriteriaDataSource,
		connectionhandler.NewConnectionHandlerDataSource,
		consentdefinition.NewConsentDefinitionDataSource,
		consentdefinitionlocalization.NewConsentDefinitionLocalizationDataSource,
		consentservice.NewConsentServiceDataSource,
		constructedattribute.NewConstructedAttributeDataSource,
		correlatedldapdataview.NewCorrelatedLdapDataViewDataSource,
		cryptomanager.NewCryptoManagerDataSource,
		customloggedstats.NewCustomLoggedStatsDataSource,
		datasecurityauditor.NewDataSecurityAuditorDataSource,
		debugtarget.NewDebugTargetDataSource,
		delegatedadminattribute.NewDelegatedAdminAttributeDataSource,
		delegatedadminattributecategory.NewDelegatedAdminAttributeCategoryDataSource,
		delegatedadmincorrelatedrestresource.NewDelegatedAdminCorrelatedRestResourceDataSource,
		delegatedadminresourcerights.NewDelegatedAdminResourceRightsDataSource,
		delegatedadminrights.NewDelegatedAdminRightsDataSource,
		dnmap.NewDnMapDataSource,
		entrycache.NewEntryCacheDataSource,
		extendedoperationhandler.NewExtendedOperationHandlerDataSource,
		externalserver.NewExternalServerDataSource,
		failurelockoutaction.NewFailureLockoutActionDataSource,
		gauge.NewGaugeDataSource,
		gauge.NewGaugesDataSource,
		gaugedatasource.NewGaugeDataSourceDataSource,
		globalconfiguration.NewGlobalConfigurationDataSource,
		groupimplementation.NewGroupImplementationDataSource,
		httpconfiguration.NewHttpConfigurationDataSource,
		httpservletcrossoriginpolicy.NewHttpServletCrossOriginPolicyDataSource,
		httpservletextension.NewHttpServletExtensionDataSource,
		identitymapper.NewIdentityMapperDataSource,
		idtokenvalidator.NewIdTokenValidatorDataSource,
		interserverauthenticationinfo.NewInterServerAuthenticationInfoDataSource,
		jsonattributeconstraints.NewJsonAttributeConstraintsDataSource,
		jsonfieldconstraints.NewJsonFieldConstraintsDataSource,
		keymanagerprovider.NewKeyManagerProviderDataSource,
		keypair.NewKeyPairDataSource,
		ldapcorrelationattributepair.NewLdapCorrelationAttributePairDataSource,
		ldapsdkdebuglogger.NewLdapSdkDebugLoggerDataSource,
		license.NewLicenseDataSource,
		localdbcompositeindex.NewLocalDbCompositeIndexDataSource,
		localdbindex.NewLocalDbIndexDataSource,
		localdbvlvindex.NewLocalDbVlvIndexDataSource,
		location.NewLocationDataSource,
		logfieldbehavior.NewLogFieldBehaviorDataSource,
		logfieldmapping.NewLogFieldMappingDataSource,
		logfieldsyntax.NewLogFieldSyntaxDataSource,
		logfilerotationlistener.NewLogFileRotationListenerDataSource,
		logpublisher.NewLogPublisherDataSource,
		logretentionpolicy.NewLogRetentionPolicyDataSource,
		logrotationpolicy.NewLogRotationPolicyDataSource,
		macsecretkey.NewMacSecretKeyDataSource,
		matchingrule.NewMatchingRuleDataSource,
		monitoringendpoint.NewMonitoringEndpointDataSource,
		monitorprovider.NewMonitorProviderDataSource,
		notificationmanager.NewNotificationManagerDataSource,
		oauthtokenhandler.NewOauthTokenHandlerDataSource,
		obscuredvalue.NewObscuredValueDataSource,
		otpdeliverymechanism.NewOtpDeliveryMechanismDataSource,
		passphraseprovider.NewPassphraseProviderDataSource,
		passthroughauthenticationhandler.NewPassThroughAuthenticationHandlerDataSource,
		passwordgenerator.NewPasswordGeneratorDataSource,
		passwordpolicy.NewPasswordPolicyDataSource,
		passwordstoragescheme.NewPasswordStorageSchemeDataSource,
		passwordvalidator.NewPasswordValidatorDataSource,
		plugin.NewPluginDataSource,
		pluginroot.NewPluginRootDataSource,
		prometheusmonitorattributemetric.NewPrometheusMonitorAttributeMetricDataSource,
		recurringtask.NewRecurringTaskDataSource,
		recurringtaskchain.NewRecurringTaskChainDataSource,
		replicationassurancepolicy.NewReplicationAssurancePolicyDataSource,
		replicationdomain.NewReplicationDomainDataSource,
		replicationserver.NewReplicationServerDataSource,
		requestcriteria.NewRequestCriteriaDataSource,
		restresourcetype.NewRestResourceTypeDataSource,
		resultcodemap.NewResultCodeMapDataSource,
		resultcriteria.NewResultCriteriaDataSource,
		rootdn.NewRootDnDataSource,
		rootdnuser.NewRootDnUserDataSource,
		rootdsebackend.NewRootDseBackendDataSource,
		saslmechanismhandler.NewSaslMechanismHandlerDataSource,
		scimattribute.NewScimAttributeDataSource,
		scimattributemapping.NewScimAttributeMappingDataSource,
		scimresourcetype.NewScimResourceTypeDataSource,
		scimschema.NewScimSchemaDataSource,
		scimsubattribute.NewScimSubattributeDataSource,
		searchentrycriteria.NewSearchEntryCriteriaDataSource,
		searchreferencecriteria.NewSearchReferenceCriteriaDataSource,
		sensitiveattribute.NewSensitiveAttributeDataSource,
		servergroup.NewServerGroupDataSource,
		serverinstance.NewServerInstanceDataSource,
		serverinstancelistener.NewServerInstanceListenerDataSource,
		softdeletepolicy.NewSoftDeletePolicyDataSource,
		synchronizationprovider.NewSynchronizationProviderDataSource,
		tokenclaimvalidation.NewTokenClaimValidationDataSource,
		topologyadminuser.NewTopologyAdminUserDataSource,
		trustedcertificate.NewTrustedCertificateDataSource,
		trustmanagerprovider.NewTrustManagerProviderDataSource,
		uncachedattributecriteria.NewUncachedAttributeCriteriaDataSource,
		uncachedentrycriteria.NewUncachedEntryCriteriaDataSource,
		vaultauthenticationmethod.NewVaultAuthenticationMethodDataSource,
		velocitycontextprovider.NewVelocityContextProviderDataSource,
		velocitytemplateloader.NewVelocityTemplateLoaderDataSource,
		virtualattribute.NewVirtualAttributeDataSource,
		webapplicationextension.NewWebApplicationExtensionDataSource,
		workqueue.NewWorkQueueDataSource,
	}
}

// Resources defines the resources implemented in the provider.
// Maintain alphabetical order for ease of management
func (p *pingdirectoryProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		accesscontrolhandler.NewAccessControlHandlerResource,
		accesstokenvalidator.NewAccessTokenValidatorResource,
		accesstokenvalidator.NewDefaultAccessTokenValidatorResource,
		accountstatusnotificationhandler.NewAccountStatusNotificationHandlerResource,
		accountstatusnotificationhandler.NewDefaultAccountStatusNotificationHandlerResource,
		alarmmanager.NewAlarmManagerResource,
		alerthandler.NewAlertHandlerResource,
		alerthandler.NewDefaultAlertHandlerResource,
		attributesyntax.NewAttributeSyntaxResource,
		azureauthenticationmethod.NewAzureAuthenticationMethodResource,
		azureauthenticationmethod.NewDefaultAzureAuthenticationMethodResource,
		backend.NewBackendResource,
		backend.NewDefaultBackendResource,
		certificatemapper.NewCertificateMapperResource,
		certificatemapper.NewDefaultCertificateMapperResource,
		changesubscription.NewChangeSubscriptionResource,
		changesubscription.NewDefaultChangeSubscriptionResource,
		changesubscriptionhandler.NewChangeSubscriptionHandlerResource,
		changesubscriptionhandler.NewDefaultChangeSubscriptionHandlerResource,
		ciphersecretkey.NewCipherSecretKeyResource,
		cipherstreamprovider.NewCipherStreamProviderResource,
		cipherstreamprovider.NewDefaultCipherStreamProviderResource,
		clientconnectionpolicy.NewClientConnectionPolicyResource,
		clientconnectionpolicy.NewDefaultClientConnectionPolicyResource,
		conjurauthenticationmethod.NewConjurAuthenticationMethodResource,
		conjurauthenticationmethod.NewDefaultConjurAuthenticationMethodResource,
		connectioncriteria.NewConnectionCriteriaResource,
		connectioncriteria.NewDefaultConnectionCriteriaResource,
		connectionhandler.NewConnectionHandlerResource,
		connectionhandler.NewDefaultConnectionHandlerResource,
		consentdefinition.NewConsentDefinitionResource,
		consentdefinition.NewDefaultConsentDefinitionResource,
		consentdefinitionlocalization.NewConsentDefinitionLocalizationResource,
		consentdefinitionlocalization.NewDefaultConsentDefinitionLocalizationResource,
		consentservice.NewConsentServiceResource,
		constructedattribute.NewConstructedAttributeResource,
		constructedattribute.NewDefaultConstructedAttributeResource,
		correlatedldapdataview.NewCorrelatedLdapDataViewResource,
		correlatedldapdataview.NewDefaultCorrelatedLdapDataViewResource,
		cryptomanager.NewCryptoManagerResource,
		customloggedstats.NewCustomLoggedStatsResource,
		customloggedstats.NewDefaultCustomLoggedStatsResource,
		datasecurityauditor.NewDataSecurityAuditorResource,
		datasecurityauditor.NewDefaultDataSecurityAuditorResource,
		debugtarget.NewDebugTargetResource,
		debugtarget.NewDefaultDebugTargetResource,
		delegatedadminattribute.NewDefaultDelegatedAdminAttributeResource,
		delegatedadminattribute.NewDelegatedAdminAttributeResource,
		delegatedadminattributecategory.NewDefaultDelegatedAdminAttributeCategoryResource,
		delegatedadminattributecategory.NewDelegatedAdminAttributeCategoryResource,
		delegatedadmincorrelatedrestresource.NewDefaultDelegatedAdminCorrelatedRestResourceResource,
		delegatedadmincorrelatedrestresource.NewDelegatedAdminCorrelatedRestResourceResource,
		delegatedadminresourcerights.NewDefaultDelegatedAdminResourceRightsResource,
		delegatedadminresourcerights.NewDelegatedAdminResourceRightsResource,
		delegatedadminrights.NewDefaultDelegatedAdminRightsResource,
		delegatedadminrights.NewDelegatedAdminRightsResource,
		dnmap.NewDefaultDnMapResource,
		dnmap.NewDnMapResource,
		entrycache.NewDefaultEntryCacheResource,
		entrycache.NewEntryCacheResource,
		extendedoperationhandler.NewDefaultExtendedOperationHandlerResource,
		extendedoperationhandler.NewExtendedOperationHandlerResource,
		externalserver.NewDefaultExternalServerResource,
		externalserver.NewExternalServerResource,
		failurelockoutaction.NewDefaultFailureLockoutActionResource,
		failurelockoutaction.NewFailureLockoutActionResource,
		gauge.NewDefaultGaugeResource,
		gauge.NewGaugeResource,
		gaugedatasource.NewDefaultGaugeDataSourceResource,
		gaugedatasource.NewGaugeDataSourceResource,
		globalconfiguration.NewGlobalConfigurationResource,
		groupimplementation.NewGroupImplementationResource,
		httpconfiguration.NewHttpConfigurationResource,
		httpservletcrossoriginpolicy.NewDefaultHttpServletCrossOriginPolicyResource,
		httpservletcrossoriginpolicy.NewHttpServletCrossOriginPolicyResource,
		httpservletextension.NewDefaultHttpServletExtensionResource,
		httpservletextension.NewHttpServletExtensionResource,
		identitymapper.NewDefaultIdentityMapperResource,
		identitymapper.NewIdentityMapperResource,
		idtokenvalidator.NewDefaultIdTokenValidatorResource,
		idtokenvalidator.NewIdTokenValidatorResource,
		interserverauthenticationinfo.NewInterServerAuthenticationInfoResource,
		jsonattributeconstraints.NewDefaultJsonAttributeConstraintsResource,
		jsonattributeconstraints.NewJsonAttributeConstraintsResource,
		jsonfieldconstraints.NewDefaultJsonFieldConstraintsResource,
		jsonfieldconstraints.NewJsonFieldConstraintsResource,
		keymanagerprovider.NewDefaultKeyManagerProviderResource,
		keymanagerprovider.NewKeyManagerProviderResource,
		keypair.NewDefaultKeyPairResource,
		keypair.NewKeyPairResource,
		ldapcorrelationattributepair.NewDefaultLdapCorrelationAttributePairResource,
		ldapcorrelationattributepair.NewLdapCorrelationAttributePairResource,
		ldapsdkdebuglogger.NewLdapSdkDebugLoggerResource,
		license.NewLicenseResource,
		localdbcompositeindex.NewDefaultLocalDbCompositeIndexResource,
		localdbcompositeindex.NewLocalDbCompositeIndexResource,
		localdbindex.NewDefaultLocalDbIndexResource,
		localdbindex.NewLocalDbIndexResource,
		localdbvlvindex.NewDefaultLocalDbVlvIndexResource,
		localdbvlvindex.NewLocalDbVlvIndexResource,
		location.NewDefaultLocationResource,
		location.NewLocationResource,
		logfieldbehavior.NewDefaultLogFieldBehaviorResource,
		logfieldbehavior.NewLogFieldBehaviorResource,
		logfieldmapping.NewDefaultLogFieldMappingResource,
		logfieldmapping.NewLogFieldMappingResource,
		logfieldsyntax.NewLogFieldSyntaxResource,
		logfilerotationlistener.NewDefaultLogFileRotationListenerResource,
		logfilerotationlistener.NewLogFileRotationListenerResource,
		logpublisher.NewDefaultLogPublisherResource,
		logpublisher.NewLogPublisherResource,
		logretentionpolicy.NewDefaultLogRetentionPolicyResource,
		logretentionpolicy.NewLogRetentionPolicyResource,
		logrotationpolicy.NewDefaultLogRotationPolicyResource,
		logrotationpolicy.NewLogRotationPolicyResource,
		macsecretkey.NewMacSecretKeyResource,
		matchingrule.NewMatchingRuleResource,
		monitoringendpoint.NewDefaultMonitoringEndpointResource,
		monitoringendpoint.NewMonitoringEndpointResource,
		monitorprovider.NewDefaultMonitorProviderResource,
		monitorprovider.NewMonitorProviderResource,
		notificationmanager.NewDefaultNotificationManagerResource,
		notificationmanager.NewNotificationManagerResource,
		oauthtokenhandler.NewDefaultOauthTokenHandlerResource,
		oauthtokenhandler.NewOauthTokenHandlerResource,
		obscuredvalue.NewDefaultObscuredValueResource,
		obscuredvalue.NewObscuredValueResource,
		otpdeliverymechanism.NewDefaultOtpDeliveryMechanismResource,
		otpdeliverymechanism.NewOtpDeliveryMechanismResource,
		passphraseprovider.NewDefaultPassphraseProviderResource,
		passphraseprovider.NewPassphraseProviderResource,
		passthroughauthenticationhandler.NewDefaultPassThroughAuthenticationHandlerResource,
		passthroughauthenticationhandler.NewPassThroughAuthenticationHandlerResource,
		passwordgenerator.NewDefaultPasswordGeneratorResource,
		passwordgenerator.NewPasswordGeneratorResource,
		passwordpolicy.NewDefaultPasswordPolicyResource,
		passwordpolicy.NewPasswordPolicyResource,
		passwordstoragescheme.NewDefaultPasswordStorageSchemeResource,
		passwordstoragescheme.NewPasswordStorageSchemeResource,
		passwordvalidator.NewDefaultPasswordValidatorResource,
		passwordvalidator.NewPasswordValidatorResource,
		plugin.NewDefaultPluginResource,
		plugin.NewPluginResource,
		pluginroot.NewPluginRootResource,
		prometheusmonitorattributemetric.NewDefaultPrometheusMonitorAttributeMetricResource,
		prometheusmonitorattributemetric.NewPrometheusMonitorAttributeMetricResource,
		recurringtask.NewDefaultRecurringTaskResource,
		recurringtask.NewRecurringTaskResource,
		recurringtaskchain.NewDefaultRecurringTaskChainResource,
		recurringtaskchain.NewRecurringTaskChainResource,
		replicationassurancepolicy.NewDefaultReplicationAssurancePolicyResource,
		replicationassurancepolicy.NewReplicationAssurancePolicyResource,
		replicationdomain.NewReplicationDomainResource,
		replicationserver.NewReplicationServerResource,
		requestcriteria.NewDefaultRequestCriteriaResource,
		requestcriteria.NewRequestCriteriaResource,
		restresourcetype.NewDefaultRestResourceTypeResource,
		restresourcetype.NewRestResourceTypeResource,
		resultcodemap.NewDefaultResultCodeMapResource,
		resultcodemap.NewResultCodeMapResource,
		resultcriteria.NewDefaultResultCriteriaResource,
		resultcriteria.NewResultCriteriaResource,
		rootdn.NewRootDnResource,
		rootdnuser.NewDefaultRootDnUserResource,
		rootdnuser.NewRootDnUserResource,
		rootdsebackend.NewRootDseBackendResource,
		saslmechanismhandler.NewDefaultSaslMechanismHandlerResource,
		saslmechanismhandler.NewSaslMechanismHandlerResource,
		scimattribute.NewDefaultScimAttributeResource,
		scimattribute.NewScimAttributeResource,
		scimattributemapping.NewDefaultScimAttributeMappingResource,
		scimattributemapping.NewScimAttributeMappingResource,
		scimresourcetype.NewDefaultScimResourceTypeResource,
		scimresourcetype.NewScimResourceTypeResource,
		scimschema.NewDefaultScimSchemaResource,
		scimschema.NewScimSchemaResource,
		scimsubattribute.NewDefaultScimSubattributeResource,
		scimsubattribute.NewScimSubattributeResource,
		searchentrycriteria.NewDefaultSearchEntryCriteriaResource,
		searchentrycriteria.NewSearchEntryCriteriaResource,
		searchreferencecriteria.NewDefaultSearchReferenceCriteriaResource,
		searchreferencecriteria.NewSearchReferenceCriteriaResource,
		sensitiveattribute.NewDefaultSensitiveAttributeResource,
		sensitiveattribute.NewSensitiveAttributeResource,
		servergroup.NewDefaultServerGroupResource,
		servergroup.NewServerGroupResource,
		serverinstance.NewServerInstanceResource,
		serverinstancelistener.NewServerInstanceListenerResource,
		softdeletepolicy.NewDefaultSoftDeletePolicyResource,
		softdeletepolicy.NewSoftDeletePolicyResource,
		synchronizationprovider.NewSynchronizationProviderResource,
		tokenclaimvalidation.NewDefaultTokenClaimValidationResource,
		tokenclaimvalidation.NewTokenClaimValidationResource,
		topologyadminuser.NewDefaultTopologyAdminUserResource,
		topologyadminuser.NewTopologyAdminUserResource,
		trustedcertificate.NewDefaultTrustedCertificateResource,
		trustedcertificate.NewTrustedCertificateResource,
		trustmanagerprovider.NewDefaultTrustManagerProviderResource,
		trustmanagerprovider.NewTrustManagerProviderResource,
		uncachedattributecriteria.NewDefaultUncachedAttributeCriteriaResource,
		uncachedattributecriteria.NewUncachedAttributeCriteriaResource,
		uncachedentrycriteria.NewDefaultUncachedEntryCriteriaResource,
		uncachedentrycriteria.NewUncachedEntryCriteriaResource,
		vaultauthenticationmethod.NewDefaultVaultAuthenticationMethodResource,
		vaultauthenticationmethod.NewVaultAuthenticationMethodResource,
		velocitycontextprovider.NewDefaultVelocityContextProviderResource,
		velocitycontextprovider.NewVelocityContextProviderResource,
		velocitytemplateloader.NewDefaultVelocityTemplateLoaderResource,
		velocitytemplateloader.NewVelocityTemplateLoaderResource,
		virtualattribute.NewDefaultVirtualAttributeResource,
		virtualattribute.NewVirtualAttributeResource,
		webapplicationextension.NewDefaultWebApplicationExtensionResource,
		webapplicationextension.NewWebApplicationExtensionResource,
		workqueue.NewWorkQueueResource,
	}
}
