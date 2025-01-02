package logrotationpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
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
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultLogRotationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type logRotationPolicyResourceModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
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

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *logRotationPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var planModel logRotationPolicyResourceModel
	req.Plan.Get(ctx, &planModel)
	planModel.setNotApplicableAttrsNull()
	resp.Plan.Set(ctx, &planModel)
}

func (model *logRotationPolicyResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "time-limit" {
		model.TimeOfDay, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "never-rotate" {
		model.TimeOfDay, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "size-limit" {
		model.TimeOfDay, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsLogRotationPolicy() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("rotation_interval"),
			path.MatchRoot("type"),
			[]string{"time-limit"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("time_of_day"),
			path.MatchRoot("type"),
			[]string{"fixed-time"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("file_size_limit"),
			path.MatchRoot("type"),
			[]string{"size-limit"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"time-limit",
			[]path.Expression{path.MatchRoot("rotation_interval")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"fixed-time",
			[]path.Expression{path.MatchRoot("time_of_day")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"size-limit",
			[]path.Expression{path.MatchRoot("file_size_limit")},
		),
	}
}

// Add config validators
func (r logRotationPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsLogRotationPolicy()
}

// Add config validators
func (r defaultLogRotationPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsLogRotationPolicy()
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
func populateLogRotationPolicyUnknownValues(model *logRotationPolicyResourceModel) {
	if model.TimeOfDay.IsUnknown() || model.TimeOfDay.IsNull() {
		model.TimeOfDay, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *logRotationPolicyResourceModel) populateAllComputedStringAttributes() {
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.RotationInterval.IsUnknown() || model.RotationInterval.IsNull() {
		model.RotationInterval = types.StringValue("")
	}
	if model.FileSizeLimit.IsUnknown() || model.FileSizeLimit.IsNull() {
		model.FileSizeLimit = types.StringValue("")
	}
}

// Read a TimeLimitLogRotationPolicyResponse object into the model struct
func readTimeLimitLogRotationPolicyResponse(ctx context.Context, r *client.TimeLimitLogRotationPolicyResponse, state *logRotationPolicyResourceModel, expectedValues *logRotationPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("time-limit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RotationInterval = types.StringValue(r.RotationInterval)
	config.CheckMismatchedPDFormattedAttributes("rotation_interval",
		expectedValues.RotationInterval, state.RotationInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogRotationPolicyUnknownValues(state)
}

// Read a FixedTimeLogRotationPolicyResponse object into the model struct
func readFixedTimeLogRotationPolicyResponse(ctx context.Context, r *client.FixedTimeLogRotationPolicyResponse, state *logRotationPolicyResourceModel, expectedValues *logRotationPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("fixed-time")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.TimeOfDay = internaltypes.GetStringSet(r.TimeOfDay)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogRotationPolicyUnknownValues(state)
}

// Read a NeverRotateLogRotationPolicyResponse object into the model struct
func readNeverRotateLogRotationPolicyResponse(ctx context.Context, r *client.NeverRotateLogRotationPolicyResponse, state *logRotationPolicyResourceModel, expectedValues *logRotationPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("never-rotate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogRotationPolicyUnknownValues(state)
}

// Read a SizeLimitLogRotationPolicyResponse object into the model struct
func readSizeLimitLogRotationPolicyResponse(ctx context.Context, r *client.SizeLimitLogRotationPolicyResponse, state *logRotationPolicyResourceModel, expectedValues *logRotationPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("size-limit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.FileSizeLimit = types.StringValue(r.FileSizeLimit)
	config.CheckMismatchedPDFormattedAttributes("file_size_limit",
		expectedValues.FileSizeLimit, state.FileSizeLimit, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogRotationPolicyUnknownValues(state)
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
	addRequest := client.NewAddTimeLimitLogRotationPolicyRequest([]client.EnumtimeLimitLogRotationPolicySchemaUrn{client.ENUMTIMELIMITLOGROTATIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_ROTATION_POLICYTIME_LIMIT},
		plan.RotationInterval.ValueString(),
		plan.Name.ValueString())
	addOptionalTimeLimitLogRotationPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRotationPolicyAPI.AddLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRotationPolicyRequest(
		client.AddTimeLimitLogRotationPolicyRequestAsAddLogRotationPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRotationPolicyAPI.AddLogRotationPolicyExecute(apiAddRequest)
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
	addRequest := client.NewAddFixedTimeLogRotationPolicyRequest([]client.EnumfixedTimeLogRotationPolicySchemaUrn{client.ENUMFIXEDTIMELOGROTATIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_ROTATION_POLICYFIXED_TIME},
		TimeOfDaySlice,
		plan.Name.ValueString())
	addOptionalFixedTimeLogRotationPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRotationPolicyAPI.AddLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRotationPolicyRequest(
		client.AddFixedTimeLogRotationPolicyRequestAsAddLogRotationPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRotationPolicyAPI.AddLogRotationPolicyExecute(apiAddRequest)
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
	addRequest := client.NewAddNeverRotateLogRotationPolicyRequest([]client.EnumneverRotateLogRotationPolicySchemaUrn{client.ENUMNEVERROTATELOGROTATIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_ROTATION_POLICYNEVER_ROTATE},
		plan.Name.ValueString())
	addOptionalNeverRotateLogRotationPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRotationPolicyAPI.AddLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRotationPolicyRequest(
		client.AddNeverRotateLogRotationPolicyRequestAsAddLogRotationPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRotationPolicyAPI.AddLogRotationPolicyExecute(apiAddRequest)
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
	addRequest := client.NewAddSizeLimitLogRotationPolicyRequest([]client.EnumsizeLimitLogRotationPolicySchemaUrn{client.ENUMSIZELIMITLOGROTATIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_ROTATION_POLICYSIZE_LIMIT},
		plan.FileSizeLimit.ValueString(),
		plan.Name.ValueString())
	addOptionalSizeLimitLogRotationPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRotationPolicyAPI.AddLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRotationPolicyRequest(
		client.AddSizeLimitLogRotationPolicyRequestAsAddLogRotationPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRotationPolicyAPI.AddLogRotationPolicyExecute(apiAddRequest)
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

	readResponse, httpResp, err := r.apiClient.LogRotationPolicyAPI.GetLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
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

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogRotationPolicyAPI.UpdateLogRotationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createLogRotationPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogRotationPolicyAPI.UpdateLogRotationPolicyExecute(updateRequest)
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
		if updateResponse.TimeLimitLogRotationPolicyResponse != nil {
			readTimeLimitLogRotationPolicyResponse(ctx, updateResponse.TimeLimitLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FixedTimeLogRotationPolicyResponse != nil {
			readFixedTimeLogRotationPolicyResponse(ctx, updateResponse.FixedTimeLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NeverRotateLogRotationPolicyResponse != nil {
			readNeverRotateLogRotationPolicyResponse(ctx, updateResponse.NeverRotateLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SizeLimitLogRotationPolicyResponse != nil {
			readSizeLimitLogRotationPolicyResponse(ctx, updateResponse.SizeLimitLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *logRotationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultLogRotationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogRotationPolicy(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readLogRotationPolicy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state logRotationPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogRotationPolicyAPI.GetLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Log Rotation Policy", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Rotation Policy", err, httpResp)
		}
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

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
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
	updateRequest := apiClient.LogRotationPolicyAPI.UpdateLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createLogRotationPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogRotationPolicyAPI.UpdateLogRotationPolicyExecute(updateRequest)
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
		if updateResponse.TimeLimitLogRotationPolicyResponse != nil {
			readTimeLimitLogRotationPolicyResponse(ctx, updateResponse.TimeLimitLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FixedTimeLogRotationPolicyResponse != nil {
			readFixedTimeLogRotationPolicyResponse(ctx, updateResponse.FixedTimeLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NeverRotateLogRotationPolicyResponse != nil {
			readNeverRotateLogRotationPolicyResponse(ctx, updateResponse.NeverRotateLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SizeLimitLogRotationPolicyResponse != nil {
			readSizeLimitLogRotationPolicyResponse(ctx, updateResponse.SizeLimitLogRotationPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
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

	httpResp, err := r.apiClient.LogRotationPolicyAPI.DeleteLogRotationPolicyExecute(r.apiClient.LogRotationPolicyAPI.DeleteLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
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
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
