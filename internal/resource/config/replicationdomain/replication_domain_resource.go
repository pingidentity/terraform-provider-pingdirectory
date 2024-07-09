package replicationdomain

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	_ resource.Resource                = &replicationDomainResource{}
	_ resource.ResourceWithConfigure   = &replicationDomainResource{}
	_ resource.ResourceWithImportState = &replicationDomainResource{}
)

// Create a Replication Domain resource
func NewReplicationDomainResource() resource.Resource {
	return &replicationDomainResource{}
}

// replicationDomainResource is the resource implementation.
type replicationDomainResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *replicationDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_replication_domain"
}

// Configure adds the provider configured client to the resource.
func (r *replicationDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type replicationDomainResourceModel struct {
	Id                                        types.String `tfsdk:"id"`
	Name                                      types.String `tfsdk:"name"`
	Notifications                             types.Set    `tfsdk:"notifications"`
	RequiredActions                           types.Set    `tfsdk:"required_actions"`
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
	MissingChangesPolicy                      types.String `tfsdk:"missing_changes_policy"`
}

// GetSchema defines the schema for the resource.
func (r *replicationDomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Replication Domain.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Replication Domain resource. Options are ['replication-domain']",
				Optional:    false,
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"replication-domain"}...),
				},
			},
			"synchronization_provider_name": schema.StringAttribute{
				Description: "Name of the parent Synchronization Provider",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server_id": schema.Int64Attribute{
				Description: "Specifies a unique identifier for the Directory Server within the Replication Domain.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"base_dn": schema.StringAttribute{
				Description: "Specifies the base DN of the replicated data.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"window_size": schema.Int64Attribute{
				Description: "Specifies the maximum number of replication updates the Directory Server can have outstanding from the Replication Server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
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
			"sync_hist_purge_delay": schema.StringAttribute{
				Description: "The time in seconds after which historical information used in replication conflict resolution is purged. The information is removed from entries when they are modified after the purge delay has elapsed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"restricted": schema.BoolAttribute{
				Description: "When set to true, changes are only replicated with server instances that belong to the same replication set.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"on_replay_failure_wait_for_dependent_ops_timeout": schema.StringAttribute{
				Description: "Defines the maximum time to retry a failed operation. An operation will be retried only if it appears that the failure might be dependent on an earlier operation from a different server that hasn't replicated yet. The frequency of the retry is determined by the dependent-ops-replay-failure-wait-time property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dependent_ops_replay_failure_wait_time": schema.StringAttribute{
				Description: "Defines how long to wait before retrying certain operations, specifically operations that might have failed because they depend on an operation from a different server that has not yet replicated to this instance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"missing_changes_policy": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. Determines how the server responds when replication detects that some changes might have been missed. Each missing changes policy is a set of missing changes actions to take for a set of missing changes types. The value configured here only applies to this particular replication domain.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"maximum-integrity", "favor-integrity", "favor-availability", "maximum-availability", "use-server-default"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *replicationDomainResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingDirectory10000)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model replicationDomainResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsNonEmptyString(model.MissingChangesPolicy) {
		resp.Diagnostics.AddError("Attribute 'missing_changes_policy' not supported by PingDirectory version "+r.providerConfig.ProductVersion, "")
	}
}

// Read a ReplicationDomainResponse object into the model struct
func readReplicationDomainResponse(ctx context.Context, r *client.ReplicationDomainResponse, state *replicationDomainResourceModel, expectedValues *replicationDomainResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("replication-domain")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ServerID = types.Int64Value(r.ServerID)
	state.BaseDN = types.StringValue(r.BaseDN)
	state.WindowSize = internaltypes.Int64TypeOrNil(r.WindowSize)
	state.HeartbeatInterval = internaltypes.StringTypeOrNil(r.HeartbeatInterval, true)
	config.CheckMismatchedPDFormattedAttributes("heartbeat_interval",
		expectedValues.HeartbeatInterval, state.HeartbeatInterval, diagnostics)
	state.SyncHistPurgeDelay = internaltypes.StringTypeOrNil(r.SyncHistPurgeDelay, true)
	config.CheckMismatchedPDFormattedAttributes("sync_hist_purge_delay",
		expectedValues.SyncHistPurgeDelay, state.SyncHistPurgeDelay, diagnostics)
	state.Restricted = internaltypes.BoolTypeOrNil(r.Restricted)
	state.OnReplayFailureWaitForDependentOpsTimeout = internaltypes.StringTypeOrNil(r.OnReplayFailureWaitForDependentOpsTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("on_replay_failure_wait_for_dependent_ops_timeout",
		expectedValues.OnReplayFailureWaitForDependentOpsTimeout, state.OnReplayFailureWaitForDependentOpsTimeout, diagnostics)
	state.DependentOpsReplayFailureWaitTime = internaltypes.StringTypeOrNil(r.DependentOpsReplayFailureWaitTime, true)
	config.CheckMismatchedPDFormattedAttributes("dependent_ops_replay_failure_wait_time",
		expectedValues.DependentOpsReplayFailureWaitTime, state.DependentOpsReplayFailureWaitTime, diagnostics)
	state.MissingChangesPolicy = internaltypes.StringTypeOrNil(
		client.StringPointerEnumreplicationDomainMissingChangesPolicyProp(r.MissingChangesPolicy), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *replicationDomainResourceModel) setStateValuesNotReturnedByAPI(expectedValues *replicationDomainResourceModel) {
	if !expectedValues.SynchronizationProviderName.IsUnknown() {
		state.SynchronizationProviderName = expectedValues.SynchronizationProviderName
	}
}

// Create any update operations necessary to make the state match the plan
func createReplicationDomainOperations(plan replicationDomainResourceModel, state replicationDomainResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerID, state.ServerID, "server-id")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddInt64OperationIfNecessary(&ops, plan.WindowSize, state.WindowSize, "window-size")
	operations.AddStringOperationIfNecessary(&ops, plan.HeartbeatInterval, state.HeartbeatInterval, "heartbeat-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.SyncHistPurgeDelay, state.SyncHistPurgeDelay, "sync-hist-purge-delay")
	operations.AddBoolOperationIfNecessary(&ops, plan.Restricted, state.Restricted, "restricted")
	operations.AddStringOperationIfNecessary(&ops, plan.OnReplayFailureWaitForDependentOpsTimeout, state.OnReplayFailureWaitForDependentOpsTimeout, "on-replay-failure-wait-for-dependent-ops-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.DependentOpsReplayFailureWaitTime, state.DependentOpsReplayFailureWaitTime, "dependent-ops-replay-failure-wait-time")
	operations.AddStringOperationIfNecessary(&ops, plan.MissingChangesPolicy, state.MissingChangesPolicy, "missing-changes-policy")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *replicationDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan replicationDomainResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ReplicationDomainAPI.GetReplicationDomain(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.SynchronizationProviderName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Replication Domain", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state replicationDomainResourceModel
	readReplicationDomainResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ReplicationDomainAPI.UpdateReplicationDomain(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.SynchronizationProviderName.ValueString())
	ops := createReplicationDomainOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ReplicationDomainAPI.UpdateReplicationDomainExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Replication Domain", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readReplicationDomainResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *replicationDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state replicationDomainResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ReplicationDomainAPI.GetReplicationDomain(
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
	readReplicationDomainResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *replicationDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan replicationDomainResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state replicationDomainResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ReplicationDomainAPI.UpdateReplicationDomain(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.SynchronizationProviderName.ValueString())

	// Determine what update operations are necessary
	ops := createReplicationDomainOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ReplicationDomainAPI.UpdateReplicationDomainExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Replication Domain", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readReplicationDomainResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *replicationDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *replicationDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [synchronization-provider-name]/[replication-domain-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("synchronization_provider_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), split[1])...)
}
