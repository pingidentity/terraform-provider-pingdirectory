// Copyright © 2025 Ping Identity Corporation

package groupimplementation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &groupImplementationDataSource{}
	_ datasource.DataSourceWithConfigure = &groupImplementationDataSource{}
)

// Create a Group Implementation data source
func NewGroupImplementationDataSource() datasource.DataSource {
	return &groupImplementationDataSource{}
}

// groupImplementationDataSource is the datasource implementation.
type groupImplementationDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *groupImplementationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_implementation"
}

// Configure adds the provider configured client to the data source.
func (r *groupImplementationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type groupImplementationDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *groupImplementationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Group Implementation.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Group Implementation resource. Options are ['static', 'inverted-static', 'virtual-static', 'dynamic']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Group Implementation",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Group Implementation is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a StaticGroupImplementationResponse object into the model struct
func readStaticGroupImplementationResponseDataSource(ctx context.Context, r *client.StaticGroupImplementationResponse, state *groupImplementationDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("static")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a InvertedStaticGroupImplementationResponse object into the model struct
func readInvertedStaticGroupImplementationResponseDataSource(ctx context.Context, r *client.InvertedStaticGroupImplementationResponse, state *groupImplementationDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("inverted-static")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a VirtualStaticGroupImplementationResponse object into the model struct
func readVirtualStaticGroupImplementationResponseDataSource(ctx context.Context, r *client.VirtualStaticGroupImplementationResponse, state *groupImplementationDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("virtual-static")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a DynamicGroupImplementationResponse object into the model struct
func readDynamicGroupImplementationResponseDataSource(ctx context.Context, r *client.DynamicGroupImplementationResponse, state *groupImplementationDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("dynamic")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *groupImplementationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state groupImplementationDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.GroupImplementationAPI.GetGroupImplementation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Group Implementation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.StaticGroupImplementationResponse != nil {
		readStaticGroupImplementationResponseDataSource(ctx, readResponse.StaticGroupImplementationResponse, &state, &resp.Diagnostics)
	}
	if readResponse.InvertedStaticGroupImplementationResponse != nil {
		readInvertedStaticGroupImplementationResponseDataSource(ctx, readResponse.InvertedStaticGroupImplementationResponse, &state, &resp.Diagnostics)
	}
	if readResponse.VirtualStaticGroupImplementationResponse != nil {
		readVirtualStaticGroupImplementationResponseDataSource(ctx, readResponse.VirtualStaticGroupImplementationResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DynamicGroupImplementationResponse != nil {
		readDynamicGroupImplementationResponseDataSource(ctx, readResponse.DynamicGroupImplementationResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
