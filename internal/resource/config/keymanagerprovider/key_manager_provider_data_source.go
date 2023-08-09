package keymanagerprovider

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
	_ datasource.DataSource              = &keyManagerProviderDataSource{}
	_ datasource.DataSourceWithConfigure = &keyManagerProviderDataSource{}
)

// Create a Key Manager Provider data source
func NewKeyManagerProviderDataSource() datasource.DataSource {
	return &keyManagerProviderDataSource{}
}

// keyManagerProviderDataSource is the datasource implementation.
type keyManagerProviderDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *keyManagerProviderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_manager_provider"
}

// Configure adds the provider configured client to the data source.
func (r *keyManagerProviderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type keyManagerProviderDataSourceModel struct {
	Id                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	Type                            types.String `tfsdk:"type"`
	ExtensionClass                  types.String `tfsdk:"extension_class"`
	ExtensionArgument               types.Set    `tfsdk:"extension_argument"`
	Pkcs11ProviderClass             types.String `tfsdk:"pkcs11_provider_class"`
	Pkcs11ProviderConfigurationFile types.String `tfsdk:"pkcs11_provider_configuration_file"`
	Pkcs11KeyStoreType              types.String `tfsdk:"pkcs11_key_store_type"`
	Pkcs11MaxCacheDuration          types.String `tfsdk:"pkcs11_max_cache_duration"`
	KeyStoreFile                    types.String `tfsdk:"key_store_file"`
	KeyStoreType                    types.String `tfsdk:"key_store_type"`
	KeyStorePin                     types.String `tfsdk:"key_store_pin"`
	KeyStorePinFile                 types.String `tfsdk:"key_store_pin_file"`
	KeyStorePinPassphraseProvider   types.String `tfsdk:"key_store_pin_passphrase_provider"`
	PrivateKeyPin                   types.String `tfsdk:"private_key_pin"`
	PrivateKeyPinFile               types.String `tfsdk:"private_key_pin_file"`
	PrivateKeyPinPassphraseProvider types.String `tfsdk:"private_key_pin_passphrase_provider"`
	Description                     types.String `tfsdk:"description"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *keyManagerProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Key Manager Provider.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Key Manager Provider resource. Options are ['file-based', 'custom', 'pkcs11', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Key Manager Provider.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Key Manager Provider. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"pkcs11_provider_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java security provider class that implements support for interacting with PKCS #11 tokens.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"pkcs11_provider_configuration_file": schema.StringAttribute{
				Description: "The path to the file to use to configure the security provider that implements support for interacting with PKCS #11 tokens.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"pkcs11_key_store_type": schema.StringAttribute{
				Description: "The key store type to use when obtaining an instance of a key store for interacting with a PKCS #11 token.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"pkcs11_max_cache_duration": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.1+. The maximum length of time that data retrieved from PKCS #11 tokens may be cached for reuse. Caching might be necessary if there is noticable latency when accessing the token, for example if the token uses a remote key store. A value of zero milliseconds indicates that no caching should be performed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_store_file": schema.StringAttribute{
				Description: "Specifies the path to the file that contains the private key information. This may be an absolute path, or a path that is relative to the Directory Server instance root.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_store_type": schema.StringAttribute{
				Description: "Specifies the format for the data in the key store file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_store_pin": schema.StringAttribute{
				Description: " When the `type` attribute is set to `file-based`: Specifies the PIN needed to access the File Based Key Manager Provider. When the `type` attribute is set to `pkcs11`: Specifies the PIN needed to access the PKCS11 Key Manager Provider.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"key_store_pin_file": schema.StringAttribute{
				Description: " When the `type` attribute is set to `file-based`: Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the File Based Key Manager Provider. When the `type` attribute is set to `pkcs11`: Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the PKCS11 Key Manager Provider.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_store_pin_passphrase_provider": schema.StringAttribute{
				Description: " When the `type` attribute is set to `file-based`: The passphrase provider to use to obtain the clear-text PIN needed to access the File Based Key Manager Provider. When the `type` attribute is set to `pkcs11`: The passphrase provider to use to obtain the clear-text PIN needed to access the PKCS11 Key Manager Provider.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"private_key_pin": schema.StringAttribute{
				Description: "Specifies the clear-text PIN needed to access the File Based Key Manager Provider private key. If no private key PIN is specified the PIN defaults to the key store PIN.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"private_key_pin_file": schema.StringAttribute{
				Description: "Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the File Based Key Manager Provider private key. If no private key PIN is specified the PIN defaults to the key store PIN.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"private_key_pin_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider to use to obtain the clear-text PIN needed to access the File Based Key Manager Provider private key. If no private key PIN is specified the PIN defaults to the key store PIN.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Key Manager Provider",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Key Manager Provider is enabled for use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a FileBasedKeyManagerProviderResponse object into the model struct
func readFileBasedKeyManagerProviderResponseDataSource(ctx context.Context, r *client.FileBasedKeyManagerProviderResponse, state *keyManagerProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.KeyStoreFile = types.StringValue(r.KeyStoreFile)
	state.KeyStoreType = internaltypes.StringTypeOrNil(r.KeyStoreType, false)
	state.KeyStorePinFile = internaltypes.StringTypeOrNil(r.KeyStorePinFile, false)
	state.KeyStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.KeyStorePinPassphraseProvider, false)
	state.PrivateKeyPinFile = internaltypes.StringTypeOrNil(r.PrivateKeyPinFile, false)
	state.PrivateKeyPinPassphraseProvider = internaltypes.StringTypeOrNil(r.PrivateKeyPinPassphraseProvider, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a CustomKeyManagerProviderResponse object into the model struct
func readCustomKeyManagerProviderResponseDataSource(ctx context.Context, r *client.CustomKeyManagerProviderResponse, state *keyManagerProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a Pkcs11KeyManagerProviderResponse object into the model struct
func readPkcs11KeyManagerProviderResponseDataSource(ctx context.Context, r *client.Pkcs11KeyManagerProviderResponse, state *keyManagerProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("pkcs11")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Pkcs11ProviderClass = internaltypes.StringTypeOrNil(r.Pkcs11ProviderClass, false)
	state.Pkcs11ProviderConfigurationFile = internaltypes.StringTypeOrNil(r.Pkcs11ProviderConfigurationFile, false)
	state.Pkcs11KeyStoreType = internaltypes.StringTypeOrNil(r.Pkcs11KeyStoreType, false)
	state.Pkcs11MaxCacheDuration = internaltypes.StringTypeOrNil(r.Pkcs11MaxCacheDuration, false)
	state.KeyStorePinFile = internaltypes.StringTypeOrNil(r.KeyStorePinFile, false)
	state.KeyStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.KeyStorePinPassphraseProvider, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyKeyManagerProviderResponse object into the model struct
func readThirdPartyKeyManagerProviderResponseDataSource(ctx context.Context, r *client.ThirdPartyKeyManagerProviderResponse, state *keyManagerProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *keyManagerProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state keyManagerProviderDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.KeyManagerProviderApi.GetKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Key Manager Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.FileBasedKeyManagerProviderResponse != nil {
		readFileBasedKeyManagerProviderResponseDataSource(ctx, readResponse.FileBasedKeyManagerProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CustomKeyManagerProviderResponse != nil {
		readCustomKeyManagerProviderResponseDataSource(ctx, readResponse.CustomKeyManagerProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.Pkcs11KeyManagerProviderResponse != nil {
		readPkcs11KeyManagerProviderResponseDataSource(ctx, readResponse.Pkcs11KeyManagerProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyKeyManagerProviderResponse != nil {
		readThirdPartyKeyManagerProviderResponseDataSource(ctx, readResponse.ThirdPartyKeyManagerProviderResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
