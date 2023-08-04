package failurelockoutaction

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	_ resource.Resource                = &failureLockoutActionResource{}
	_ resource.ResourceWithConfigure   = &failureLockoutActionResource{}
	_ resource.ResourceWithImportState = &failureLockoutActionResource{}
	_ resource.Resource                = &defaultFailureLockoutActionResource{}
	_ resource.ResourceWithConfigure   = &defaultFailureLockoutActionResource{}
	_ resource.ResourceWithImportState = &defaultFailureLockoutActionResource{}
)

// Create a Failure Lockout Action resource
func NewFailureLockoutActionResource() resource.Resource {
	return &failureLockoutActionResource{}
}

func NewDefaultFailureLockoutActionResource() resource.Resource {
	return &defaultFailureLockoutActionResource{}
}

// failureLockoutActionResource is the resource implementation.
type failureLockoutActionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultFailureLockoutActionResource is the resource implementation.
type defaultFailureLockoutActionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *failureLockoutActionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_failure_lockout_action"
}

func (r *defaultFailureLockoutActionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_failure_lockout_action"
}

// Configure adds the provider configured client to the resource.
func (r *failureLockoutActionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultFailureLockoutActionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type failureLockoutActionResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	LastUpdated                       types.String `tfsdk:"last_updated"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
	Type                              types.String `tfsdk:"type"`
	Delay                             types.String `tfsdk:"delay"`
	AllowBlockingDelay                types.Bool   `tfsdk:"allow_blocking_delay"`
	GenerateAccountStatusNotification types.Bool   `tfsdk:"generate_account_status_notification"`
	Description                       types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *failureLockoutActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	failureLockoutActionSchema(ctx, req, resp, false)
}

func (r *defaultFailureLockoutActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	failureLockoutActionSchema(ctx, req, resp, true)
}

func failureLockoutActionSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Failure Lockout Action.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Failure Lockout Action resource. Options are ['delay-bind-response', 'no-operation', 'lock-account']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"delay-bind-response", "no-operation", "lock-account"}...),
				},
			},
			"delay": schema.StringAttribute{
				Description: "The length of time to delay the bind response for accounts with too many failed authentication attempts.",
				Optional:    true,
			},
			"allow_blocking_delay": schema.BoolAttribute{
				Description: "Indicates whether to delay the response for authentication attempts even if that delay may block the thread being used to process the attempt.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"generate_account_status_notification": schema.BoolAttribute{
				Description: "Indicates whether to generate an account status notification for cases in which a bind response is delayed because of failure lockout.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Failure Lockout Action",
				Optional:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{}
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputed(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *failureLockoutActionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFailureLockoutActionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanFailureLockoutAction(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model failureLockoutActionResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.GenerateAccountStatusNotification) && model.Type.ValueString() != "delay-bind-response" && model.Type.ValueString() != "no-operation" {
		resp.Diagnostics.AddError("Attribute 'generate_account_status_notification' not supported by pingdirectory_failure_lockout_action resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'generate_account_status_notification', the 'type' attribute must be one of ['delay-bind-response', 'no-operation']")
	}
	if internaltypes.IsDefined(model.AllowBlockingDelay) && model.Type.ValueString() != "delay-bind-response" {
		resp.Diagnostics.AddError("Attribute 'allow_blocking_delay' not supported by pingdirectory_failure_lockout_action resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'allow_blocking_delay', the 'type' attribute must be one of ['delay-bind-response']")
	}
	if internaltypes.IsDefined(model.Delay) && model.Type.ValueString() != "delay-bind-response" {
		resp.Diagnostics.AddError("Attribute 'delay' not supported by pingdirectory_failure_lockout_action resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'delay', the 'type' attribute must be one of ['delay-bind-response']")
	}
}

// Add optional fields to create request for delay-bind-response failure-lockout-action
func addOptionalDelayBindResponseFailureLockoutActionFields(ctx context.Context, addRequest *client.AddDelayBindResponseFailureLockoutActionRequest, plan failureLockoutActionResourceModel) {
	if internaltypes.IsDefined(plan.AllowBlockingDelay) {
		addRequest.AllowBlockingDelay = plan.AllowBlockingDelay.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.GenerateAccountStatusNotification) {
		addRequest.GenerateAccountStatusNotification = plan.GenerateAccountStatusNotification.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for no-operation failure-lockout-action
func addOptionalNoOperationFailureLockoutActionFields(ctx context.Context, addRequest *client.AddNoOperationFailureLockoutActionRequest, plan failureLockoutActionResourceModel) {
	if internaltypes.IsDefined(plan.GenerateAccountStatusNotification) {
		addRequest.GenerateAccountStatusNotification = plan.GenerateAccountStatusNotification.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for lock-account failure-lockout-action
func addOptionalLockAccountFailureLockoutActionFields(ctx context.Context, addRequest *client.AddLockAccountFailureLockoutActionRequest, plan failureLockoutActionResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a DelayBindResponseFailureLockoutActionResponse object into the model struct
func readDelayBindResponseFailureLockoutActionResponse(ctx context.Context, r *client.DelayBindResponseFailureLockoutActionResponse, state *failureLockoutActionResourceModel, expectedValues *failureLockoutActionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("delay-bind-response")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Delay = types.StringValue(r.Delay)
	config.CheckMismatchedPDFormattedAttributes("delay",
		expectedValues.Delay, state.Delay, diagnostics)
	state.AllowBlockingDelay = internaltypes.BoolTypeOrNil(r.AllowBlockingDelay)
	state.GenerateAccountStatusNotification = internaltypes.BoolTypeOrNil(r.GenerateAccountStatusNotification)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a NoOperationFailureLockoutActionResponse object into the model struct
func readNoOperationFailureLockoutActionResponse(ctx context.Context, r *client.NoOperationFailureLockoutActionResponse, state *failureLockoutActionResourceModel, expectedValues *failureLockoutActionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("no-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.GenerateAccountStatusNotification = internaltypes.BoolTypeOrNil(r.GenerateAccountStatusNotification)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a LockAccountFailureLockoutActionResponse object into the model struct
func readLockAccountFailureLockoutActionResponse(ctx context.Context, r *client.LockAccountFailureLockoutActionResponse, state *failureLockoutActionResourceModel, expectedValues *failureLockoutActionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("lock-account")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createFailureLockoutActionOperations(plan failureLockoutActionResourceModel, state failureLockoutActionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Delay, state.Delay, "delay")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowBlockingDelay, state.AllowBlockingDelay, "allow-blocking-delay")
	operations.AddBoolOperationIfNecessary(&ops, plan.GenerateAccountStatusNotification, state.GenerateAccountStatusNotification, "generate-account-status-notification")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a delay-bind-response failure-lockout-action
func (r *failureLockoutActionResource) CreateDelayBindResponseFailureLockoutAction(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan failureLockoutActionResourceModel) (*failureLockoutActionResourceModel, error) {
	addRequest := client.NewAddDelayBindResponseFailureLockoutActionRequest(plan.Name.ValueString(),
		[]client.EnumdelayBindResponseFailureLockoutActionSchemaUrn{client.ENUMDELAYBINDRESPONSEFAILURELOCKOUTACTIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0FAILURE_LOCKOUT_ACTIONDELAY_BIND_RESPONSE},
		plan.Delay.ValueString())
	addOptionalDelayBindResponseFailureLockoutActionFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.FailureLockoutActionApi.AddFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddFailureLockoutActionRequest(
		client.AddDelayBindResponseFailureLockoutActionRequestAsAddFailureLockoutActionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.AddFailureLockoutActionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Failure Lockout Action", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state failureLockoutActionResourceModel
	readDelayBindResponseFailureLockoutActionResponse(ctx, addResponse.DelayBindResponseFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a no-operation failure-lockout-action
func (r *failureLockoutActionResource) CreateNoOperationFailureLockoutAction(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan failureLockoutActionResourceModel) (*failureLockoutActionResourceModel, error) {
	addRequest := client.NewAddNoOperationFailureLockoutActionRequest(plan.Name.ValueString(),
		[]client.EnumnoOperationFailureLockoutActionSchemaUrn{client.ENUMNOOPERATIONFAILURELOCKOUTACTIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0FAILURE_LOCKOUT_ACTIONNO_OPERATION})
	addOptionalNoOperationFailureLockoutActionFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.FailureLockoutActionApi.AddFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddFailureLockoutActionRequest(
		client.AddNoOperationFailureLockoutActionRequestAsAddFailureLockoutActionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.AddFailureLockoutActionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Failure Lockout Action", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state failureLockoutActionResourceModel
	readNoOperationFailureLockoutActionResponse(ctx, addResponse.NoOperationFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a lock-account failure-lockout-action
func (r *failureLockoutActionResource) CreateLockAccountFailureLockoutAction(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan failureLockoutActionResourceModel) (*failureLockoutActionResourceModel, error) {
	addRequest := client.NewAddLockAccountFailureLockoutActionRequest(plan.Name.ValueString(),
		[]client.EnumlockAccountFailureLockoutActionSchemaUrn{client.ENUMLOCKACCOUNTFAILURELOCKOUTACTIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0FAILURE_LOCKOUT_ACTIONLOCK_ACCOUNT})
	addOptionalLockAccountFailureLockoutActionFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.FailureLockoutActionApi.AddFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddFailureLockoutActionRequest(
		client.AddLockAccountFailureLockoutActionRequestAsAddFailureLockoutActionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.AddFailureLockoutActionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Failure Lockout Action", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state failureLockoutActionResourceModel
	readLockAccountFailureLockoutActionResponse(ctx, addResponse.LockAccountFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *failureLockoutActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan failureLockoutActionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *failureLockoutActionResourceModel
	var err error
	if plan.Type.ValueString() == "delay-bind-response" {
		state, err = r.CreateDelayBindResponseFailureLockoutAction(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "no-operation" {
		state, err = r.CreateNoOperationFailureLockoutAction(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "lock-account" {
		state, err = r.CreateLockAccountFailureLockoutAction(ctx, req, resp, plan)
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
func (r *defaultFailureLockoutActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan failureLockoutActionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.GetFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Failure Lockout Action", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state failureLockoutActionResourceModel
	if readResponse.DelayBindResponseFailureLockoutActionResponse != nil {
		readDelayBindResponseFailureLockoutActionResponse(ctx, readResponse.DelayBindResponseFailureLockoutActionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NoOperationFailureLockoutActionResponse != nil {
		readNoOperationFailureLockoutActionResponse(ctx, readResponse.NoOperationFailureLockoutActionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LockAccountFailureLockoutActionResponse != nil {
		readLockAccountFailureLockoutActionResponse(ctx, readResponse.LockAccountFailureLockoutActionResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.FailureLockoutActionApi.UpdateFailureLockoutAction(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createFailureLockoutActionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.UpdateFailureLockoutActionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Failure Lockout Action", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.DelayBindResponseFailureLockoutActionResponse != nil {
			readDelayBindResponseFailureLockoutActionResponse(ctx, updateResponse.DelayBindResponseFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NoOperationFailureLockoutActionResponse != nil {
			readNoOperationFailureLockoutActionResponse(ctx, updateResponse.NoOperationFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LockAccountFailureLockoutActionResponse != nil {
			readLockAccountFailureLockoutActionResponse(ctx, updateResponse.LockAccountFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *failureLockoutActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFailureLockoutActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readFailureLockoutAction(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state failureLockoutActionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.FailureLockoutActionApi.GetFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Failure Lockout Action", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.DelayBindResponseFailureLockoutActionResponse != nil {
		readDelayBindResponseFailureLockoutActionResponse(ctx, readResponse.DelayBindResponseFailureLockoutActionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NoOperationFailureLockoutActionResponse != nil {
		readNoOperationFailureLockoutActionResponse(ctx, readResponse.NoOperationFailureLockoutActionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LockAccountFailureLockoutActionResponse != nil {
		readLockAccountFailureLockoutActionResponse(ctx, readResponse.LockAccountFailureLockoutActionResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *failureLockoutActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFailureLockoutActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateFailureLockoutAction(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan failureLockoutActionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state failureLockoutActionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.FailureLockoutActionApi.UpdateFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createFailureLockoutActionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.FailureLockoutActionApi.UpdateFailureLockoutActionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Failure Lockout Action", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.DelayBindResponseFailureLockoutActionResponse != nil {
			readDelayBindResponseFailureLockoutActionResponse(ctx, updateResponse.DelayBindResponseFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NoOperationFailureLockoutActionResponse != nil {
			readNoOperationFailureLockoutActionResponse(ctx, updateResponse.NoOperationFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LockAccountFailureLockoutActionResponse != nil {
			readLockAccountFailureLockoutActionResponse(ctx, updateResponse.LockAccountFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultFailureLockoutActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *failureLockoutActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state failureLockoutActionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.FailureLockoutActionApi.DeleteFailureLockoutActionExecute(r.apiClient.FailureLockoutActionApi.DeleteFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Failure Lockout Action", err, httpResp)
		return
	}
}

func (r *failureLockoutActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFailureLockoutAction(ctx, req, resp)
}

func (r *defaultFailureLockoutActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFailureLockoutAction(ctx, req, resp)
}

func importFailureLockoutAction(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
