package plugin

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

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
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &cleanUpInactivePingfederatePersistentSessionsPluginResource{}
	_ resource.ResourceWithConfigure   = &cleanUpInactivePingfederatePersistentSessionsPluginResource{}
	_ resource.ResourceWithImportState = &cleanUpInactivePingfederatePersistentSessionsPluginResource{}
	_ resource.Resource                = &defaultCleanUpInactivePingfederatePersistentSessionsPluginResource{}
	_ resource.ResourceWithConfigure   = &defaultCleanUpInactivePingfederatePersistentSessionsPluginResource{}
	_ resource.ResourceWithImportState = &defaultCleanUpInactivePingfederatePersistentSessionsPluginResource{}
)

// Create a Clean Up Inactive Pingfederate Persistent Sessions Plugin resource
func NewCleanUpInactivePingfederatePersistentSessionsPluginResource() resource.Resource {
	return &cleanUpInactivePingfederatePersistentSessionsPluginResource{}
}

func NewDefaultCleanUpInactivePingfederatePersistentSessionsPluginResource() resource.Resource {
	return &defaultCleanUpInactivePingfederatePersistentSessionsPluginResource{}
}

// cleanUpInactivePingfederatePersistentSessionsPluginResource is the resource implementation.
type cleanUpInactivePingfederatePersistentSessionsPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultCleanUpInactivePingfederatePersistentSessionsPluginResource is the resource implementation.
type defaultCleanUpInactivePingfederatePersistentSessionsPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *cleanUpInactivePingfederatePersistentSessionsPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clean_up_inactive_pingfederate_persistent_sessions_plugin"
}

