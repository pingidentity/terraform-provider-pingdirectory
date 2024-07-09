package replicationserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &replicationServerResource{}
	_ resource.ResourceWithConfigure   = &replicationServerResource{}
	_ resource.ResourceWithImportState = &replicationServerResource{}
)

// Create a Replication Server resource
func NewReplicationServerResource() resource.Resource {
	return &replicationServerResource{}
}

// replicationServerResource is the resource implementation.
type replicationServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *replicationServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_replication_server"
}

// Configure adds the provider configured client to the resource.
func (r *replicationServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type replicationServerResourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	Notifications                                types.Set    `tfsdk:"notifications"`
	RequiredActions                              types.Set    `tfsdk:"required_actions"`
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

// GetSchema defines the schema for the resource.
func (r *replicationServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Replication Server.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Replication Server resource. Options are ['replication-server']",
				Optional:    false,
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"replication-server"}...),
				},
			},
			"synchronization_provider_name": schema.StringAttribute{
				Description: "Name of the parent Synchronization Provider",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"replication_server_id": schema.Int64Attribute{
				Description: "Specifies a unique identifier for the Replication Server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"replication_db_directory": schema.StringAttribute{
				Description: "The path where the Replication Server stores all persistent information.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"je_property": schema.SetAttribute{
				Description: "Specifies the database and environment properties for the Berkeley DB Java Edition database for the replication changelog.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"replication_purge_delay": schema.StringAttribute{
				Description: "Changes are guaranteed to be maintained in the changelog database for at least this duration. Setting target-database-size can allow additional changes to be maintained up to the configured size on disk.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"target_database_size": schema.StringAttribute{
				Description: "The replication changelog database is allowed to grow up to this size even if changes are older than the configured replication-purge-delay.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"replication_port": schema.Int64Attribute{
				Description: "The port on which this Replication Server waits for connections from other Replication Servers or Directory Server instances.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"listen_on_all_addresses": schema.BoolAttribute{
				Description: "Indicates whether the Replication Server should listen on all addresses for this host. If set to FALSE, then the Replication Server will listen only to the address resolved from the hostname provided.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"compression_criteria": schema.StringAttribute{
				Description: "Specifies when the replication traffic should be compressed.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"always", "remote", "never"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"heartbeat_interval": schema.StringAttribute{
				Description: "Specifies the heartbeat interval that the Directory Server will use when communicating with Replication Servers.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"remote_monitor_update_interval": schema.StringAttribute{
				Description: "Specifies the duration that topology monitor data will be cached before it is requested again from a remote server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"restricted_domain": schema.SetAttribute{
				Description: "Specifies the base DN of domains that are only replicated between server instances that belong to the same replication set.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"gateway_priority": schema.Int64Attribute{
				Description: "Specifies the gateway priority of the Replication Server in the current location.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"missing_changes_alert_threshold_percent": schema.Int64Attribute{
				Description: "Supported in PingDirectory product version 9.3.0.0+. Specifies the missing changes alert threshold as a percentage of the total pending changes. For instance, a value of 80 indicates that the replica is 80% of the way to losing changes.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"missing_changes_policy": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. Determines how the server responds when replication detects that some changes might have been missed. Each missing changes policy is a set of missing changes actions to take for a set of missing changes types. The value configured here acts as a default for all replication domains on this replication server.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"maximum-integrity", "favor-integrity", "favor-availability", "maximum-availability"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"include_all_remote_servers_state_in_monitor_message": schema.BoolAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. Indicates monitor messages should include information about remote servers.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *replicationServerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingDirectory10000)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model replicationServerResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsNonEmptyString(model.MissingChangesPolicy) {
		resp.Diagnostics.AddError("Attribute 'missing_changes_policy' not supported by PingDirectory version "+r.providerConfig.ProductVersion, "")
	}
	if internaltypes.IsDefined(model.IncludeAllRemoteServersStateInMonitorMessage) {
		resp.Diagnostics.AddError("Attribute 'include_all_remote_servers_state_in_monitor_message' not supported by PingDirectory version "+r.providerConfig.ProductVersion, "")
	}
	compare, err = version.Compare(r.providerConfig.ProductVersion, version.PingDirectory9300)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	if internaltypes.IsDefined(model.MissingChangesAlertThresholdPercent) {
		resp.Diagnostics.AddError("Attribute 'missing_changes_alert_threshold_percent' not supported by PingDirectory version "+r.providerConfig.ProductVersion, "")
	}
}

