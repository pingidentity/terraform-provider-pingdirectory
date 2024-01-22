package trustedcertificate

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
	_ datasource.DataSource              = &trustedCertificateDataSource{}
	_ datasource.DataSourceWithConfigure = &trustedCertificateDataSource{}
)

// Create a Trusted Certificate data source
func NewTrustedCertificateDataSource() datasource.DataSource {
	return &trustedCertificateDataSource{}
}

// trustedCertificateDataSource is the datasource implementation.
type trustedCertificateDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *trustedCertificateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trusted_certificate"
}

// Configure adds the provider configured client to the data source.
func (r *trustedCertificateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type trustedCertificateDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Certificate types.String `tfsdk:"certificate"`
}

// GetSchema defines the schema for the datasource.
func (r *trustedCertificateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Trusted Certificate.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Trusted Certificate resource. Options are ['trusted-certificate']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"certificate": schema.StringAttribute{
				Description: "The PEM-encoded X.509v3 certificate.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a TrustedCertificateResponse object into the model struct
func readTrustedCertificateResponseDataSource(ctx context.Context, r *client.TrustedCertificateResponse, state *trustedCertificateDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("trusted-certificate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Certificate = types.StringValue(r.Certificate)
}

// Read resource information
func (r *trustedCertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state trustedCertificateDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.TrustedCertificateAPI.GetTrustedCertificate(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Trusted Certificate", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readTrustedCertificateResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
