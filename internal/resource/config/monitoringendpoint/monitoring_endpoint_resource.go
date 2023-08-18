package monitoringendpoint

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &monitoringEndpointResource{}
	_ resource.ResourceWithConfigure   = &monitoringEndpointResource{}
	_ resource.ResourceWithImportState = &monitoringEndpointResource{}
	_ resource.Resource                = &defaultMonitoringEndpointResource{}
	_ resource.ResourceWithConfigure   = &defaultMonitoringEndpointResource{}
	_ resource.ResourceWithImportState = &defaultMonitoringEndpointResource{}
)

// Create a Monitoring Endpoint resource
func NewMonitoringEndpointResource() resource.Resource {
	return &monitoringEndpointResource{}
}

func NewDefaultMonitoringEndpointResource() resource.Resource {
	return &defaultMonitoringEndpointResource{}
}

// monitoringEndpointResource is the resource implementation.
type monitoringEndpointResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultMonitoringEndpointResource is the resource implementation.
type defaultMonitoringEndpointResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *monitoringEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitoring_endpoint"
}

func (r *defaultMonitoringEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_monitoring_endpoint"
}

// Configure adds the provider configured client to the resource.
func (r *monitoringEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultMonitoringEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type monitoringEndpointResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	Notifications        types.Set    `tfsdk:"notifications"`
	RequiredActions      types.Set    `tfsdk:"required_actions"`
	Type                 types.String `tfsdk:"type"`
	Hostname             types.String `tfsdk:"hostname"`
	ServerPort           types.Int64  `tfsdk:"server_port"`
	ConnectionType       types.String `tfsdk:"connection_type"`
	TrustManagerProvider types.String `tfsdk:"trust_manager_provider"`
	AdditionalTags       types.Set    `tfsdk:"additional_tags"`
	Enabled              types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *monitoringEndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	monitoringEndpointSchema(ctx, req, resp, false)
}

func (r *defaultMonitoringEndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	monitoringEndpointSchema(ctx, req, resp, true)
}

func monitoringEndpointSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Monitoring Endpoint.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Monitoring Endpoint resource. Options are ['statsd']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("statsd"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"statsd"}...),
				},
			},
			"hostname": schema.StringAttribute{
				Description: "The name of the host where this StatsD Monitoring Endpoint should send metric data.",
				Required:    true,
			},
			"server_port": schema.Int64Attribute{
				Description: "Specifies the port number of the endpoint where metric data should be sent.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(8125),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"connection_type": schema.StringAttribute{
				Description: "Specifies the protocol and security that this StatsD Monitoring Endpoint should use to connect to the configured endpoint.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("unencrypted-udp"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "The trust manager provider to use if SSL over TCP is to be used for connection-level security.",
				Optional:    true,
			},
			"additional_tags": schema.SetAttribute{
				Description: "Specifies any optional additional tags to include in StatsD messages. Any additional tags will be appended to the end of each StatsD message, separated by commas. Tags should be written in a [key]:[value] format (\"host:server1\", for example).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Monitoring Endpoint is enabled for use in the Directory Server.",
				Required:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for statsd monitoring-endpoint
func addOptionalStatsdMonitoringEndpointFields(ctx context.Context, addRequest *client.AddStatsdMonitoringEndpointRequest, plan monitoringEndpointResourceModel) error {
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

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *monitoringEndpointResourceModel) populateAllComputedStringAttributes() {
	if model.TrustManagerProvider.IsUnknown() || model.TrustManagerProvider.IsNull() {
		model.TrustManagerProvider = types.StringValue("")
	}
	if model.ConnectionType.IsUnknown() || model.ConnectionType.IsNull() {
		model.ConnectionType = types.StringValue("")
	}
	if model.Hostname.IsUnknown() || model.Hostname.IsNull() {
		model.Hostname = types.StringValue("")
	}
}

// Read a StatsdMonitoringEndpointResponse object into the model struct
func readStatsdMonitoringEndpointResponse(ctx context.Context, r *client.StatsdMonitoringEndpointResponse, state *monitoringEndpointResourceModel, expectedValues *monitoringEndpointResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("statsd")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Hostname = types.StringValue(r.Hostname)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.ConnectionType = types.StringValue(r.ConnectionType.String())
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, internaltypes.IsEmptyString(expectedValues.TrustManagerProvider))
	state.AdditionalTags = internaltypes.GetStringSet(r.AdditionalTags)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createMonitoringEndpointOperations(plan monitoringEndpointResourceModel, state monitoringEndpointResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Hostname, state.Hostname, "hostname")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerPort, state.ServerPort, "server-port")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionType, state.ConnectionType, "connection-type")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustManagerProvider, state.TrustManagerProvider, "trust-manager-provider")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AdditionalTags, state.AdditionalTags, "additional-tags")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a statsd monitoring-endpoint
func (r *monitoringEndpointResource) CreateStatsdMonitoringEndpoint(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan monitoringEndpointResourceModel) (*monitoringEndpointResourceModel, error) {
	addRequest := client.NewAddStatsdMonitoringEndpointRequest(plan.Name.ValueString(),
		[]client.EnumstatsdMonitoringEndpointSchemaUrn{client.ENUMSTATSDMONITORINGENDPOINTSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0MONITORING_ENDPOINTSTATSD},
		plan.Hostname.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalStatsdMonitoringEndpointFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Monitoring Endpoint", err.Error())
		return nil, err
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Monitoring Endpoint", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state monitoringEndpointResourceModel
	readStatsdMonitoringEndpointResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *monitoringEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan monitoringEndpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateStatsdMonitoringEndpoint(ctx, req, resp, plan)
	if err != nil {
		return
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, *state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultMonitoringEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan monitoringEndpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.MonitoringEndpointApi.GetMonitoringEndpoint(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Monitoring Endpoint", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state monitoringEndpointResourceModel
	readStatsdMonitoringEndpointResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.MonitoringEndpointApi.UpdateMonitoringEndpoint(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createMonitoringEndpointOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.MonitoringEndpointApi.UpdateMonitoringEndpointExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Monitoring Endpoint", err, httpResp)
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

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *monitoringEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readMonitoringEndpoint(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultMonitoringEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readMonitoringEndpoint(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readMonitoringEndpoint(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state monitoringEndpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.MonitoringEndpointApi.GetMonitoringEndpoint(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Monitoring Endpoint", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Monitoring Endpoint", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readStatsdMonitoringEndpointResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *monitoringEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateMonitoringEndpoint(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultMonitoringEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateMonitoringEndpoint(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateMonitoringEndpoint(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan monitoringEndpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state monitoringEndpointResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.MonitoringEndpointApi.UpdateMonitoringEndpoint(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createMonitoringEndpointOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.MonitoringEndpointApi.UpdateMonitoringEndpointExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Monitoring Endpoint", err, httpResp)
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
func (r *defaultMonitoringEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *monitoringEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state monitoringEndpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.MonitoringEndpointApi.DeleteMonitoringEndpointExecute(r.apiClient.MonitoringEndpointApi.DeleteMonitoringEndpoint(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Monitoring Endpoint", err, httpResp)
		return
	}
}

func (r *monitoringEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importMonitoringEndpoint(ctx, req, resp)
}

func (r *defaultMonitoringEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importMonitoringEndpoint(ctx, req, resp)
}

func importMonitoringEndpoint(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
