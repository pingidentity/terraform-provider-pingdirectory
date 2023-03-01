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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &snmpSubagentPluginResource{}
	_ resource.ResourceWithConfigure   = &snmpSubagentPluginResource{}
	_ resource.ResourceWithImportState = &snmpSubagentPluginResource{}
	_ resource.Resource                = &defaultSnmpSubagentPluginResource{}
	_ resource.ResourceWithConfigure   = &defaultSnmpSubagentPluginResource{}
	_ resource.ResourceWithImportState = &defaultSnmpSubagentPluginResource{}
)

// Create a Snmp Subagent Plugin resource
func NewSnmpSubagentPluginResource() resource.Resource {
	return &snmpSubagentPluginResource{}
}

func NewDefaultSnmpSubagentPluginResource() resource.Resource {
	return &defaultSnmpSubagentPluginResource{}
}

// snmpSubagentPluginResource is the resource implementation.
type snmpSubagentPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSnmpSubagentPluginResource is the resource implementation.
type defaultSnmpSubagentPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *snmpSubagentPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snmp_subagent_plugin"
}

func (r *defaultSnmpSubagentPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_snmp_subagent_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *snmpSubagentPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultSnmpSubagentPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type snmpSubagentPluginResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	LastUpdated                 types.String `tfsdk:"last_updated"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	ContextName                 types.String `tfsdk:"context_name"`
	AgentxAddress               types.String `tfsdk:"agentx_address"`
	AgentxPort                  types.Int64  `tfsdk:"agentx_port"`
	NumWorkerThreads            types.Int64  `tfsdk:"num_worker_threads"`
	SessionTimeout              types.String `tfsdk:"session_timeout"`
	ConnectRetryMaxWait         types.String `tfsdk:"connect_retry_max_wait"`
	PingInterval                types.String `tfsdk:"ping_interval"`
	Description                 types.String `tfsdk:"description"`
	Enabled                     types.Bool   `tfsdk:"enabled"`
	InvokeForInternalOperations types.Bool   `tfsdk:"invoke_for_internal_operations"`
}

// GetSchema defines the schema for the resource.
func (r *snmpSubagentPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	snmpSubagentPluginSchema(ctx, req, resp, false)
}

func (r *defaultSnmpSubagentPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	snmpSubagentPluginSchema(ctx, req, resp, true)
}

func snmpSubagentPluginSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Snmp Subagent Plugin.",
		Attributes: map[string]schema.Attribute{
			"context_name": schema.StringAttribute{
				Description: "The SNMP context name for this sub-agent. The context name must not be longer than 30 ASCII characters. Each server in a topology must have a unique SNMP context name.",
				Optional:    true,
				Computed:    true,
			},
			"agentx_address": schema.StringAttribute{
				Description: "The hostname or IP address of the SNMP master agent.",
				Optional:    true,
				Computed:    true,
			},
			"agentx_port": schema.Int64Attribute{
				Description: "The port number on which the SNMP master agent will be contacted.",
				Optional:    true,
				Computed:    true,
			},
			"num_worker_threads": schema.Int64Attribute{
				Description: "The number of worker threads to use to handle SNMP requests.",
				Optional:    true,
				Computed:    true,
			},
			"session_timeout": schema.StringAttribute{
				Description: "Specifies the maximum amount of time to wait for a session to the master agent to be established.",
				Optional:    true,
				Computed:    true,
			},
			"connect_retry_max_wait": schema.StringAttribute{
				Description: "The maximum amount of time to wait between attempts to establish a connection to the master agent.",
				Optional:    true,
				Computed:    true,
			},
			"ping_interval": schema.StringAttribute{
				Description: "The amount of time between consecutive pings sent by the sub-agent on its connection to the master agent. A value of zero disables the sending of pings by the sub-agent.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Plugin",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the plug-in is enabled for use.",
				Required:    true,
			},
			"invoke_for_internal_operations": schema.BoolAttribute{
				Description: "Indicates whether the plug-in should be invoked for internal operations.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	if setOptionalToComputed {
		config.SetOptionalAttributesToComputed(&schema)
	}
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalSnmpSubagentPluginFields(ctx context.Context, addRequest *client.AddSnmpSubagentPluginRequest, plan snmpSubagentPluginResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ContextName) {
		stringVal := plan.ContextName.ValueString()
		addRequest.ContextName = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AgentxAddress) {
		stringVal := plan.AgentxAddress.ValueString()
		addRequest.AgentxAddress = &stringVal
	}
	if internaltypes.IsDefined(plan.AgentxPort) {
		intVal := int32(plan.AgentxPort.ValueInt64())
		addRequest.AgentxPort = &intVal
	}
	if internaltypes.IsDefined(plan.NumWorkerThreads) {
		intVal := int32(plan.NumWorkerThreads.ValueInt64())
		addRequest.NumWorkerThreads = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SessionTimeout) {
		stringVal := plan.SessionTimeout.ValueString()
		addRequest.SessionTimeout = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectRetryMaxWait) {
		stringVal := plan.ConnectRetryMaxWait.ValueString()
		addRequest.ConnectRetryMaxWait = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PingInterval) {
		stringVal := plan.PingInterval.ValueString()
		addRequest.PingInterval = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		boolVal := plan.InvokeForInternalOperations.ValueBool()
		addRequest.InvokeForInternalOperations = &boolVal
	}
}

// Read a SnmpSubagentPluginResponse object into the model struct
func readSnmpSubagentPluginResponse(ctx context.Context, r *client.SnmpSubagentPluginResponse, state *snmpSubagentPluginResourceModel, expectedValues *snmpSubagentPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ContextName = internaltypes.StringTypeOrNil(r.ContextName, internaltypes.IsEmptyString(expectedValues.ContextName))
	state.AgentxAddress = types.StringValue(r.AgentxAddress)
	state.AgentxPort = types.Int64Value(int64(r.AgentxPort))
	state.NumWorkerThreads = internaltypes.Int64TypeOrNil(r.NumWorkerThreads)
	state.SessionTimeout = internaltypes.StringTypeOrNil(r.SessionTimeout, internaltypes.IsEmptyString(expectedValues.SessionTimeout))
	config.CheckMismatchedPDFormattedAttributes("session_timeout",
		expectedValues.SessionTimeout, state.SessionTimeout, diagnostics)
	state.ConnectRetryMaxWait = internaltypes.StringTypeOrNil(r.ConnectRetryMaxWait, internaltypes.IsEmptyString(expectedValues.ConnectRetryMaxWait))
	config.CheckMismatchedPDFormattedAttributes("connect_retry_max_wait",
		expectedValues.ConnectRetryMaxWait, state.ConnectRetryMaxWait, diagnostics)
	state.PingInterval = internaltypes.StringTypeOrNil(r.PingInterval, internaltypes.IsEmptyString(expectedValues.PingInterval))
	config.CheckMismatchedPDFormattedAttributes("ping_interval",
		expectedValues.PingInterval, state.PingInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSnmpSubagentPluginOperations(plan snmpSubagentPluginResourceModel, state snmpSubagentPluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ContextName, state.ContextName, "context-name")
	operations.AddStringOperationIfNecessary(&ops, plan.AgentxAddress, state.AgentxAddress, "agentx-address")
	operations.AddInt64OperationIfNecessary(&ops, plan.AgentxPort, state.AgentxPort, "agentx-port")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumWorkerThreads, state.NumWorkerThreads, "num-worker-threads")
	operations.AddStringOperationIfNecessary(&ops, plan.SessionTimeout, state.SessionTimeout, "session-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectRetryMaxWait, state.ConnectRetryMaxWait, "connect-retry-max-wait")
	operations.AddStringOperationIfNecessary(&ops, plan.PingInterval, state.PingInterval, "ping-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.InvokeForInternalOperations, state.InvokeForInternalOperations, "invoke-for-internal-operations")
	return ops
}

// Create a new resource
func (r *snmpSubagentPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan snmpSubagentPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddSnmpSubagentPluginRequest(plan.Id.ValueString(),
		[]client.EnumsnmpSubagentPluginSchemaUrn{client.ENUMSNMPSUBAGENTPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINSNMP_SUBAGENT},
		plan.Enabled.ValueBool())
	addOptionalSnmpSubagentPluginFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddSnmpSubagentPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Snmp Subagent Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state snmpSubagentPluginResourceModel
	readSnmpSubagentPluginResponse(ctx, addResponse.SnmpSubagentPluginResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSnmpSubagentPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan snmpSubagentPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Snmp Subagent Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state snmpSubagentPluginResourceModel
	readSnmpSubagentPluginResponse(ctx, readResponse.SnmpSubagentPluginResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSnmpSubagentPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Snmp Subagent Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSnmpSubagentPluginResponse(ctx, updateResponse.SnmpSubagentPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *snmpSubagentPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSnmpSubagentPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSnmpSubagentPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSnmpSubagentPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSnmpSubagentPlugin(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state snmpSubagentPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Snmp Subagent Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSnmpSubagentPluginResponse(ctx, readResponse.SnmpSubagentPluginResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *snmpSubagentPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSnmpSubagentPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSnmpSubagentPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSnmpSubagentPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSnmpSubagentPlugin(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan snmpSubagentPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state snmpSubagentPluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSnmpSubagentPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Snmp Subagent Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSnmpSubagentPluginResponse(ctx, updateResponse.SnmpSubagentPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSnmpSubagentPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *snmpSubagentPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state snmpSubagentPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PluginApi.DeletePluginExecute(r.apiClient.PluginApi.DeletePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Snmp Subagent Plugin", err, httpResp)
		return
	}
}

func (r *snmpSubagentPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSnmpSubagentPlugin(ctx, req, resp)
}

func (r *defaultSnmpSubagentPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSnmpSubagentPlugin(ctx, req, resp)
}

func importSnmpSubagentPlugin(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