// Read a ReplicationServerResponse object into the model struct
func readReplicationServerResponse(ctx context.Context, r *client.ReplicationServerResponse, state *replicationServerResourceModel, expectedValues *replicationServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("replication-server")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.ReplicationServerID = types.Int64Value(r.ReplicationServerID)
	state.ReplicationDBDirectory = types.StringValue(r.ReplicationDBDirectory)
	state.JeProperty = internaltypes.GetStringSet(r.JeProperty)
	state.ReplicationPurgeDelay = internaltypes.StringTypeOrNil(r.ReplicationPurgeDelay, true)
	config.CheckMismatchedPDFormattedAttributes("replication_purge_delay",
		expectedValues.ReplicationPurgeDelay, state.ReplicationPurgeDelay, diagnostics)
	state.TargetDatabaseSize = internaltypes.StringTypeOrNil(r.TargetDatabaseSize, true)
	config.CheckMismatchedPDFormattedAttributes("target_database_size",
		expectedValues.TargetDatabaseSize, state.TargetDatabaseSize, diagnostics)
	state.ReplicationPort = types.Int64Value(r.ReplicationPort)
	state.ListenOnAllAddresses = internaltypes.BoolTypeOrNil(r.ListenOnAllAddresses)
	state.CompressionCriteria = internaltypes.StringTypeOrNil(
		client.StringPointerEnumreplicationServerCompressionCriteriaProp(r.CompressionCriteria), true)
	state.HeartbeatInterval = internaltypes.StringTypeOrNil(r.HeartbeatInterval, true)
	config.CheckMismatchedPDFormattedAttributes("heartbeat_interval",
		expectedValues.HeartbeatInterval, state.HeartbeatInterval, diagnostics)
	state.RemoteMonitorUpdateInterval = internaltypes.StringTypeOrNil(r.RemoteMonitorUpdateInterval, true)
	config.CheckMismatchedPDFormattedAttributes("remote_monitor_update_interval",
		expectedValues.RemoteMonitorUpdateInterval, state.RemoteMonitorUpdateInterval, diagnostics)
	state.RestrictedDomain = internaltypes.GetStringSet(r.RestrictedDomain)
	state.GatewayPriority = types.Int64Value(r.GatewayPriority)
	state.MissingChangesAlertThresholdPercent = internaltypes.Int64TypeOrNil(r.MissingChangesAlertThresholdPercent)
	state.MissingChangesPolicy = internaltypes.StringTypeOrNil(
		client.StringPointerEnumreplicationServerMissingChangesPolicyProp(r.MissingChangesPolicy), true)
	state.IncludeAllRemoteServersStateInMonitorMessage = internaltypes.BoolTypeOrNil(r.IncludeAllRemoteServersStateInMonitorMessage)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *replicationServerResourceModel) setStateValuesNotReturnedByAPI(expectedValues *replicationServerResourceModel) {
	if !expectedValues.SynchronizationProviderName.IsUnknown() {
		state.SynchronizationProviderName = expectedValues.SynchronizationProviderName
	}
}

// Create any update operations necessary to make the state match the plan
func createReplicationServerOperations(plan replicationServerResourceModel, state replicationServerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddInt64OperationIfNecessary(&ops, plan.ReplicationServerID, state.ReplicationServerID, "replication-server-id")
	operations.AddStringOperationIfNecessary(&ops, plan.ReplicationDBDirectory, state.ReplicationDBDirectory, "replication-db-directory")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.JeProperty, state.JeProperty, "je-property")
	operations.AddStringOperationIfNecessary(&ops, plan.ReplicationPurgeDelay, state.ReplicationPurgeDelay, "replication-purge-delay")
	operations.AddStringOperationIfNecessary(&ops, plan.TargetDatabaseSize, state.TargetDatabaseSize, "target-database-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.ReplicationPort, state.ReplicationPort, "replication-port")
	operations.AddBoolOperationIfNecessary(&ops, plan.ListenOnAllAddresses, state.ListenOnAllAddresses, "listen-on-all-addresses")
	operations.AddStringOperationIfNecessary(&ops, plan.CompressionCriteria, state.CompressionCriteria, "compression-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.HeartbeatInterval, state.HeartbeatInterval, "heartbeat-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.RemoteMonitorUpdateInterval, state.RemoteMonitorUpdateInterval, "remote-monitor-update-interval")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RestrictedDomain, state.RestrictedDomain, "restricted-domain")
	operations.AddInt64OperationIfNecessary(&ops, plan.GatewayPriority, state.GatewayPriority, "gateway-priority")
	operations.AddInt64OperationIfNecessary(&ops, plan.MissingChangesAlertThresholdPercent, state.MissingChangesAlertThresholdPercent, "missing-changes-alert-threshold-percent")
	operations.AddStringOperationIfNecessary(&ops, plan.MissingChangesPolicy, state.MissingChangesPolicy, "missing-changes-policy")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeAllRemoteServersStateInMonitorMessage, state.IncludeAllRemoteServersStateInMonitorMessage, "include-all-remote-servers-state-in-monitor-message")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *replicationServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan replicationServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ReplicationServerAPI.GetReplicationServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.SynchronizationProviderName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Replication Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state replicationServerResourceModel
	readReplicationServerResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ReplicationServerAPI.UpdateReplicationServer(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.SynchronizationProviderName.ValueString())
	ops := createReplicationServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ReplicationServerAPI.UpdateReplicationServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Replication Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readReplicationServerResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *replicationServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state replicationServerResourceModel
	diags := req.State.Get(ctx, &state)
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
	readReplicationServerResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *replicationServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan replicationServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state replicationServerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ReplicationServerAPI.UpdateReplicationServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.SynchronizationProviderName.ValueString())

	// Determine what update operations are necessary
	ops := createReplicationServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ReplicationServerAPI.UpdateReplicationServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Replication Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readReplicationServerResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *replicationServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *replicationServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve parent name and save to synchronization_provider_name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("synchronization_provider_name"), req, resp)
}
