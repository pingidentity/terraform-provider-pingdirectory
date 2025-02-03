// Copyright © 2025 Ping Identity Corporation

package synchronizationprovider

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
	_ datasource.DataSource              = &synchronizationProviderDataSource{}
	_ datasource.DataSourceWithConfigure = &synchronizationProviderDataSource{}
)

// Create a Synchronization Provider data source
func NewSynchronizationProviderDataSource() datasource.DataSource {
	return &synchronizationProviderDataSource{}
}

// synchronizationProviderDataSource is the datasource implementation.
type synchronizationProviderDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *synchronizationProviderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synchronization_provider"
}

// Configure adds the provider configured client to the data source.
func (r *synchronizationProviderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type synchronizationProviderDataSourceModel struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Type                   types.String `tfsdk:"type"`
	NumUpdateReplayThreads types.Int64  `tfsdk:"num_update_replay_threads"`
	Description            types.String `tfsdk:"description"`
	Enabled                types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *synchronizationProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Synchronization Provider.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Synchronization Provider resource. Options are ['replication', 'custom']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"num_update_replay_threads": schema.Int64Attribute{
				Description: "Specifies the number of update replay threads.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Synchronization Provider",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Synchronization Provider is enabled for use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a ReplicationSynchronizationProviderResponse object into the model struct
func readReplicationSynchronizationProviderResponseDataSource(ctx context.Context, r *client.ReplicationSynchronizationProviderResponse, state *synchronizationProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("replication")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.NumUpdateReplayThreads = internaltypes.Int64TypeOrNil(r.NumUpdateReplayThreads)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a CustomSynchronizationProviderResponse object into the model struct
func readCustomSynchronizationProviderResponseDataSource(ctx context.Context, r *client.CustomSynchronizationProviderResponse, state *synchronizationProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *synchronizationProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state synchronizationProviderDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SynchronizationProviderAPI.GetSynchronizationProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Synchronization Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.ReplicationSynchronizationProviderResponse != nil {
		readReplicationSynchronizationProviderResponseDataSource(ctx, readResponse.ReplicationSynchronizationProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CustomSynchronizationProviderResponse != nil {
		readCustomSynchronizationProviderResponseDataSource(ctx, readResponse.CustomSynchronizationProviderResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
