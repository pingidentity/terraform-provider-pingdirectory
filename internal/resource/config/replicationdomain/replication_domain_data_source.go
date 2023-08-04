package replicationdomain

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
	_ datasource.DataSource              = &replicationDomainDataSource{}
	_ datasource.DataSourceWithConfigure = &replicationDomainDataSource{}
)

// Create a Replication Domain data source
func NewReplicationDomainDataSource() datasource.DataSource {
	return &replicationDomainDataSource{}
}

// replicationDomainDataSource is the datasource implementation.
type replicationDomainDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *replicationDomainDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication_domain"
}

// Configure adds the provider configured client to the data source.
func (r *replicationDomainDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type replicationDomainDataSourceModel struct {
	Id                                        types.String `tfsdk:"id"`
	Name                                      types.String `tfsdk:"name"`
	Type                                      types.String `tfsdk:"type"`
	SynchronizationProviderName               types.String `tfsdk:"synchronization_provider_name"`
	ServerID                                  types.Int64  `tfsdk:"server_id"`
	BaseDN                                    types.String `tfsdk:"base_dn"`
	WindowSize                                types.Int64  `tfsdk:"window_size"`
	HeartbeatInterval                         types.String `tfsdk:"heartbeat_interval"`
	SyncHistPurgeDelay                        types.String `tfsdk:"sync_hist_purge_delay"`
	Restricted                                types.Bool   `tfsdk:"restricted"`
	OnReplayFailureWaitForDependentOpsTimeout types.String `tfsdk:"on_replay_failure_wait_for_dependent_ops_timeout"`
	DependentOpsReplayFailureWaitTime         types.String `tfsdk:"dependent_ops_replay_failure_wait_time"`
}

// GetSchema defines the schema for the datasource.
func (r *replicationDomainDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Replication Domain.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Replication Domain resource. Options are ['replication-domain']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"synchronization_provider_name": schema.StringAttribute{
				Description: "Name of the parent Synchronization Provider",
				Required:    true,
			},
			"server_id": schema.Int64Attribute{
				Description: "Specifies a unique identifier for the Directory Server within the Replication Domain.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"base_dn": schema.StringAttribute{
				Description: "Specifies the base DN of the replicated data.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"window_size": schema.Int64Attribute{
				Description: "Specifies the maximum number of replication updates the Directory Server can have outstanding from the Replication Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"heartbeat_interval": schema.StringAttribute{
				Description: "Specifies the heartbeat interval that the Directory Server will use when communicating with Replication Servers.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"sync_hist_purge_delay": schema.StringAttribute{
				Description: "The time in seconds after which historical information used in replication conflict resolution is purged. The information is removed from entries when they are modified after the purge delay has elapsed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"restricted": schema.BoolAttribute{
				Description: "When set to true, changes are only replicated with server instances that belong to the same replication set.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"on_replay_failure_wait_for_dependent_ops_timeout": schema.StringAttribute{
				Description: "Defines the maximum time to retry a failed operation. An operation will be retried only if it appears that the failure might be dependent on an earlier operation from a different server that hasn't replicated yet. The frequency of the retry is determined by the dependent-ops-replay-failure-wait-time property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"dependent_ops_replay_failure_wait_time": schema.StringAttribute{
				Description: "Defines how long to wait before retrying certain operations, specifically operations that might have failed because they depend on an operation from a different server that has not yet replicated to this instance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a ReplicationDomainResponse object into the model struct
func readReplicationDomainResponseDataSource(ctx context.Context, r *client.ReplicationDomainResponse, state *replicationDomainDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("replication-domain")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ServerID = types.Int64Value(r.ServerID)
	state.BaseDN = types.StringValue(r.BaseDN)
	state.WindowSize = internaltypes.Int64TypeOrNil(r.WindowSize)
	state.HeartbeatInterval = internaltypes.StringTypeOrNil(r.HeartbeatInterval, false)
	state.SyncHistPurgeDelay = internaltypes.StringTypeOrNil(r.SyncHistPurgeDelay, false)
	state.Restricted = internaltypes.BoolTypeOrNil(r.Restricted)
	state.OnReplayFailureWaitForDependentOpsTimeout = internaltypes.StringTypeOrNil(r.OnReplayFailureWaitForDependentOpsTimeout, false)
	state.DependentOpsReplayFailureWaitTime = internaltypes.StringTypeOrNil(r.DependentOpsReplayFailureWaitTime, false)
}

// Read resource information
func (r *replicationDomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state replicationDomainDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ReplicationDomainApi.GetReplicationDomain(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.SynchronizationProviderName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Replication Domain", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readReplicationDomainResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
