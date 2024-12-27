package webapplicationextension

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &webApplicationExtensionDataSource{}
	_ datasource.DataSourceWithConfigure = &webApplicationExtensionDataSource{}
)

// Create a Web Application Extension data source
func NewWebApplicationExtensionDataSource() datasource.DataSource {
	return &webApplicationExtensionDataSource{}
}

// webApplicationExtensionDataSource is the datasource implementation.
type webApplicationExtensionDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *webApplicationExtensionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_application_extension"
}

// Configure adds the provider configured client to the data source.
func (r *webApplicationExtensionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type webApplicationExtensionDataSourceModel struct {
	Id                                  types.String `tfsdk:"id"`
	Name                                types.String `tfsdk:"name"`
	Type                                types.String `tfsdk:"type"`
	SsoEnabled                          types.Bool   `tfsdk:"sso_enabled"`
	OidcClientID                        types.String `tfsdk:"oidc_client_id"`
	OidcClientSecret                    types.String `tfsdk:"oidc_client_secret"`
	OidcClientSecretPassphraseProvider  types.String `tfsdk:"oidc_client_secret_passphrase_provider"`
	OidcIssuerURL                       types.String `tfsdk:"oidc_issuer_url"`
	OidcTrustStoreFile                  types.String `tfsdk:"oidc_trust_store_file"`
	OidcTrustStoreType                  types.String `tfsdk:"oidc_trust_store_type"`
	OidcTrustStorePinPassphraseProvider types.String `tfsdk:"oidc_trust_store_pin_passphrase_provider"`
	OidcStrictHostnameVerification      types.Bool   `tfsdk:"oidc_strict_hostname_verification"`
	OidcTrustAll                        types.Bool   `tfsdk:"oidc_trust_all"`
	LdapServer                          types.String `tfsdk:"ldap_server"`
	TrustStoreFile                      types.String `tfsdk:"trust_store_file"`
	TrustStoreType                      types.String `tfsdk:"trust_store_type"`
	TrustStorePinPassphraseProvider     types.String `tfsdk:"trust_store_pin_passphrase_provider"`
	LogFile                             types.String `tfsdk:"log_file"`
	Complexity                          types.String `tfsdk:"complexity"`
	Description                         types.String `tfsdk:"description"`
	BaseContextPath                     types.String `tfsdk:"base_context_path"`
	WarFile                             types.String `tfsdk:"war_file"`
	DocumentRootDirectory               types.String `tfsdk:"document_root_directory"`
	DeploymentDescriptorFile            types.String `tfsdk:"deployment_descriptor_file"`
	TemporaryDirectory                  types.String `tfsdk:"temporary_directory"`
	InitParameter                       types.Set    `tfsdk:"init_parameter"`
}

// GetSchema defines the schema for the datasource.
func (r *webApplicationExtensionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Web Application Extension.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Web Application Extension resource. Options are ['console', 'generic']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"sso_enabled": schema.BoolAttribute{
				Description: "Indicates that SSO login into the Administrative Console is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"oidc_client_id": schema.StringAttribute{
				Description: "The client ID to use when authenticating to the OpenID Connect provider.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"oidc_client_secret": schema.StringAttribute{
				Description: "The client secret to use when authenticating to the OpenID Connect provider.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"oidc_client_secret_passphrase_provider": schema.StringAttribute{
				Description: "A passphrase provider that may be used to obtain the client secret to use when authenticating to the OpenID Connect provider.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"oidc_issuer_url": schema.StringAttribute{
				Description: "The issuer URL of the OpenID Connect provider.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"oidc_trust_store_file": schema.StringAttribute{
				Description: "Specifies the path to the truststore file used by this application to evaluate OIDC provider certificates. If this field is left blank, the default JVM trust store will be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"oidc_trust_store_type": schema.StringAttribute{
				Description: "Specifies the format for the data in the OIDC trust store file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"oidc_trust_store_pin_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider that may be used to obtain the PIN for the trust store used with OIDC providers. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"oidc_strict_hostname_verification": schema.BoolAttribute{
				Description: "Controls whether or not hostname verification is performed, which checks if the hostname of the OIDC provider matches the name(s) stored inside the certificate it provides. This property should only be set to false for testing purposes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"oidc_trust_all": schema.BoolAttribute{
				Description: "Controls whether or not this application will always trust any certificate that is presented to it, regardless of its contents. This property should only be set to true for testing purposes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ldap_server": schema.StringAttribute{
				Description: "The LDAP URL used to connect to the managed server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"trust_store_file": schema.StringAttribute{
				Description: "Specifies the path to the truststore file, which is used by this application to establish trust of managed servers.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"trust_store_type": schema.StringAttribute{
				Description: "Specifies the format for the data in the trust store file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"trust_store_pin_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider that may be used to obtain the PIN for the trust store used with managed LDAP servers. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_file": schema.StringAttribute{
				Description: "The path to the log file for the web application.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"complexity": schema.StringAttribute{
				Description: "Specifies the maximum complexity level for managed configuration elements.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Web Application Extension",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"base_context_path": schema.StringAttribute{
				Description: "Specifies the base context path that should be used by HTTP clients to reference content. The value must start with a forward slash and at least one additional character and must represent a valid HTTP context path.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"war_file": schema.StringAttribute{
				Description: "Specifies the path to a standard web application archive (WAR) file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"document_root_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory on the local filesystem containing the files to be served by this Web Application Extension. The path must exist, and it must be a directory.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"deployment_descriptor_file": schema.StringAttribute{
				Description: "Specifies the path to the deployment descriptor file when used with document-root-directory.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"temporary_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory that may be used to store temporary files such as extracted WAR files and compiled JSP files.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"init_parameter": schema.SetAttribute{
				Description: "Specifies an initialization parameter to pass into the web application during startup.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a ConsoleWebApplicationExtensionResponse object into the model struct
func readConsoleWebApplicationExtensionResponseDataSource(ctx context.Context, r *client.ConsoleWebApplicationExtensionResponse, state *webApplicationExtensionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("console")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SsoEnabled = internaltypes.BoolTypeOrNil(r.SsoEnabled)
	state.OidcClientID = internaltypes.StringTypeOrNil(r.OidcClientID, false)
	state.OidcClientSecretPassphraseProvider = internaltypes.StringTypeOrNil(r.OidcClientSecretPassphraseProvider, false)
	state.OidcIssuerURL = internaltypes.StringTypeOrNil(r.OidcIssuerURL, false)
	state.OidcTrustStoreFile = internaltypes.StringTypeOrNil(r.OidcTrustStoreFile, false)
	state.OidcTrustStoreType = internaltypes.StringTypeOrNil(r.OidcTrustStoreType, false)
	state.OidcTrustStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.OidcTrustStorePinPassphraseProvider, false)
	state.OidcStrictHostnameVerification = internaltypes.BoolTypeOrNil(r.OidcStrictHostnameVerification)
	state.OidcTrustAll = internaltypes.BoolTypeOrNil(r.OidcTrustAll)
	state.LdapServer = internaltypes.StringTypeOrNil(r.LdapServer, false)
	state.TrustStoreFile = internaltypes.StringTypeOrNil(r.TrustStoreFile, false)
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, false)
	state.TrustStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.TrustStorePinPassphraseProvider, false)
	state.LogFile = internaltypes.StringTypeOrNil(r.LogFile, false)
	state.Complexity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumwebApplicationExtensionComplexityProp(r.Complexity), false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.WarFile = internaltypes.StringTypeOrNil(r.WarFile, false)
	state.DocumentRootDirectory = internaltypes.StringTypeOrNil(r.DocumentRootDirectory, false)
	state.DeploymentDescriptorFile = internaltypes.StringTypeOrNil(r.DeploymentDescriptorFile, false)
	state.TemporaryDirectory = internaltypes.StringTypeOrNil(r.TemporaryDirectory, false)
	state.InitParameter = internaltypes.GetStringSet(r.InitParameter)
}

// Read a GenericWebApplicationExtensionResponse object into the model struct
func readGenericWebApplicationExtensionResponseDataSource(ctx context.Context, r *client.GenericWebApplicationExtensionResponse, state *webApplicationExtensionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generic")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.WarFile = internaltypes.StringTypeOrNil(r.WarFile, false)
	state.DocumentRootDirectory = internaltypes.StringTypeOrNil(r.DocumentRootDirectory, false)
	state.DeploymentDescriptorFile = internaltypes.StringTypeOrNil(r.DeploymentDescriptorFile, false)
	state.TemporaryDirectory = internaltypes.StringTypeOrNil(r.TemporaryDirectory, false)
	state.InitParameter = internaltypes.GetStringSet(r.InitParameter)
}

// Read resource information
func (r *webApplicationExtensionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state webApplicationExtensionDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.WebApplicationExtensionAPI.GetWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Web Application Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.ConsoleWebApplicationExtensionResponse != nil {
		readConsoleWebApplicationExtensionResponseDataSource(ctx, readResponse.ConsoleWebApplicationExtensionResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GenericWebApplicationExtensionResponse != nil {
		readGenericWebApplicationExtensionResponseDataSource(ctx, readResponse.GenericWebApplicationExtensionResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
