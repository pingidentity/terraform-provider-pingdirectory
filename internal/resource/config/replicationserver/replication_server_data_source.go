package replicationserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &replicationServerDataSource{}
	_ datasource.DataSourceWithConfigure = &replicationServerDataSource{}
)

// Create a Replication Server data source
func NewReplicationServerDataSource() datasource.DataSource {
	return &replicationServerDataSource{}
}

// replicationServerDataSource is the datasource implementation.
type replicationServerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *replicationServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication_server"
}

// Configure adds the provider configured client to the data source.
func (r *replicationServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type replicationServerDataSourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	Type                                         types.String `tfsdk:"type"`
	SynchronizationProviderName                  types.String `tfsdk:"synchronization_provider_name"`
	ReplicationServerID                          types.Int64  `tfsdk:"replication_server_id"`
	ReplicationDBDirectory                       types.String `tfsdk:"replication_db_directory"`
	JeProperty                                   types.Set    `tfsdk:"je_property"`
	ReplicationPurgeDelay                        types.String `tfsdk:"replication_purge_delay"`
	TargetDatabaseSize                           types.String `tfsdk:"target_database_size"`
	ReplicationPort                              types.Int64  `tfsdk:"replication_port"`
	ListenOnAllAddresses                         types.Bool   `tfsdk:"listen_on_all_addresses"`
	CompressionCriteria                          types.String `tfsdk:"compression_criteria"`
	HeartbeatInterval                            types.String `tfsdk:"heartbeat_interval"`
	RemoteMonitorUpdateInterval                  types.String `tfsdk:"remote_monitor_update_interval"`
	RestrictedDomain                             types.Set    `tfsdk:"restricted_domain"`
	GatewayPriority                              types.Int64  `tfsdk:"gateway_priority"`
	MissingChangesAlertThresholdPercent          types.Int64  `tfsdk:"missing_changes_alert_threshold_percent"`
	MissingChangesPolicy                         types.String `tfsdk:"missing_changes_policy"`
	IncludeAllRemoteServersStateInMonitorMessage types.Bool   `tfsdk:"include_all_remote_servers_state_in_monitor_message"`
}

// GetSchema defines the schema for the datasource.
func (r *replicationServerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Replication Server.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Replication Server resource. Options are ['replication-server']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"synchronization_provider_name": schema.StringAttribute{
				Description: "Name of the parent Synchronization Provider",
				Required:    true,
			},
			"replication_server_id": schema.Int64Attribute{
				Description: "Specifies a unique identifier for the Replication Server.",
				Required:    true,
			},
			"replication_db_directory": schema.StringAttribute{
				Description: "The path where the Replication Server stores all persistent information.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"je_property": schema.SetAttribute{
				Description: "Specifies the database and environment properties for the Berkeley DB Java Edition database for the replication changelog.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"replication_purge_delay": schema.StringAttribute{
				Description: "Changes are guaranteed to be maintained in the changelog database for at least this duration. Setting target-database-size can allow additional changes to be maintained up to the configured size on disk.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"target_database_size": schema.StringAttribute{
				Description: "The replication changelog database is allowed to grow up to this size even if changes are older than the configured replication-purge-delay.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"replication_port": schema.Int64Attribute{
				Description: "The port on which this Replication Server waits for connections from other Replication Servers or Directory Server instances.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"listen_on_all_addresses": schema.BoolAttribute{
				Description: "Indicates whether the Replication Server should listen on all addresses for this host. If set to FALSE, then the Replication Server will listen only to the address resolved from the hostname provided.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"compression_criteria": schema.StringAttribute{
				Description: "Specifies when the replication traffic should be compressed.",
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
			"remote_monitor_update_interval": schema.StringAttribute{
				Description: "Specifies the duration that topology monitor data will be cached before it is requested again from a remote server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"restricted_domain": schema.SetAttribute{
				Description: "Specifies the base DN of domains that are only replicated between server instances that belong to the same replication set.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"gateway_priority": schema.Int64Attribute{
				Description: "Specifies the gateway priority of the Replication Server in the current location.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"missing_changes_alert_threshold_percent": schema.Int64Attribute{
				Description: "Supported in PingDirectory product version 9.3.0.0+. Specifies the missing changes alert threshold as a percentage of the total pending changes. For instance, a value of 80 indicates that the replica is 80% of the way to losing changes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"missing_changes_policy": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. Determines how the server responds when replication detects that some changes might have been missed. Each missing changes policy is a set of missing changes actions to take for a set of missing changes types. The value configured here acts as a default for all replication domains on this replication server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_all_remote_servers_state_in_monitor_message": schema.BoolAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. Indicates monitor messages should include information about remote servers.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a ReplicationServerResponse object into the model struct
func readReplicationServerResponseDataSource(ctx context.Context, r *client.ReplicationServerResponse, state *replicationServerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("replication-server")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.ReplicationServerID = types.Int64Value(r.ReplicationServerID)
	state.ReplicationDBDirectory = types.StringValue(r.ReplicationDBDirectory)
	state.JeProperty = internaltypes.GetStringSet(r.JeProperty)
	state.ReplicationPurgeDelay = internaltypes.StringTypeOrNil(r.ReplicationPurgeDelay, false)
	state.TargetDatabaseSize = internaltypes.StringTypeOrNil(r.TargetDatabaseSize, false)
	state.ReplicationPort = types.Int64Value(r.ReplicationPort)
	state.ListenOnAllAddresses = internaltypes.BoolTypeOrNil(r.ListenOnAllAddresses)
	state.CompressionCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumreplicationServerCompressionCriteriaProp(r.CompressionCriteria), false)
	state.HeartbeatInterval = internaltypes.StringTypeOrNil(r.HeartbeatInterval, false)
	state.RemoteMonitorUpdateInterval = internaltypes.StringTypeOrNil(r.RemoteMonitorUpdateInterval, false)
	state.RestrictedDomain = internaltypes.GetStringSet(r.RestrictedDomain)
	state.GatewayPriority = types.Int64Value(r.GatewayPriority)
	state.MissingChangesAlertThresholdPercent = internaltypes.Int64TypeOrNil(r.MissingChangesAlertThresholdPercent)
	state.MissingChangesPolicy = internaltypes.StringTypeOrNil(
		client.StringPointerEnumreplicationServerMissingChangesPolicyProp(r.MissingChangesPolicy), false)
	state.IncludeAllRemoteServersStateInMonitorMessage = internaltypes.BoolTypeOrNil(r.IncludeAllRemoteServersStateInMonitorMessage)
}

// Read resource information
func (r *replicationServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state replicationServerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ReplicationServerAPI.GetReplicationServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.SynchronizationProviderName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Replication Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readReplicationServerResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
