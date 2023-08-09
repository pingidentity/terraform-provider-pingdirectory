package passphraseprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &passphraseProviderDataSource{}
	_ datasource.DataSourceWithConfigure = &passphraseProviderDataSource{}
)

// Create a Passphrase Provider data source
func NewPassphraseProviderDataSource() datasource.DataSource {
	return &passphraseProviderDataSource{}
}

// passphraseProviderDataSource is the datasource implementation.
type passphraseProviderDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *passphraseProviderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_passphrase_provider"
}

// Configure adds the provider configured client to the data source.
func (r *passphraseProviderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type passphraseProviderDataSourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Type                      types.String `tfsdk:"type"`
	ExtensionClass            types.String `tfsdk:"extension_class"`
	ExtensionArgument         types.Set    `tfsdk:"extension_argument"`
	VaultExternalServer       types.String `tfsdk:"vault_external_server"`
	VaultSecretPath           types.String `tfsdk:"vault_secret_path"`
	VaultSecretFieldName      types.String `tfsdk:"vault_secret_field_name"`
	ConjurExternalServer      types.String `tfsdk:"conjur_external_server"`
	ConjurSecretRelativePath  types.String `tfsdk:"conjur_secret_relative_path"`
	PasswordFile              types.String `tfsdk:"password_file"`
	KeyVaultURI               types.String `tfsdk:"key_vault_uri"`
	AzureAuthenticationMethod types.String `tfsdk:"azure_authentication_method"`
	HttpProxyExternalServer   types.String `tfsdk:"http_proxy_external_server"`
	SecretName                types.String `tfsdk:"secret_name"`
	ObscuredValue             types.String `tfsdk:"obscured_value"`
	AwsExternalServer         types.String `tfsdk:"aws_external_server"`
	SecretID                  types.String `tfsdk:"secret_id"`
	SecretFieldName           types.String `tfsdk:"secret_field_name"`
	SecretVersionID           types.String `tfsdk:"secret_version_id"`
	SecretVersionStage        types.String `tfsdk:"secret_version_stage"`
	MaxCacheDuration          types.String `tfsdk:"max_cache_duration"`
	EnvironmentVariable       types.String `tfsdk:"environment_variable"`
	Description               types.String `tfsdk:"description"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *passphraseProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Passphrase Provider.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Passphrase Provider resource. Options are ['environment-variable', 'amazon-secrets-manager', 'obscured-value', 'azure-key-vault', 'file-based', 'conjur', 'vault', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Passphrase Provider.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Passphrase Provider. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"vault_external_server": schema.StringAttribute{
				Description: "An external server definition with information needed to connect and authenticate to the Vault instance containing the passphrase.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"vault_secret_path": schema.StringAttribute{
				Description: "The path to the desired secret in the Vault service. This will be appended to the value of the base-url property for the associated Vault external server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"vault_secret_field_name": schema.StringAttribute{
				Description: "The name of the field in the Vault secret record that contains the passphrase to use to generate the encryption key.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"conjur_external_server": schema.StringAttribute{
				Description: "An external server definition with information needed to connect and authenticate to the Conjur instance containing the passphrase.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"conjur_secret_relative_path": schema.StringAttribute{
				Description: "The portion of the path that follows the account name in the URI needed to obtain the desired secret. Any special characters in the path must be URL-encoded.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_file": schema.StringAttribute{
				Description: "The path to the file containing the passphrase.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_vault_uri": schema.StringAttribute{
				Description: "The URI that identifies the Azure Key Vault from which the secret is to be retrieved.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"azure_authentication_method": schema.StringAttribute{
				Description: "The mechanism used to authenticate to the Azure service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.0+. A reference to an HTTP proxy server that should be used for requests sent to the Azure service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"secret_name": schema.StringAttribute{
				Description: "The name of the secret to retrieve.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"obscured_value": schema.StringAttribute{
				Description: "The value to be stored in an obscured form.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"aws_external_server": schema.StringAttribute{
				Description: "The external server with information to use when interacting with the AWS Secrets Manager.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"secret_id": schema.StringAttribute{
				Description: "The Amazon Resource Name (ARN) or the user-friendly name of the secret to be retrieved.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"secret_field_name": schema.StringAttribute{
				Description: "The name of the JSON field whose value is the passphrase that will be retrieved.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"secret_version_id": schema.StringAttribute{
				Description: "The unique identifier for the version of the secret to be retrieved.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"secret_version_stage": schema.StringAttribute{
				Description: "The staging label for the version of the secret to be retrieved.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_cache_duration": schema.StringAttribute{
				Description: " When the `type` value is one of [`azure-key-vault`]: The maximum length of time that the passphrase provider may cache the passphrase that has been read from Azure Key Vault. A value of zero seconds indicates that the provider should always attempt to read the passphrase from the Azure service. When the `type` value is one of [`file-based`]: The maximum length of time that the passphrase provider may cache the passphrase that has been read from the target file. A value of zero seconds indicates that the provider should always attempt to read the passphrase from the file. When the `type` value is one of [`conjur`]: The maximum length of time that the passphrase provider may cache the passphrase that has been read from Conjur. A value of zero seconds indicates that the provider should always attempt to read the passphrase from Conjur. When the `type` value is one of [`amazon-secrets-manager`, `vault`]: The maximum length of time that the passphrase provider may cache the passphrase that has been read from Vault. A value of zero seconds indicates that the provider should always attempt to read the passphrase from Vault.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"environment_variable": schema.StringAttribute{
				Description: "The name of the environment variable that is expected to hold the passphrase.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Passphrase Provider",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Passphrase Provider is enabled for use in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a EnvironmentVariablePassphraseProviderResponse object into the model struct
func readEnvironmentVariablePassphraseProviderResponseDataSource(ctx context.Context, r *client.EnvironmentVariablePassphraseProviderResponse, state *passphraseProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("environment-variable")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnvironmentVariable = types.StringValue(r.EnvironmentVariable)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a AmazonSecretsManagerPassphraseProviderResponse object into the model struct
func readAmazonSecretsManagerPassphraseProviderResponseDataSource(ctx context.Context, r *client.AmazonSecretsManagerPassphraseProviderResponse, state *passphraseProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("amazon-secrets-manager")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AwsExternalServer = types.StringValue(r.AwsExternalServer)
	state.SecretID = types.StringValue(r.SecretID)
	state.SecretFieldName = types.StringValue(r.SecretFieldName)
	state.SecretVersionID = internaltypes.StringTypeOrNil(r.SecretVersionID, false)
	state.SecretVersionStage = internaltypes.StringTypeOrNil(r.SecretVersionStage, false)
	state.MaxCacheDuration = internaltypes.StringTypeOrNil(r.MaxCacheDuration, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ObscuredValuePassphraseProviderResponse object into the model struct
func readObscuredValuePassphraseProviderResponseDataSource(ctx context.Context, r *client.ObscuredValuePassphraseProviderResponse, state *passphraseProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("obscured-value")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a AzureKeyVaultPassphraseProviderResponse object into the model struct
func readAzureKeyVaultPassphraseProviderResponseDataSource(ctx context.Context, r *client.AzureKeyVaultPassphraseProviderResponse, state *passphraseProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("azure-key-vault")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.KeyVaultURI = types.StringValue(r.KeyVaultURI)
	state.AzureAuthenticationMethod = types.StringValue(r.AzureAuthenticationMethod)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, false)
	state.SecretName = types.StringValue(r.SecretName)
	state.MaxCacheDuration = internaltypes.StringTypeOrNil(r.MaxCacheDuration, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a FileBasedPassphraseProviderResponse object into the model struct
func readFileBasedPassphraseProviderResponseDataSource(ctx context.Context, r *client.FileBasedPassphraseProviderResponse, state *passphraseProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PasswordFile = types.StringValue(r.PasswordFile)
	state.MaxCacheDuration = internaltypes.StringTypeOrNil(r.MaxCacheDuration, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ConjurPassphraseProviderResponse object into the model struct
func readConjurPassphraseProviderResponseDataSource(ctx context.Context, r *client.ConjurPassphraseProviderResponse, state *passphraseProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("conjur")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConjurExternalServer = types.StringValue(r.ConjurExternalServer)
	state.ConjurSecretRelativePath = types.StringValue(r.ConjurSecretRelativePath)
	state.MaxCacheDuration = internaltypes.StringTypeOrNil(r.MaxCacheDuration, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a VaultPassphraseProviderResponse object into the model struct
func readVaultPassphraseProviderResponseDataSource(ctx context.Context, r *client.VaultPassphraseProviderResponse, state *passphraseProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("vault")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.VaultExternalServer = types.StringValue(r.VaultExternalServer)
	state.VaultSecretPath = types.StringValue(r.VaultSecretPath)
	state.VaultSecretFieldName = types.StringValue(r.VaultSecretFieldName)
	state.MaxCacheDuration = internaltypes.StringTypeOrNil(r.MaxCacheDuration, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyPassphraseProviderResponse object into the model struct
func readThirdPartyPassphraseProviderResponseDataSource(ctx context.Context, r *client.ThirdPartyPassphraseProviderResponse, state *passphraseProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *passphraseProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state passphraseProviderDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PassphraseProviderApi.GetPassphraseProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Passphrase Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.EnvironmentVariablePassphraseProviderResponse != nil {
		readEnvironmentVariablePassphraseProviderResponseDataSource(ctx, readResponse.EnvironmentVariablePassphraseProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AmazonSecretsManagerPassphraseProviderResponse != nil {
		readAmazonSecretsManagerPassphraseProviderResponseDataSource(ctx, readResponse.AmazonSecretsManagerPassphraseProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ObscuredValuePassphraseProviderResponse != nil {
		readObscuredValuePassphraseProviderResponseDataSource(ctx, readResponse.ObscuredValuePassphraseProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AzureKeyVaultPassphraseProviderResponse != nil {
		readAzureKeyVaultPassphraseProviderResponseDataSource(ctx, readResponse.AzureKeyVaultPassphraseProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedPassphraseProviderResponse != nil {
		readFileBasedPassphraseProviderResponseDataSource(ctx, readResponse.FileBasedPassphraseProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ConjurPassphraseProviderResponse != nil {
		readConjurPassphraseProviderResponseDataSource(ctx, readResponse.ConjurPassphraseProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.VaultPassphraseProviderResponse != nil {
		readVaultPassphraseProviderResponseDataSource(ctx, readResponse.VaultPassphraseProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPassphraseProviderResponse != nil {
		readThirdPartyPassphraseProviderResponseDataSource(ctx, readResponse.ThirdPartyPassphraseProviderResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
