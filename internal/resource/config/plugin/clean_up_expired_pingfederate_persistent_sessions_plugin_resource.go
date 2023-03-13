package plugin

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &cleanUpExpiredPingfederatePersistentSessionsPluginResource{}
	_ resource.ResourceWithConfigure   = &cleanUpExpiredPingfederatePersistentSessionsPluginResource{}
	_ resource.ResourceWithImportState = &cleanUpExpiredPingfederatePersistentSessionsPluginResource{}
	_ resource.Resource                = &defaultCleanUpExpiredPingfederatePersistentSessionsPluginResource{}
	_ resource.ResourceWithConfigure   = &defaultCleanUpExpiredPingfederatePersistentSessionsPluginResource{}
	_ resource.ResourceWithImportState = &defaultCleanUpExpiredPingfederatePersistentSessionsPluginResource{}
)

// Create a Clean Up Expired Pingfederate Persistent Sessions Plugin resource
func NewCleanUpExpiredPingfederatePersistentSessionsPluginResource() resource.Resource {
	return &cleanUpExpiredPingfederatePersistentSessionsPluginResource{}
}

func NewDefaultCleanUpExpiredPingfederatePersistentSessionsPluginResource() resource.Resource {
	return &defaultCleanUpExpiredPingfederatePersistentSessionsPluginResource{}
}

// cleanUpExpiredPingfederatePersistentSessionsPluginResource is the resource implementation.
type cleanUpExpiredPingfederatePersistentSessionsPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultCleanUpExpiredPingfederatePersistentSessionsPluginResource is the resource implementation.
type defaultCleanUpExpiredPingfederatePersistentSessionsPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *cleanUpExpiredPingfederatePersistentSessionsPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clean_up_expired_pingfederate_persistent_sessions_plugin"
}