func (r *defaultCleanUpInactivePingfederatePersistentSessionsPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_clean_up_inactive_pingfederate_persistent_sessions_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *cleanUpInactivePingfederatePersistentSessionsPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultCleanUpInactivePingfederatePersistentSessionsPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type cleanUpInactivePingfederatePersistentSessionsPluginResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	LastUpdated             types.String `tfsdk:"last_updated"`
	Notifications           types.Set    `tfsdk:"notifications"`
	RequiredActions         types.Set    `tfsdk:"required_actions"`
	ExpirationOffset        types.String `tfsdk:"expiration_offset"`
	PollingInterval         types.String `tfsdk:"polling_interval"`
	PeerServerPriorityIndex types.Int64  `tfsdk:"peer_server_priority_index"`
	BaseDN                  types.String `tfsdk:"base_dn"`
	MaxUpdatesPerSecond     types.Int64  `tfsdk:"max_updates_per_second"`
	NumDeleteThreads        types.Int64  `tfsdk:"num_delete_threads"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *cleanUpInactivePingfederatePersistentSessionsPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	cleanUpInactivePingfederatePersistentSessionsPluginSchema(ctx, req, resp, false)
}

func (r *defaultCleanUpInactivePingfederatePersistentSessionsPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	cleanUpInactivePingfederatePersistentSessionsPluginSchema(ctx, req, resp, true)
}

func cleanUpInactivePingfederatePersistentSessionsPluginSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Clean Up Inactive Pingfederate Persistent Sessions Plugin.",
		Attributes: map[string]schema.Attribute{
			"expiration_offset": schema.StringAttribute{
				Description: "Sessions whose last activity timestamp is older than this offset will be removed.",
				Required:    true,
			},
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
func addOptionalCleanUpInactivePingfederatePersistentSessionsPluginFields(ctx context.Context, addRequest *client.AddCleanUpInactivePingfederatePersistentSessionsPluginRequest, plan cleanUpInactivePingfederatePersistentSessionsPluginResourceModel) {
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

// Read a CleanUpInactivePingfederatePersistentSessionsPluginResponse object into the model struct
func readCleanUpInactivePingfederatePersistentSessionsPluginResponse(ctx context.Context, r *client.CleanUpInactivePingfederatePersistentSessionsPluginResponse, state *cleanUpInactivePingfederatePersistentSessionsPluginResourceModel, expectedValues *cleanUpInactivePingfederatePersistentSessionsPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ExpirationOffset = types.StringValue(r.ExpirationOffset)
	config.CheckMismatchedPDFormattedAttributes("expiration_offset",
		expectedValues.ExpirationOffset, state.ExpirationOffset, diagnostics)
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
func createCleanUpInactivePingfederatePersistentSessionsPluginOperations(plan cleanUpInactivePingfederatePersistentSessionsPluginResourceModel, state cleanUpInactivePingfederatePersistentSessionsPluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExpirationOffset, state.ExpirationOffset, "expiration-offset")
	operations.AddStringOperationIfNecessary(&ops, plan.PollingInterval, state.PollingInterval, "polling-interval")
	operations.AddInt64OperationIfNecessary(&ops, plan.PeerServerPriorityIndex, state.PeerServerPriorityIndex, "peer-server-priority-index")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxUpdatesPerSecond, state.MaxUpdatesPerSecond, "max-updates-per-second")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumDeleteThreads, state.NumDeleteThreads, "num-delete-threads")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *cleanUpInactivePingfederatePersistentSessionsPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cleanUpInactivePingfederatePersistentSessionsPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddCleanUpInactivePingfederatePersistentSessionsPluginRequest(plan.Id.ValueString(),
		[]client.EnumcleanUpInactivePingfederatePersistentSessionsPluginSchemaUrn{client.ENUMCLEANUPINACTIVEPINGFEDERATEPERSISTENTSESSIONSPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINCLEAN_UP_INACTIVE_PINGFEDERATE_PERSISTENT_SESSIONS},
		plan.ExpirationOffset.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalCleanUpInactivePingfederatePersistentSessionsPluginFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddCleanUpInactivePingfederatePersistentSessionsPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Clean Up Inactive Pingfederate Persistent Sessions Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state cleanUpInactivePingfederatePersistentSessionsPluginResourceModel
	readCleanUpInactivePingfederatePersistentSessionsPluginResponse(ctx, addResponse.CleanUpInactivePingfederatePersistentSessionsPluginResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultCleanUpInactivePingfederatePersistentSessionsPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cleanUpInactivePingfederatePersistentSessionsPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Clean Up Inactive Pingfederate Persistent Sessions Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state cleanUpInactivePingfederatePersistentSessionsPluginResourceModel
	readCleanUpInactivePingfederatePersistentSessionsPluginResponse(ctx, readResponse.CleanUpInactivePingfederatePersistentSessionsPluginResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createCleanUpInactivePingfederatePersistentSessionsPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Clean Up Inactive Pingfederate Persistent Sessions Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCleanUpInactivePingfederatePersistentSessionsPluginResponse(ctx, updateResponse.CleanUpInactivePingfederatePersistentSessionsPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *cleanUpInactivePingfederatePersistentSessionsPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCleanUpInactivePingfederatePersistentSessionsPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCleanUpInactivePingfederatePersistentSessionsPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCleanUpInactivePingfederatePersistentSessionsPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readCleanUpInactivePingfederatePersistentSessionsPlugin(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state cleanUpInactivePingfederatePersistentSessionsPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Clean Up Inactive Pingfederate Persistent Sessions Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCleanUpInactivePingfederatePersistentSessionsPluginResponse(ctx, readResponse.CleanUpInactivePingfederatePersistentSessionsPluginResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *cleanUpInactivePingfederatePersistentSessionsPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCleanUpInactivePingfederatePersistentSessionsPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCleanUpInactivePingfederatePersistentSessionsPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCleanUpInactivePingfederatePersistentSessionsPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateCleanUpInactivePingfederatePersistentSessionsPlugin(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan cleanUpInactivePingfederatePersistentSessionsPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state cleanUpInactivePingfederatePersistentSessionsPluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createCleanUpInactivePingfederatePersistentSessionsPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Clean Up Inactive Pingfederate Persistent Sessions Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCleanUpInactivePingfederatePersistentSessionsPluginResponse(ctx, updateResponse.CleanUpInactivePingfederatePersistentSessionsPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultCleanUpInactivePingfederatePersistentSessionsPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *cleanUpInactivePingfederatePersistentSessionsPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state cleanUpInactivePingfederatePersistentSessionsPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PluginApi.DeletePluginExecute(r.apiClient.PluginApi.DeletePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Clean Up Inactive Pingfederate Persistent Sessions Plugin", err, httpResp)
		return
	}
}

func (r *cleanUpInactivePingfederatePersistentSessionsPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCleanUpInactivePingfederatePersistentSessionsPlugin(ctx, req, resp)
}

func (r *defaultCleanUpInactivePingfederatePersistentSessionsPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCleanUpInactivePingfederatePersistentSessionsPlugin(ctx, req, resp)
}

func importCleanUpInactivePingfederatePersistentSessionsPlugin(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
