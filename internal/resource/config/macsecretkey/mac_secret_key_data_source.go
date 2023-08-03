package macsecretkey

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
	_ datasource.DataSource              = &macSecretKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &macSecretKeyDataSource{}
)

// Create a Mac Secret Key data source
func NewMacSecretKeyDataSource() datasource.DataSource {
	return &macSecretKeyDataSource{}
}

// macSecretKeyDataSource is the datasource implementation.
type macSecretKeyDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *macSecretKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mac_secret_key"
}

// Configure adds the provider configured client to the data source.
func (r *macSecretKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type macSecretKeyDataSourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	ServerInstanceName types.String `tfsdk:"server_instance_name"`
	MacAlgorithmName   types.String `tfsdk:"mac_algorithm_name"`
	KeyID              types.String `tfsdk:"key_id"`
	IsCompromised      types.Bool   `tfsdk:"is_compromised"`
	SymmetricKey       types.Set    `tfsdk:"symmetric_key"`
	KeyLengthBits      types.Int64  `tfsdk:"key_length_bits"`
}

// GetSchema defines the schema for the datasource.
func (r *macSecretKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Mac Secret Key.",
		Attributes: map[string]schema.Attribute{
			"server_instance_name": schema.StringAttribute{
				Description: "Name of the parent Server Instance",
				Required:    true,
			},
			"mac_algorithm_name": schema.StringAttribute{
				Description: "The algorithm name used to generate this MAC key, e.g. HmacMD5, HmacSHA1, HMacSHA256, etc.",
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

// Read a MacSecretKeyResponse object into the model struct
func readMacSecretKeyResponseDataSource(ctx context.Context, r *client.MacSecretKeyResponse, state *macSecretKeyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MacAlgorithmName = internaltypes.StringTypeOrNil(r.MacAlgorithmName, false)
	state.KeyID = types.StringValue(r.KeyID)
	state.IsCompromised = internaltypes.BoolTypeOrNil(r.IsCompromised)
	state.SymmetricKey = internaltypes.GetStringSet(r.SymmetricKey)
	state.KeyLengthBits = types.Int64Value(r.KeyLengthBits)
}

// Read resource information
func (r *macSecretKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state macSecretKeyDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.MacSecretKeyApi.GetMacSecretKey(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.ServerInstanceName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Mac Secret Key", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readMacSecretKeyResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
