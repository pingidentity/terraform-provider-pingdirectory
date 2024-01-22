package ciphersecretkey

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
	_ datasource.DataSource              = &cipherSecretKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &cipherSecretKeyDataSource{}
)

// Create a Cipher Secret Key data source
func NewCipherSecretKeyDataSource() datasource.DataSource {
	return &cipherSecretKeyDataSource{}
}

// cipherSecretKeyDataSource is the datasource implementation.
type cipherSecretKeyDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *cipherSecretKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cipher_secret_key"
}

// Configure adds the provider configured client to the data source.
func (r *cipherSecretKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type cipherSecretKeyDataSourceModel struct {
	Id                             types.String `tfsdk:"id"`
	Name                           types.String `tfsdk:"name"`
	Type                           types.String `tfsdk:"type"`
	ServerInstanceName             types.String `tfsdk:"server_instance_name"`
	CipherTransformationName       types.String `tfsdk:"cipher_transformation_name"`
	InitializationVectorLengthBits types.Int64  `tfsdk:"initialization_vector_length_bits"`
	KeyID                          types.String `tfsdk:"key_id"`
	IsCompromised                  types.Bool   `tfsdk:"is_compromised"`
	SymmetricKey                   types.Set    `tfsdk:"symmetric_key"`
	KeyLengthBits                  types.Int64  `tfsdk:"key_length_bits"`
}

// GetSchema defines the schema for the datasource.
func (r *cipherSecretKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Cipher Secret Key.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Cipher Secret Key resource. Options are ['cipher-secret-key']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_instance_name": schema.StringAttribute{
				Description: "Name of the parent Server Instance",
				Required:    true,
			},
			"cipher_transformation_name": schema.StringAttribute{
				Description: "The algorithm name used to produce this cipher, e.g. AES/CBC/PKCS5Padding.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"initialization_vector_length_bits": schema.Int64Attribute{
				Description: "The initialization vector length of the cipher in bits.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_id": schema.StringAttribute{
				Description: "The unique system-generated identifier for the Secret Key.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"is_compromised": schema.BoolAttribute{
				Description: "If the key is compromised, an administrator may set this flag to immediately trigger the creation of a new secret key. After the new key is generated, the value of this property will be reset to false.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"symmetric_key": schema.SetAttribute{
				Description: "The symmetric key that is used for both encryption of plain text and decryption of cipher text. This stores the secret key for each server instance encrypted with that server's inter-server certificate.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"key_length_bits": schema.Int64Attribute{
				Description: "The length of the key in bits.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a CipherSecretKeyResponse object into the model struct
func readCipherSecretKeyResponseDataSource(ctx context.Context, r *client.CipherSecretKeyResponse, state *cipherSecretKeyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("cipher-secret-key")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CipherTransformationName = internaltypes.StringTypeOrNil(r.CipherTransformationName, false)
	state.InitializationVectorLengthBits = internaltypes.Int64TypeOrNil(r.InitializationVectorLengthBits)
	state.KeyID = types.StringValue(r.KeyID)
	state.IsCompromised = internaltypes.BoolTypeOrNil(r.IsCompromised)
	state.SymmetricKey = internaltypes.GetStringSet(r.SymmetricKey)
	state.KeyLengthBits = types.Int64Value(r.KeyLengthBits)
}

// Read resource information
func (r *cipherSecretKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state cipherSecretKeyDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CipherSecretKeyAPI.GetCipherSecretKey(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.ServerInstanceName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Cipher Secret Key", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCipherSecretKeyResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
