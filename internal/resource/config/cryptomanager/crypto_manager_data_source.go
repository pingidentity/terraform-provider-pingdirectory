package cryptomanager

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
	_ datasource.DataSource              = &cryptoManagerDataSource{}
	_ datasource.DataSourceWithConfigure = &cryptoManagerDataSource{}
)

// Create a Crypto Manager data source
func NewCryptoManagerDataSource() datasource.DataSource {
	return &cryptoManagerDataSource{}
}

// cryptoManagerDataSource is the datasource implementation.
type cryptoManagerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *cryptoManagerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_crypto_manager"
}

// Configure adds the provider configured client to the data source.
func (r *cryptoManagerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type cryptoManagerDataSourceModel struct {
	Id                               types.String `tfsdk:"id"`
	Type                             types.String `tfsdk:"type"`
	DigestAlgorithm                  types.String `tfsdk:"digest_algorithm"`
	MacAlgorithm                     types.String `tfsdk:"mac_algorithm"`
	MacKeyLength                     types.Int64  `tfsdk:"mac_key_length"`
	SigningEncryptionSettingsID      types.String `tfsdk:"signing_encryption_settings_id"`
	CipherTransformation             types.String `tfsdk:"cipher_transformation"`
	CipherKeyLength                  types.Int64  `tfsdk:"cipher_key_length"`
	KeyWrappingTransformation        types.String `tfsdk:"key_wrapping_transformation"`
	SslProtocol                      types.Set    `tfsdk:"ssl_protocol"`
	SslCipherSuite                   types.Set    `tfsdk:"ssl_cipher_suite"`
	OutboundSSLProtocol              types.Set    `tfsdk:"outbound_ssl_protocol"`
	OutboundSSLCipherSuite           types.Set    `tfsdk:"outbound_ssl_cipher_suite"`
	EnableSha1CipherSuites           types.Bool   `tfsdk:"enable_sha_1_cipher_suites"`
	EnableRsaKeyExchangeCipherSuites types.Bool   `tfsdk:"enable_rsa_key_exchange_cipher_suites"`
	SslCertNickname                  types.String `tfsdk:"ssl_cert_nickname"`
}

// GetSchema defines the schema for the datasource.
func (r *cryptoManagerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Crypto Manager.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Crypto Manager resource. Options are ['crypto-manager']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"digest_algorithm": schema.StringAttribute{
				Description: "Specifies the preferred message digest algorithm for the Directory Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"mac_algorithm": schema.StringAttribute{
				Description: "Specifies the preferred MAC algorithm for the Directory Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"mac_key_length": schema.Int64Attribute{
				Description: "Specifies the key length in bits for the preferred MAC algorithm.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"signing_encryption_settings_id": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.0+. The ID of the encryption settings definition to use for generating digital signatures. If this is not specified, then the server's preferred encryption settings definition will be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"cipher_transformation": schema.StringAttribute{
				Description: "Specifies the cipher for the Directory Server using the syntax algorithm/mode/padding.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"cipher_key_length": schema.Int64Attribute{
				Description: "Specifies the key length in bits for the preferred cipher.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_wrapping_transformation": schema.StringAttribute{
				Description: "The preferred key wrapping transformation for the Directory Server. This value must be the same for all server instances in a replication topology.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ssl_protocol": schema.SetAttribute{
				Description: "Specifies the names of TLS protocols that are allowed for use in secure communication.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"ssl_cipher_suite": schema.SetAttribute{
				Description: "Specifies the names of the TLS cipher suites that are allowed for use in secure communication.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"outbound_ssl_protocol": schema.SetAttribute{
				Description: "Specifies the names of the TLS protocols that will be enabled for outbound connections initiated by the Directory Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"outbound_ssl_cipher_suite": schema.SetAttribute{
				Description: "Specifies the names of the TLS cipher suites that will be enabled for outbound connections initiated by the Directory Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"enable_sha_1_cipher_suites": schema.BoolAttribute{
				Description: "Indicates whether to enable support for TLS cipher suites that use the SHA-1 digest algorithm. The SHA-1 digest algorithm is no longer considered secure and is not recommended for use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enable_rsa_key_exchange_cipher_suites": schema.BoolAttribute{
				Description: "Indicates whether to enable support for TLS cipher suites that use the RSA key exchange algorithm. Cipher suites that rely on RSA key exchange are not recommended because they do not support forward secrecy, which means that if the private key is compromised, then any communication negotiated using that private key should also be considered compromised.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ssl_cert_nickname": schema.StringAttribute{
				Description: "Specifies the nickname (also called the alias) of the certificate that the Crypto Manager should use when performing SSL communication.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a CryptoManagerResponse object into the model struct
func readCryptoManagerResponseDataSource(ctx context.Context, r *client.CryptoManagerResponse, state *cryptoManagerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("crypto-manager")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.DigestAlgorithm = internaltypes.StringTypeOrNil(r.DigestAlgorithm, false)
	state.MacAlgorithm = internaltypes.StringTypeOrNil(r.MacAlgorithm, false)
	state.MacKeyLength = internaltypes.Int64TypeOrNil(r.MacKeyLength)
	state.SigningEncryptionSettingsID = internaltypes.StringTypeOrNil(r.SigningEncryptionSettingsID, false)
	state.CipherTransformation = internaltypes.StringTypeOrNil(r.CipherTransformation, false)
	state.CipherKeyLength = internaltypes.Int64TypeOrNil(r.CipherKeyLength)
	state.KeyWrappingTransformation = internaltypes.StringTypeOrNil(r.KeyWrappingTransformation, false)
	state.SslProtocol = internaltypes.GetStringSet(r.SslProtocol)
	state.SslCipherSuite = internaltypes.GetStringSet(r.SslCipherSuite)
	state.OutboundSSLProtocol = internaltypes.GetStringSet(r.OutboundSSLProtocol)
	state.OutboundSSLCipherSuite = internaltypes.GetStringSet(r.OutboundSSLCipherSuite)
	state.EnableSha1CipherSuites = internaltypes.BoolTypeOrNil(r.EnableSha1CipherSuites)
	state.EnableRsaKeyExchangeCipherSuites = internaltypes.BoolTypeOrNil(r.EnableRsaKeyExchangeCipherSuites)
	state.SslCertNickname = internaltypes.StringTypeOrNil(r.SslCertNickname, false)
}

// Read resource information
func (r *cryptoManagerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state cryptoManagerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CryptoManagerApi.GetCryptoManager(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Crypto Manager", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCryptoManagerResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
