package logrotationpolicy

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
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
	_ resource.Resource                = &logRotationPolicyResource{}
	_ resource.ResourceWithConfigure   = &logRotationPolicyResource{}
	_ resource.ResourceWithImportState = &logRotationPolicyResource{}
	_ resource.Resource                = &defaultLogRotationPolicyResource{}
	_ resource.ResourceWithConfigure   = &defaultLogRotationPolicyResource{}
	_ resource.ResourceWithImportState = &defaultLogRotationPolicyResource{}
)

// Create a Log Rotation Policy resource
func NewLogRotationPolicyResource() resource.Resource {
	return &logRotationPolicyResource{}
}

func NewDefaultLogRotationPolicyResource() resource.Resource {
	return &defaultLogRotationPolicyResource{}
}

// logRotationPolicyResource is the resource implementation.
type logRotationPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLogRotationPolicyResource is the resource implementation.
type defaultLogRotationPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *logRotationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_rotation_policy"
}

func (r *defaultLogRotationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_log_rotation_policy"
}

// Configure adds the provider configured client to the resource.
func (r *logRotationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultLogRotationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type logRotationPolicyResourceModel struct {
	Id               types.String `tfsdk:"id"`
	LastUpdated      types.String `tfsdk:"last_updated"`
	Notifications    types.Set    `tfsdk:"notifications"`
	RequiredActions  types.Set    `tfsdk:"required_actions"`
	Type             types.String `tfsdk:"type"`
	FileSizeLimit    types.String `tfsdk:"file_size_limit"`
	TimeOfDay        types.Set    `tfsdk:"time_of_day"`
	RotationInterval types.String `tfsdk:"rotation_interval"`
	Description      types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *logRotationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	logRotationPolicySchema(ctx, req, resp, false)
}

func (r *defaultLogRotationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	logRotationPolicySchema(ctx, req, resp, true)
}

func logRotationPolicySchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Log Rotation Policy.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Log Rotation Policy resource. Options are ['time-limit', 'fixed-time', 'never-rotate', 'size-limit']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"time-limit", "fixed-time", "never-rotate", "size-limit"}...),
				},
			},
			"file_size_limit": schema.StringAttribute{
				Description: "Specifies the maximum size that a log file can reach before it is rotated.",
				Optional:    true,
			},
			"time_of_day": schema.SetAttribute{
				Description: "Specifies the time of day at which log rotation should occur.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"rotation_interval": schema.StringAttribute{
				Description: "Specifies the time interval between rotations.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Rotation Policy",
				Optional:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"time-limit", "fixed-time", "never-rotate", "size-limit"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *logRotationPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLogRotationPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanLogRotationPolicy(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model logRotationPolicyResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.RotationInterval) && model.Type.ValueString() != "time-limit" {
		resp.Diagnostics.AddError("Attribute 'rotation_interval' not supported by pingdirectory_log_rotation_policy resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'rotation_interval', the 'type' attribute must be one of ['time-limit']")
	}
	if internaltypes.IsDefined(model.FileSizeLimit) && model.Type.ValueString() != "size-limit" {
		resp.Diagnostics.AddError("Attribute 'file_size_limit' not supported by pingdirectory_log_rotation_policy resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'file_size_limit', the 'type' attribute must be one of ['size-limit']")
	}
	if internaltypes.IsDefined(model.TimeOfDay) && model.Type.ValueString() != "fixed-time" {
		resp.Diagnostics.AddError("Attribute 'time_of_day' not supported by pingdirectory_log_rotation_policy resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'time_of_day', the 'type' attribute must be one of ['fixed-time']")
	}
}

