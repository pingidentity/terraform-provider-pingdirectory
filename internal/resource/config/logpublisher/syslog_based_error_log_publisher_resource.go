package logpublisher

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &syslogBasedErrorLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &syslogBasedErrorLogPublisherResource{}
	_ resource.ResourceWithImportState = &syslogBasedErrorLogPublisherResource{}
	_ resource.Resource                = &defaultSyslogBasedErrorLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &defaultSyslogBasedErrorLogPublisherResource{}
	_ resource.ResourceWithImportState = &defaultSyslogBasedErrorLogPublisherResource{}
)

// Create a Syslog Based Error Log Publisher resource
func NewSyslogBasedErrorLogPublisherResource() resource.Resource {
	return &syslogBasedErrorLogPublisherResource{}
}

func NewDefaultSyslogBasedErrorLogPublisherResource() resource.Resource {
	return &defaultSyslogBasedErrorLogPublisherResource{}
}

// syslogBasedErrorLogPublisherResource is the resource implementation.
type syslogBasedErrorLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSyslogBasedErrorLogPublisherResource is the resource implementation.
type defaultSyslogBasedErrorLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *syslogBasedErrorLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_syslog_based_error_log_publisher"
}

func (r *defaultSyslogBasedErrorLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_syslog_based_error_log_publisher"
}

// Configure adds the provider configured client to the resource.
func (r *syslogBasedErrorLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultSyslogBasedErrorLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type syslogBasedErrorLogPublisherResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	Notifications        types.Set    `tfsdk:"notifications"`
	RequiredActions      types.Set    `tfsdk:"required_actions"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	ServerHostName       types.String `tfsdk:"server_host_name"`
	ServerPort           types.Int64  `tfsdk:"server_port"`
	SyslogFacility       types.Int64  `tfsdk:"syslog_facility"`
	AutoFlush            types.Bool   `tfsdk:"auto_flush"`
	Asynchronous         types.Bool   `tfsdk:"asynchronous"`
	QueueSize            types.Int64  `tfsdk:"queue_size"`
	DefaultSeverity      types.Set    `tfsdk:"default_severity"`
	OverrideSeverity     types.Set    `tfsdk:"override_severity"`
	Description          types.String `tfsdk:"description"`
	LoggingErrorBehavior types.String `tfsdk:"logging_error_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *syslogBasedErrorLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	syslogBasedErrorLogPublisherSchema(ctx, req, resp, false)
}

func (r *defaultSyslogBasedErrorLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	syslogBasedErrorLogPublisherSchema(ctx, req, resp, true)
}

func syslogBasedErrorLogPublisherSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Syslog Based Error Log Publisher.",
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Syslog Based Error Log Publisher is enabled for use.",
				Required:    true,
			},
			"server_host_name": schema.StringAttribute{
				Description: "Specifies the hostname or IP address of the syslogd host to log to. It is highly recommend to use localhost.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_port": schema.Int64Attribute{
				Description: "Specifies the port number of the syslogd host to log to.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"syslog_facility": schema.Int64Attribute{
				Description: "Specifies the syslog facility to use for this Syslog Based Error Log Publisher",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"auto_flush": schema.BoolAttribute{
				Description: "Specifies whether to flush the writer after every log record.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"asynchronous": schema.BoolAttribute{
				Description: "Indicates whether the Syslog Based Error Log Publisher will publish records asynchronously.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"queue_size": schema.Int64Attribute{
				Description: "The maximum number of log records that can be stored in the asynchronous queue.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"default_severity": schema.SetAttribute{
				Description: "Specifies the default severity levels for the logger.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"override_severity": schema.SetAttribute{
				Description: "Specifies the override severity levels for the logger based on the category of the messages.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Publisher",
				Optional:    true,
			},
			"logging_error_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if an error occurs during logging processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
func addOptionalSyslogBasedErrorLogPublisherFields(ctx context.Context, addRequest *client.AddSyslogBasedErrorLogPublisherRequest, plan syslogBasedErrorLogPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerHostName) {
		stringVal := plan.ServerHostName.ValueString()
		addRequest.ServerHostName = &stringVal
	}
	if internaltypes.IsDefined(plan.ServerPort) {
		intVal := int32(plan.ServerPort.ValueInt64())
		addRequest.ServerPort = &intVal
	}
	if internaltypes.IsDefined(plan.SyslogFacility) {
		intVal := int32(plan.SyslogFacility.ValueInt64())
		addRequest.SyslogFacility = &intVal
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		boolVal := plan.AutoFlush.ValueBool()
		addRequest.AutoFlush = &boolVal
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		boolVal := plan.Asynchronous.ValueBool()
		addRequest.Asynchronous = &boolVal
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		intVal := int32(plan.QueueSize.ValueInt64())
		addRequest.QueueSize = &intVal
	}
	if internaltypes.IsDefined(plan.DefaultSeverity) {
		var slice []string
		plan.DefaultSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.OverrideSeverity) {
		var slice []string
		plan.OverrideSeverity.ElementsAs(ctx, &slice, false)
		addRequest.OverrideSeverity = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Read a SyslogBasedErrorLogPublisherResponse object into the model struct
func readSyslogBasedErrorLogPublisherResponse(ctx context.Context, r *client.SyslogBasedErrorLogPublisherResponse, state *syslogBasedErrorLogPublisherResourceModel, expectedValues *syslogBasedErrorLogPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(int64(r.ServerPort))
	state.SyslogFacility = types.Int64Value(int64(r.SyslogFacility))
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSyslogBasedErrorLogPublisherOperations(plan syslogBasedErrorLogPublisherResourceModel, state syslogBasedErrorLogPublisherResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerHostName, state.ServerHostName, "server-host-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerPort, state.ServerPort, "server-port")
	operations.AddInt64OperationIfNecessary(&ops, plan.SyslogFacility, state.SyslogFacility, "syslog-facility")
	operations.AddBoolOperationIfNecessary(&ops, plan.AutoFlush, state.AutoFlush, "auto-flush")
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddInt64OperationIfNecessary(&ops, plan.QueueSize, state.QueueSize, "queue-size")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultSeverity, state.DefaultSeverity, "default-severity")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OverrideSeverity, state.OverrideSeverity, "override-severity")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	return ops
}

// Create a new resource
func (r *syslogBasedErrorLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan syslogBasedErrorLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddSyslogBasedErrorLogPublisherRequest(plan.Id.ValueString(),
		[]client.EnumsyslogBasedErrorLogPublisherSchemaUrn{client.ENUMSYSLOGBASEDERRORLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERSYSLOG_BASED_ERROR},
		plan.Enabled.ValueBool())
	err := addOptionalSyslogBasedErrorLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Syslog Based Error Log Publisher", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddSyslogBasedErrorLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Syslog Based Error Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state syslogBasedErrorLogPublisherResourceModel
	readSyslogBasedErrorLogPublisherResponse(ctx, addResponse.SyslogBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSyslogBasedErrorLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan syslogBasedErrorLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Syslog Based Error Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state syslogBasedErrorLogPublisherResourceModel
	readSyslogBasedErrorLogPublisherResponse(ctx, readResponse.SyslogBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogPublisherApi.UpdateLogPublisher(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSyslogBasedErrorLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Syslog Based Error Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSyslogBasedErrorLogPublisherResponse(ctx, updateResponse.SyslogBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *syslogBasedErrorLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSyslogBasedErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSyslogBasedErrorLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSyslogBasedErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSyslogBasedErrorLogPublisher(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state syslogBasedErrorLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Syslog Based Error Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSyslogBasedErrorLogPublisherResponse(ctx, readResponse.SyslogBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *syslogBasedErrorLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSyslogBasedErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSyslogBasedErrorLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSyslogBasedErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSyslogBasedErrorLogPublisher(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan syslogBasedErrorLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state syslogBasedErrorLogPublisherResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogPublisherApi.UpdateLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSyslogBasedErrorLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Syslog Based Error Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSyslogBasedErrorLogPublisherResponse(ctx, updateResponse.SyslogBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSyslogBasedErrorLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *syslogBasedErrorLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state syslogBasedErrorLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogPublisherApi.DeleteLogPublisherExecute(r.apiClient.LogPublisherApi.DeleteLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Syslog Based Error Log Publisher", err, httpResp)
		return
	}
}

func (r *syslogBasedErrorLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSyslogBasedErrorLogPublisher(ctx, req, resp)
}

func (r *defaultSyslogBasedErrorLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSyslogBasedErrorLogPublisher(ctx, req, resp)
}

func importSyslogBasedErrorLogPublisher(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