func (r *defaultCleanUpExpiredPingfederatePersistentSessionsPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_clean_up_expired_pingfederate_persistent_sessions_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *cleanUpExpiredPingfederatePersistentSessionsPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultCleanUpExpiredPingfederatePersistentSessionsPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type cleanUpExpiredPingfederatePersistentSessionsPluginResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	LastUpdated             types.String `tfsdk:"last_updated"`
	Notifications           types.Set    `tfsdk:"notifications"`
	RequiredActions         types.Set    `tfsdk:"required_actions"`
	PollingInterval         types.String `tfsdk:"polling_interval"`
	PeerServerPriorityIndex types.Int64  `tfsdk:"peer_server_priority_index"`
	BaseDN                  types.String `tfsdk:"base_dn"`
	MaxUpdatesPerSecond     types.Int64  `tfsdk:"max_updates_per_second"`
	NumDeleteThreads        types.Int64  `tfsdk:"num_delete_threads"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *cleanUpExpiredPingfederatePersistentSessionsPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	cleanUpExpiredPingfederatePersistentSessionsPluginSchema(ctx, req, resp, false)
}

func (r *defaultCleanUpExpiredPingfederatePersistentSessionsPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	cleanUpExpiredPingfederatePersistentSessionsPluginSchema(ctx, req, resp, true)
}

func cleanUpExpiredPingfederatePersistentSessionsPluginSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Clean Up Expired Pingfederate Persistent Sessions Plugin.",
		Attributes: map[string]schema.Attribute{
			"polling_interval": schema.StringAttribute{
				Description: "This specifies how often the plugin should check for expired data. It also controls the offset of peer servers (see the peer-server-priority-index for more information).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"peer_server_priority_index": schema.Int64Attribute{
				Description: "In a replicated environment, this determines the order in which peer servers should attempt to purge data.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"base_dn": schema.StringAttribute{
				Description: "Only entries located within the subtree specified by this base DN are eligible for purging.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_updates_per_second": schema.Int64Attribute{
				Description: "This setting smooths out the performance impact on the server by throttling the purging to the specified maximum number of updates per second. To avoid a large backlog, this value should be set comfortably above the average rate that expired data is generated. When purge-behavior is set to subtree-delete-entries, then deletion of the entire subtree is considered a single update for the purposes of throttling.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"num_delete_threads": schema.Int64Attribute{
				Description: "The number of threads used to delete expired entries.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the plug-in is enabled for use.",
				Required:    true,
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalCleanUpExpiredPingfederatePersistentSessionsPluginFields(ctx context.Context, addRequest *client.AddCleanUpExpiredPingfederatePersistentSessionsPluginRequest, plan cleanUpExpiredPingfederatePersistentSessionsPluginResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PollingInterval) {
		stringVal := plan.PollingInterval.ValueString()
		addRequest.PollingInterval = &stringVal
	}
	if internaltypes.IsDefined(plan.PeerServerPriorityIndex) {
		intVal := int32(plan.PeerServerPriorityIndex.ValueInt64())
		addRequest.PeerServerPriorityIndex = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BaseDN) {
		stringVal := plan.BaseDN.ValueString()
		addRequest.BaseDN = &stringVal
	}
	if internaltypes.IsDefined(plan.MaxUpdatesPerSecond) {
		intVal := int32(plan.MaxUpdatesPerSecond.ValueInt64())
		addRequest.MaxUpdatesPerSecond = &intVal
	}
	if internaltypes.IsDefined(plan.NumDeleteThreads) {
		intVal := int32(plan.NumDeleteThreads.ValueInt64())
		addRequest.NumDeleteThreads = &intVal
	}
}

// Read a CleanUpExpiredPingfederatePersistentSessionsPluginResponse object into the model struct
func readCleanUpExpiredPingfederatePersistentSessionsPluginResponse(ctx context.Context, r *client.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, state *cleanUpExpiredPingfederatePersistentSessionsPluginResourceModel, expectedValues *cleanUpExpiredPingfederatePersistentSessionsPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.PollingInterval = types.StringValue(r.PollingInterval)
	config.CheckMismatchedPDFormattedAttributes("polling_interval",
		expectedValues.PollingInterval, state.PollingInterval, diagnostics)
	state.PeerServerPriorityIndex = internaltypes.Int64TypeOrNil(r.PeerServerPriorityIndex)
	state.BaseDN = internaltypes.StringTypeOrNil(r.BaseDN, internaltypes.IsEmptyString(expectedValues.BaseDN))
	state.MaxUpdatesPerSecond = types.Int64Value(int64(r.MaxUpdatesPerSecond))
	state.NumDeleteThreads = types.Int64Value(int64(r.NumDeleteThreads))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createCleanUpExpiredPingfederatePersistentSessionsPluginOperations(plan cleanUpExpiredPingfederatePersistentSessionsPluginResourceModel, state cleanUpExpiredPingfederatePersistentSessionsPluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.PollingInterval, state.PollingInterval, "polling-interval")
	operations.AddInt64OperationIfNecessary(&ops, plan.PeerServerPriorityIndex, state.PeerServerPriorityIndex, "peer-server-priority-index")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxUpdatesPerSecond, state.MaxUpdatesPerSecond, "max-updates-per-second")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumDeleteThreads, state.NumDeleteThreads, "num-delete-threads")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *cleanUpExpiredPingfederatePersistentSessionsPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cleanUpExpiredPingfederatePersistentSessionsPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddCleanUpExpiredPingfederatePersistentSessionsPluginRequest(plan.Id.ValueString(),
		[]client.EnumcleanUpExpiredPingfederatePersistentSessionsPluginSchemaUrn{client.ENUMCLEANUPEXPIREDPINGFEDERATEPERSISTENTSESSIONSPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINCLEAN_UP_EXPIRED_PINGFEDERATE_PERSISTENT_SESSIONS},
		plan.Enabled.ValueBool())
	addOptionalCleanUpExpiredPingfederatePersistentSessionsPluginFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddCleanUpExpiredPingfederatePersistentSessionsPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Clean Up Expired Pingfederate Persistent Sessions Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state cleanUpExpiredPingfederatePersistentSessionsPluginResourceModel
	readCleanUpExpiredPingfederatePersistentSessionsPluginResponse(ctx, addResponse.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultCleanUpExpiredPingfederatePersistentSessionsPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cleanUpExpiredPingfederatePersistentSessionsPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Clean Up Expired Pingfederate Persistent Sessions Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state cleanUpExpiredPingfederatePersistentSessionsPluginResourceModel
	readCleanUpExpiredPingfederatePersistentSessionsPluginResponse(ctx, readResponse.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createCleanUpExpiredPingfederatePersistentSessionsPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Clean Up Expired Pingfederate Persistent Sessions Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCleanUpExpiredPingfederatePersistentSessionsPluginResponse(ctx, updateResponse.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, &state, &plan, &resp.Diagnostics)
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *cleanUpExpiredPingfederatePersistentSessionsPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCleanUpExpiredPingfederatePersistentSessionsPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCleanUpExpiredPingfederatePersistentSessionsPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCleanUpExpiredPingfederatePersistentSessionsPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readCleanUpExpiredPingfederatePersistentSessionsPlugin(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state cleanUpExpiredPingfederatePersistentSessionsPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Clean Up Expired Pingfederate Persistent Sessions Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCleanUpExpiredPingfederatePersistentSessionsPluginResponse(ctx, readResponse.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *cleanUpExpiredPingfederatePersistentSessionsPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCleanUpExpiredPingfederatePersistentSessionsPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCleanUpExpiredPingfederatePersistentSessionsPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCleanUpExpiredPingfederatePersistentSessionsPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateCleanUpExpiredPingfederatePersistentSessionsPlugin(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan cleanUpExpiredPingfederatePersistentSessionsPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state cleanUpExpiredPingfederatePersistentSessionsPluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createCleanUpExpiredPingfederatePersistentSessionsPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Clean Up Expired Pingfederate Persistent Sessions Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCleanUpExpiredPingfederatePersistentSessionsPluginResponse(ctx, updateResponse.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, &state, &plan, &resp.Diagnostics)
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultCleanUpExpiredPingfederatePersistentSessionsPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *cleanUpExpiredPingfederatePersistentSessionsPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state cleanUpExpiredPingfederatePersistentSessionsPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PluginApi.DeletePluginExecute(r.apiClient.PluginApi.DeletePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Clean Up Expired Pingfederate Persistent Sessions Plugin", err, httpResp)
		return
	}
}

func (r *cleanUpExpiredPingfederatePersistentSessionsPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCleanUpExpiredPingfederatePersistentSessionsPlugin(ctx, req, resp)
}

func (r *defaultCleanUpExpiredPingfederatePersistentSessionsPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCleanUpExpiredPingfederatePersistentSessionsPlugin(ctx, req, resp)
}

func importCleanUpExpiredPingfederatePersistentSessionsPlugin(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
