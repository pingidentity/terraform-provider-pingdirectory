package obscuredvalue

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
	_ datasource.DataSource              = &obscuredValueDataSource{}
	_ datasource.DataSourceWithConfigure = &obscuredValueDataSource{}
)

// Create a Obscured Value data source
func NewObscuredValueDataSource() datasource.DataSource {
	return &obscuredValueDataSource{}
}

// obscuredValueDataSource is the datasource implementation.
type obscuredValueDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *obscuredValueDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_obscured_value"
}

// Configure adds the provider configured client to the data source.
func (r *obscuredValueDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type obscuredValueDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	Description   types.String `tfsdk:"description"`
	ObscuredValue types.String `tfsdk:"obscured_value"`
}

// GetSchema defines the schema for the datasource.
func (r *obscuredValueDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Obscured Value.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Obscured Value",
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
		},
	}
}

// Read a ObscuredValueResponse object into the model struct
func readObscuredValueResponseDataSource(ctx context.Context, r *client.ObscuredValueResponse, state *obscuredValueDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *obscuredValueDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state obscuredValueDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ObscuredValueApi.GetObscuredValue(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Obscured Value", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readObscuredValueResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