// Add optional fields to create request for time-limit log-rotation-policy
func addOptionalTimeLimitLogRotationPolicyFields(ctx context.Context, addRequest *client.AddTimeLimitLogRotationPolicyRequest, plan logRotationPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for fixed-time log-rotation-policy
func addOptionalFixedTimeLogRotationPolicyFields(ctx context.Context, addRequest *client.AddFixedTimeLogRotationPolicyRequest, plan logRotationPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for never-rotate log-rotation-policy
func addOptionalNeverRotateLogRotationPolicyFields(ctx context.Context, addRequest *client.AddNeverRotateLogRotationPolicyRequest, plan logRotationPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for size-limit log-rotation-policy
func addOptionalSizeLimitLogRotationPolicyFields(ctx context.Context, addRequest *client.AddSizeLimitLogRotationPolicyRequest, plan logRotationPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateLogRotationPolicyUnknownValues(ctx context.Context, model *logRotationPolicyResourceModel) {
	if model.TimeOfDay.ElementType(ctx) == nil {
		model.TimeOfDay = types.SetNull(types.StringType)
	}
}

// Read a TimeLimitLogRotationPolicyResponse object into the model struct
func readTimeLimitLogRotationPolicyResponse(ctx context.Context, r *client.TimeLimitLogRotationPolicyResponse, state *logRotationPolicyResourceModel, expectedValues *logRotationPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("time-limit")
	state.Id = types.StringValue(r.Id)
	state.RotationInterval = types.StringValue(r.RotationInterval)
	config.CheckMismatchedPDFormattedAttributes("rotation_interval",
		expectedValues.RotationInterval, state.RotationInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogRotationPolicyUnknownValues(ctx, state)
}

// Read a FixedTimeLogRotationPolicyResponse object into the model struct
func readFixedTimeLogRotationPolicyResponse(ctx context.Context, r *client.FixedTimeLogRotationPolicyResponse, state *logRotationPolicyResourceModel, expectedValues *logRotationPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("fixed-time")
	state.Id = types.StringValue(r.Id)
	state.TimeOfDay = internaltypes.GetStringSet(r.TimeOfDay)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogRotationPolicyUnknownValues(ctx, state)
}

// Read a NeverRotateLogRotationPolicyResponse object into the model struct
func readNeverRotateLogRotationPolicyResponse(ctx context.Context, r *client.NeverRotateLogRotationPolicyResponse, state *logRotationPolicyResourceModel, expectedValues *logRotationPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("never-rotate")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogRotationPolicyUnknownValues(ctx, state)
}

// Read a SizeLimitLogRotationPolicyResponse object into the model struct
func readSizeLimitLogRotationPolicyResponse(ctx context.Context, r *client.SizeLimitLogRotationPolicyResponse, state *logRotationPolicyResourceModel, expectedValues *logRotationPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("size-limit")
	state.Id = types.StringValue(r.Id)
	state.FileSizeLimit = types.StringValue(r.FileSizeLimit)
	config.CheckMismatchedPDFormattedAttributes("file_size_limit",
		expectedValues.FileSizeLimit, state.FileSizeLimit, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogRotationPolicyUnknownValues(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createLogRotationPolicyOperations(plan logRotationPolicyResourceModel, state logRotationPolicyResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.FileSizeLimit, state.FileSizeLimit, "file-size-limit")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.TimeOfDay, state.TimeOfDay, "time-of-day")
	operations.AddStringOperationIfNecessary(&ops, plan.RotationInterval, state.RotationInterval, "rotation-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a time-limit log-rotation-policy
func (r *logRotationPolicyResource) CreateTimeLimitLogRotationPolicy(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logRotationPolicyResourceModel) (*logRotationPolicyResourceModel, error) {
	addRequest := client.NewAddTimeLimitLogRotationPolicyRequest(plan.Id.ValueString(),
		[]client.EnumtimeLimitLogRotationPolicySchemaUrn{client.ENUMTIMELIMITLOGROTATIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_ROTATION_POLICYTIME_LIMIT},
		plan.RotationInterval.ValueString())
	addOptionalTimeLimitLogRotationPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRotationPolicyApi.AddLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRotationPolicyRequest(
		client.AddTimeLimitLogRotationPolicyRequestAsAddLogRotationPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRotationPolicyApi.AddLogRotationPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Rotation Policy", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logRotationPolicyResourceModel
	readTimeLimitLogRotationPolicyResponse(ctx, addResponse.TimeLimitLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a fixed-time log-rotation-policy
func (r *logRotationPolicyResource) CreateFixedTimeLogRotationPolicy(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logRotationPolicyResourceModel) (*logRotationPolicyResourceModel, error) {
	var TimeOfDaySlice []string
	plan.TimeOfDay.ElementsAs(ctx, &TimeOfDaySlice, false)
	addRequest := client.NewAddFixedTimeLogRotationPolicyRequest(plan.Id.ValueString(),
		[]client.EnumfixedTimeLogRotationPolicySchemaUrn{client.ENUMFIXEDTIMELOGROTATIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_ROTATION_POLICYFIXED_TIME},
		TimeOfDaySlice)
	addOptionalFixedTimeLogRotationPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRotationPolicyApi.AddLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRotationPolicyRequest(
		client.AddFixedTimeLogRotationPolicyRequestAsAddLogRotationPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRotationPolicyApi.AddLogRotationPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Rotation Policy", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logRotationPolicyResourceModel
	readFixedTimeLogRotationPolicyResponse(ctx, addResponse.FixedTimeLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a never-rotate log-rotation-policy
func (r *logRotationPolicyResource) CreateNeverRotateLogRotationPolicy(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logRotationPolicyResourceModel) (*logRotationPolicyResourceModel, error) {
	addRequest := client.NewAddNeverRotateLogRotationPolicyRequest(plan.Id.ValueString(),
		[]client.EnumneverRotateLogRotationPolicySchemaUrn{client.ENUMNEVERROTATELOGROTATIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_ROTATION_POLICYNEVER_ROTATE})
	addOptionalNeverRotateLogRotationPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRotationPolicyApi.AddLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRotationPolicyRequest(
		client.AddNeverRotateLogRotationPolicyRequestAsAddLogRotationPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRotationPolicyApi.AddLogRotationPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Rotation Policy", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logRotationPolicyResourceModel
	readNeverRotateLogRotationPolicyResponse(ctx, addResponse.NeverRotateLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a size-limit log-rotation-policy
func (r *logRotationPolicyResource) CreateSizeLimitLogRotationPolicy(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logRotationPolicyResourceModel) (*logRotationPolicyResourceModel, error) {
	addRequest := client.NewAddSizeLimitLogRotationPolicyRequest(plan.Id.ValueString(),
		[]client.EnumsizeLimitLogRotationPolicySchemaUrn{client.ENUMSIZELIMITLOGROTATIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_ROTATION_POLICYSIZE_LIMIT},
		plan.FileSizeLimit.ValueString())
	addOptionalSizeLimitLogRotationPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRotationPolicyApi.AddLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRotationPolicyRequest(
		client.AddSizeLimitLogRotationPolicyRequestAsAddLogRotationPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRotationPolicyApi.AddLogRotationPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Rotation Policy", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logRotationPolicyResourceModel
	readSizeLimitLogRotationPolicyResponse(ctx, addResponse.SizeLimitLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *logRotationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan logRotationPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *logRotationPolicyResourceModel
	var err error
	if plan.Type.ValueString() == "time-limit" {
		state, err = r.CreateTimeLimitLogRotationPolicy(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "fixed-time" {
		state, err = r.CreateFixedTimeLogRotationPolicy(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "never-rotate" {
		state, err = r.CreateNeverRotateLogRotationPolicy(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "size-limit" {
		state, err = r.CreateSizeLimitLogRotationPolicy(ctx, req, resp, plan)
		if err != nil {
			return
		}
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
func (r *defaultLogRotationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan logRotationPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogRotationPolicyApi.GetLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Rotation Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state logRotationPolicyResourceModel
	if plan.Type.ValueString() == "time-limit" {
		readTimeLimitLogRotationPolicyResponse(ctx, readResponse.TimeLimitLogRotationPolicyResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "fixed-time" {
		readFixedTimeLogRotationPolicyResponse(ctx, readResponse.FixedTimeLogRotationPolicyResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "never-rotate" {
		readNeverRotateLogRotationPolicyResponse(ctx, readResponse.NeverRotateLogRotationPolicyResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "size-limit" {
		readSizeLimitLogRotationPolicyResponse(ctx, readResponse.SizeLimitLogRotationPolicyResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogRotationPolicyApi.UpdateLogRotationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createLogRotationPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogRotationPolicyApi.UpdateLogRotationPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Log Rotation Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "time-limit" {
			readTimeLimitLogRotationPolicyResponse(ctx, updateResponse.TimeLimitLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "fixed-time" {
			readFixedTimeLogRotationPolicyResponse(ctx, updateResponse.FixedTimeLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "never-rotate" {
			readNeverRotateLogRotationPolicyResponse(ctx, updateResponse.NeverRotateLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "size-limit" {
			readSizeLimitLogRotationPolicyResponse(ctx, updateResponse.SizeLimitLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *logRotationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLogRotationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readLogRotationPolicy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state logRotationPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogRotationPolicyApi.GetLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Rotation Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.TimeLimitLogRotationPolicyResponse != nil {
		readTimeLimitLogRotationPolicyResponse(ctx, readResponse.TimeLimitLogRotationPolicyResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FixedTimeLogRotationPolicyResponse != nil {
		readFixedTimeLogRotationPolicyResponse(ctx, readResponse.FixedTimeLogRotationPolicyResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NeverRotateLogRotationPolicyResponse != nil {
		readNeverRotateLogRotationPolicyResponse(ctx, readResponse.NeverRotateLogRotationPolicyResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SizeLimitLogRotationPolicyResponse != nil {
		readSizeLimitLogRotationPolicyResponse(ctx, readResponse.SizeLimitLogRotationPolicyResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *logRotationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLogRotationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLogRotationPolicy(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan logRotationPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state logRotationPolicyResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogRotationPolicyApi.UpdateLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createLogRotationPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogRotationPolicyApi.UpdateLogRotationPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Log Rotation Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "time-limit" {
			readTimeLimitLogRotationPolicyResponse(ctx, updateResponse.TimeLimitLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "fixed-time" {
			readFixedTimeLogRotationPolicyResponse(ctx, updateResponse.FixedTimeLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "never-rotate" {
			readNeverRotateLogRotationPolicyResponse(ctx, updateResponse.NeverRotateLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "size-limit" {
			readSizeLimitLogRotationPolicyResponse(ctx, updateResponse.SizeLimitLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultLogRotationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *logRotationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state logRotationPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogRotationPolicyApi.DeleteLogRotationPolicyExecute(r.apiClient.LogRotationPolicyApi.DeleteLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Log Rotation Policy", err, httpResp)
		return
	}
}

func (r *logRotationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLogRotationPolicy(ctx, req, resp)
}

func (r *defaultLogRotationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLogRotationPolicy(ctx, req, resp)
}

func importLogRotationPolicy(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
