// Copyright © 2025 Ping Identity Corporation

package keypair

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
	_ datasource.DataSource              = &keyPairDataSource{}
	_ datasource.DataSourceWithConfigure = &keyPairDataSource{}
)

// Create a Key Pair data source
func NewKeyPairDataSource() datasource.DataSource {
	return &keyPairDataSource{}
}

// keyPairDataSource is the datasource implementation.
type keyPairDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *keyPairDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_pair"
}

// Configure adds the provider configured client to the data source.
func (r *keyPairDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type keyPairDataSourceModel struct {
	Id                            types.String `tfsdk:"id"`
	Name                          types.String `tfsdk:"name"`
	Type                          types.String `tfsdk:"type"`
	KeyAlgorithm                  types.String `tfsdk:"key_algorithm"`
	SelfSignedCertificateValidity types.String `tfsdk:"self_signed_certificate_validity"`
	SubjectDN                     types.String `tfsdk:"subject_dn"`
	CertificateChain              types.String `tfsdk:"certificate_chain"`
	PrivateKey                    types.String `tfsdk:"private_key"`
}

// GetSchema defines the schema for the datasource.
func (r *keyPairDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Key Pair.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Key Pair resource. Options are ['key-pair']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_algorithm": schema.StringAttribute{
				Description: "The algorithm name and the length in bits of the key, e.g. RSA_2048.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"self_signed_certificate_validity": schema.StringAttribute{
				Description: "The validity period for a self-signed certificate. If not specified, the self-signed certificate will be valid for approximately 20 years. This is not used when importing an existing key-pair. The system will not automatically rotate expired certificates. It is up to the administrator to do that when that happens.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"subject_dn": schema.StringAttribute{
				Description: "The DN that should be used as the subject for the self-signed certificate and certificate signing request. This is not used when importing an existing key-pair.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"certificate_chain": schema.StringAttribute{
				Description: "The PEM-encoded X.509 certificate chain.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "The base64-encoded private key that is encrypted using the preferred encryption settings definition.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a KeyPairResponse object into the model struct
func readKeyPairResponseDataSource(ctx context.Context, r *client.KeyPairResponse, state *keyPairDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("key-pair")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.KeyAlgorithm = types.StringValue(r.KeyAlgorithm.String())
	state.SelfSignedCertificateValidity = internaltypes.StringTypeOrNil(r.SelfSignedCertificateValidity, false)
	state.SubjectDN = internaltypes.StringTypeOrNil(r.SubjectDN, false)
	state.CertificateChain = internaltypes.StringTypeOrNil(r.CertificateChain, false)
}

// Read resource information
func (r *keyPairDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state keyPairDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.KeyPairAPI.GetKeyPair(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Key Pair", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readKeyPairResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
