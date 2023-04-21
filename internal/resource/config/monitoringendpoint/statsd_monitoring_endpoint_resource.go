package monitoringendpoint

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &statsdMonitoringEndpointResource{}
	_ resource.ResourceWithConfigure   = &statsdMonitoringEndpointResource{}
	_ resource.ResourceWithImportState = &statsdMonitoringEndpointResource{}
	_ resource.Resource                = &defaultStatsdMonitoringEndpointResource{}
	_ resource.ResourceWithConfigure   = &defaultStatsdMonitoringEndpointResource{}
	_ resource.ResourceWithImportState = &defaultStatsdMonitoringEndpointResource{}
)

// Create a Statsd Monitoring Endpoint resource
func NewStatsdMonitoringEndpointResource() resource.Resource {
	return &statsdMonitoringEndpointResource{}
}

func NewDefaultStatsdMonitoringEndpointResource() resource.Resource {
	return &defaultStatsdMonitoringEndpointResource{}
}

// statsdMonitoringEndpointResource is the resource implementation.
type statsdMonitoringEndpointResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultStatsdMonitoringEndpointResource is the resource implementation.
type defaultStatsdMonitoringEndpointResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *statsdMonitoringEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_statsd_monitoring_endpoint"
}

func (r *defaultStatsdMonitoringEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_statsd_monitoring_endpoint"
}

// Configure adds the provider configured client to the resource.
func (r *statsdMonitoringEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultStatsdMonitoringEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type statsdMonitoringEndpointResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	Notifications        types.Set    `tfsdk:"notifications"`
	RequiredActions      types.Set    `tfsdk:"required_actions"`
	Hostname             types.String `tfsdk:"hostname"`
	ServerPort           types.Int64  `tfsdk:"server_port"`
	ConnectionType       types.String `tfsdk:"connection_type"`
	TrustManagerProvider types.String `tfsdk:"trust_manager_provider"`
	AdditionalTags       types.Set    `tfsdk:"additional_tags"`
	Enabled              types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *statsdMonitoringEndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	statsdMonitoringEndpointSchema(ctx, req, resp, false)
}

func (r *defaultStatsdMonitoringEndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	statsdMonitoringEndpointSchema(ctx, req, resp, true)
}

func statsdMonitoringEndpointSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Statsd Monitoring Endpoint.",
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				Description: "The name of the host where this StatsD Monitoring Endpoint should send metric data.",
				Required:    true,
			},
			"server_port": schema.Int64Attribute{
				Description: "Specifies the port number of the endpoint where metric data should be sent.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"connection_type": schema.StringAttribute{
				Description: "Specifies the protocol and security that this StatsD Monitoring Endpoint should use to connect to the configured endpoint.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "The trust manager provider to use if SSL over TCP is to be used for connection-level security.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"additional_tags": schema.SetAttribute{
				Description: "Specifies any optional additional tags to include in StatsD messages. Any additional tags will be appended to the end of each StatsD message, separated by commas. Tags should be written in a [key]:[value] format (\"host:server1\", for example).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Monitoring Endpoint is enabled for use in the Directory Server.",
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
func addOptionalStatsdMonitoringEndpointFields(ctx context.Context, addRequest *client.AddStatsdMonitoringEndpointRequest, plan statsdMonitoringEndpointResourceModel) error {
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionType) {
		connectionType, err := client.NewEnummonitoringEndpointConnectionTypePropFromValue(plan.ConnectionType.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConnectionType = connectionType
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AdditionalTags) {
		var slice []string
		plan.AdditionalTags.ElementsAs(ctx, &slice, false)
		addRequest.AdditionalTags = slice
	}
	return nil
}

// Read a StatsdMonitoringEndpointResponse object into the model struct
func readStatsdMonitoringEndpointResponse(ctx context.Context, r *client.StatsdMonitoringEndpointResponse, state *statsdMonitoringEndpointResourceModel, expectedValues *statsdMonitoringEndpointResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Hostname = types.StringValue(r.Hostname)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.ConnectionType = types.StringValue(r.ConnectionType.String())
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, internaltypes.IsEmptyString(expectedValues.TrustManagerProvider))
	state.AdditionalTags = internaltypes.GetStringSet(r.AdditionalTags)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createStatsdMonitoringEndpointOperations(plan statsdMonitoringEndpointResourceModel, state statsdMonitoringEndpointResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Hostname, state.Hostname, "hostname")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerPort, state.ServerPort, "server-port")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionType, state.ConnectionType, "connection-type")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustManagerProvider, state.TrustManagerProvider, "trust-manager-provider")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AdditionalTags, state.AdditionalTags, "additional-tags")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *statsdMonitoringEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan statsdMonitoringEndpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddStatsdMonitoringEndpointRequest(plan.Id.ValueString(),
		[]client.EnumstatsdMonitoringEndpointSchemaUrn{client.ENUMSTATSDMONITORINGENDPOINTSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0MONITORING_ENDPOINTSTATSD},
		plan.Hostname.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalStatsdMonitoringEndpointFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Statsd Monitoring Endpoint", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.MonitoringEndpointApi.AddMonitoringEndpoint(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddStatsdMonitoringEndpointRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.MonitoringEndpointApi.AddMonitoringEndpointExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Statsd Monitoring Endpoint", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state statsdMonitoringEndpointResourceModel
	readStatsdMonitoringEndpointResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultStatsdMonitoringEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan statsdMonitoringEndpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.MonitoringEndpointApi.GetMonitoringEndpoint(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Statsd Monitoring Endpoint", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state statsdMonitoringEndpointResourceModel
	readStatsdMonitoringEndpointResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.MonitoringEndpointApi.UpdateMonitoringEndpoint(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createStatsdMonitoringEndpointOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.MonitoringEndpointApi.UpdateMonitoringEndpointExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Statsd Monitoring Endpoint", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readStatsdMonitoringEndpointResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *statsdMonitoringEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readStatsdMonitoringEndpoint(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultStatsdMonitoringEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readStatsdMonitoringEndpoint(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readStatsdMonitoringEndpoint(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state statsdMonitoringEndpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.MonitoringEndpointApi.GetMonitoringEndpoint(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Statsd Monitoring Endpoint", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readStatsdMonitoringEndpointResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *statsdMonitoringEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateStatsdMonitoringEndpoint(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultStatsdMonitoringEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateStatsdMonitoringEndpoint(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateStatsdMonitoringEndpoint(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan statsdMonitoringEndpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state statsdMonitoringEndpointResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.MonitoringEndpointApi.UpdateMonitoringEndpoint(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createStatsdMonitoringEndpointOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.MonitoringEndpointApi.UpdateMonitoringEndpointExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Statsd Monitoring Endpoint", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readStatsdMonitoringEndpointResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultStatsdMonitoringEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *statsdMonitoringEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state statsdMonitoringEndpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.MonitoringEndpointApi.DeleteMonitoringEndpointExecute(r.apiClient.MonitoringEndpointApi.DeleteMonitoringEndpoint(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Statsd Monitoring Endpoint", err, httpResp)
		return
	}
}

func (r *statsdMonitoringEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStatsdMonitoringEndpoint(ctx, req, resp)
}

func (r *defaultStatsdMonitoringEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStatsdMonitoringEndpoint(ctx, req, resp)
}

func importStatsdMonitoringEndpoint(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
