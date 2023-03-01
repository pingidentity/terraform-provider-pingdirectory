package recurringtask

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
	_ resource.Resource                = &enterLockdownModeRecurringTaskResource{}
	_ resource.ResourceWithConfigure   = &enterLockdownModeRecurringTaskResource{}
	_ resource.ResourceWithImportState = &enterLockdownModeRecurringTaskResource{}
	_ resource.Resource                = &defaultEnterLockdownModeRecurringTaskResource{}
	_ resource.ResourceWithConfigure   = &defaultEnterLockdownModeRecurringTaskResource{}
	_ resource.ResourceWithImportState = &defaultEnterLockdownModeRecurringTaskResource{}
)

// Create a Enter Lockdown Mode Recurring Task resource
func NewEnterLockdownModeRecurringTaskResource() resource.Resource {
	return &enterLockdownModeRecurringTaskResource{}
}

func NewDefaultEnterLockdownModeRecurringTaskResource() resource.Resource {
	return &defaultEnterLockdownModeRecurringTaskResource{}
}

// enterLockdownModeRecurringTaskResource is the resource implementation.
type enterLockdownModeRecurringTaskResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultEnterLockdownModeRecurringTaskResource is the resource implementation.
type defaultEnterLockdownModeRecurringTaskResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *enterLockdownModeRecurringTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enter_lockdown_mode_recurring_task"
}

func (r *defaultEnterLockdownModeRecurringTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_enter_lockdown_mode_recurring_task"
}

// Configure adds the provider configured client to the resource.
func (r *enterLockdownModeRecurringTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultEnterLockdownModeRecurringTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type enterLockdownModeRecurringTaskResourceModel struct {
	Id                            types.String `tfsdk:"id"`
	LastUpdated                   types.String `tfsdk:"last_updated"`
	Notifications                 types.Set    `tfsdk:"notifications"`
	RequiredActions               types.Set    `tfsdk:"required_actions"`
	Reason                        types.String `tfsdk:"reason"`
	Description                   types.String `tfsdk:"description"`
	CancelOnTaskDependencyFailure types.Bool   `tfsdk:"cancel_on_task_dependency_failure"`
	EmailOnStart                  types.Set    `tfsdk:"email_on_start"`
	EmailOnSuccess                types.Set    `tfsdk:"email_on_success"`
	EmailOnFailure                types.Set    `tfsdk:"email_on_failure"`
	AlertOnStart                  types.Bool   `tfsdk:"alert_on_start"`
	AlertOnSuccess                types.Bool   `tfsdk:"alert_on_success"`
	AlertOnFailure                types.Bool   `tfsdk:"alert_on_failure"`
}

// GetSchema defines the schema for the resource.
func (r *enterLockdownModeRecurringTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	enterLockdownModeRecurringTaskSchema(ctx, req, resp, false)
}

func (r *defaultEnterLockdownModeRecurringTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	enterLockdownModeRecurringTaskSchema(ctx, req, resp, true)
}

func enterLockdownModeRecurringTaskSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Enter Lockdown Mode Recurring Task.",
		Attributes: map[string]schema.Attribute{
			"reason": schema.StringAttribute{
				Description: "The reason that the server is being placed in lockdown mode.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Recurring Task",
				Optional:    true,
			},
			"cancel_on_task_dependency_failure": schema.BoolAttribute{
				Description: "Indicates whether an instance of this Recurring Task should be canceled if the task immediately before it in the recurring task chain fails to complete successfully (including if it is canceled by an administrator before it starts or while it is running).",
				Optional:    true,
				Computed:    true,
			},
			"email_on_start": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent whenever an instance of this Recurring Task starts running. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"email_on_success": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent whenever an instance of this Recurring Task completes successfully. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"email_on_failure": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent if an instance of this Recurring Task fails to complete successfully. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"alert_on_start": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task starts running.",
				Optional:    true,
				Computed:    true,
			},
			"alert_on_success": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task completes successfully.",
				Optional:    true,
				Computed:    true,
			},
			"alert_on_failure": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task fails to complete successfully.",
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
func addOptionalEnterLockdownModeRecurringTaskFields(ctx context.Context, addRequest *client.AddEnterLockdownModeRecurringTaskRequest, plan enterLockdownModeRecurringTaskResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Reason) {
		stringVal := plan.Reason.ValueString()
		addRequest.Reason = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		boolVal := plan.CancelOnTaskDependencyFailure.ValueBool()
		addRequest.CancelOnTaskDependencyFailure = &boolVal
	}
	if internaltypes.IsDefined(plan.EmailOnStart) {
		var slice []string
		plan.EmailOnStart.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnStart = slice
	}
	if internaltypes.IsDefined(plan.EmailOnSuccess) {
		var slice []string
		plan.EmailOnSuccess.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnSuccess = slice
	}
	if internaltypes.IsDefined(plan.EmailOnFailure) {
		var slice []string
		plan.EmailOnFailure.ElementsAs(ctx, &slice, false)
		addRequest.EmailOnFailure = slice
	}
	if internaltypes.IsDefined(plan.AlertOnStart) {
		boolVal := plan.AlertOnStart.ValueBool()
		addRequest.AlertOnStart = &boolVal
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		boolVal := plan.AlertOnSuccess.ValueBool()
		addRequest.AlertOnSuccess = &boolVal
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		boolVal := plan.AlertOnFailure.ValueBool()
		addRequest.AlertOnFailure = &boolVal
	}
}

// Read a EnterLockdownModeRecurringTaskResponse object into the model struct
func readEnterLockdownModeRecurringTaskResponse(ctx context.Context, r *client.EnterLockdownModeRecurringTaskResponse, state *enterLockdownModeRecurringTaskResourceModel, expectedValues *enterLockdownModeRecurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Reason = internaltypes.StringTypeOrNil(r.Reason, internaltypes.IsEmptyString(expectedValues.Reason))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CancelOnTaskDependencyFailure = internaltypes.BoolTypeOrNil(r.CancelOnTaskDependencyFailure)
	state.EmailOnStart = internaltypes.GetStringSet(r.EmailOnStart)
	state.EmailOnSuccess = internaltypes.GetStringSet(r.EmailOnSuccess)
	state.EmailOnFailure = internaltypes.GetStringSet(r.EmailOnFailure)
	state.AlertOnStart = internaltypes.BoolTypeOrNil(r.AlertOnStart)
	state.AlertOnSuccess = internaltypes.BoolTypeOrNil(r.AlertOnSuccess)
	state.AlertOnFailure = internaltypes.BoolTypeOrNil(r.AlertOnFailure)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createEnterLockdownModeRecurringTaskOperations(plan enterLockdownModeRecurringTaskResourceModel, state enterLockdownModeRecurringTaskResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Reason, state.Reason, "reason")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.CancelOnTaskDependencyFailure, state.CancelOnTaskDependencyFailure, "cancel-on-task-dependency-failure")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EmailOnStart, state.EmailOnStart, "email-on-start")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EmailOnSuccess, state.EmailOnSuccess, "email-on-success")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EmailOnFailure, state.EmailOnFailure, "email-on-failure")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlertOnStart, state.AlertOnStart, "alert-on-start")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlertOnSuccess, state.AlertOnSuccess, "alert-on-success")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlertOnFailure, state.AlertOnFailure, "alert-on-failure")
	return ops
}

// Create a new resource
func (r *enterLockdownModeRecurringTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan enterLockdownModeRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddEnterLockdownModeRecurringTaskRequest(plan.Id.ValueString(),
		[]client.EnumenterLockdownModeRecurringTaskSchemaUrn{client.ENUMENTERLOCKDOWNMODERECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKENTER_LOCKDOWN_MODE})
	addOptionalEnterLockdownModeRecurringTaskFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddEnterLockdownModeRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Enter Lockdown Mode Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state enterLockdownModeRecurringTaskResourceModel
	readEnterLockdownModeRecurringTaskResponse(ctx, addResponse.EnterLockdownModeRecurringTaskResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultEnterLockdownModeRecurringTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan enterLockdownModeRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RecurringTaskApi.GetRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Enter Lockdown Mode Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state enterLockdownModeRecurringTaskResourceModel
	readEnterLockdownModeRecurringTaskResponse(ctx, readResponse.EnterLockdownModeRecurringTaskResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.RecurringTaskApi.UpdateRecurringTask(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createEnterLockdownModeRecurringTaskOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.RecurringTaskApi.UpdateRecurringTaskExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Enter Lockdown Mode Recurring Task", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readEnterLockdownModeRecurringTaskResponse(ctx, updateResponse.EnterLockdownModeRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
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
func (r *enterLockdownModeRecurringTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readEnterLockdownModeRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultEnterLockdownModeRecurringTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readEnterLockdownModeRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readEnterLockdownModeRecurringTask(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state enterLockdownModeRecurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.RecurringTaskApi.GetRecurringTask(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Enter Lockdown Mode Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readEnterLockdownModeRecurringTaskResponse(ctx, readResponse.EnterLockdownModeRecurringTaskResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *enterLockdownModeRecurringTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateEnterLockdownModeRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultEnterLockdownModeRecurringTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateEnterLockdownModeRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateEnterLockdownModeRecurringTask(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan enterLockdownModeRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state enterLockdownModeRecurringTaskResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.RecurringTaskApi.UpdateRecurringTask(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createEnterLockdownModeRecurringTaskOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.RecurringTaskApi.UpdateRecurringTaskExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Enter Lockdown Mode Recurring Task", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readEnterLockdownModeRecurringTaskResponse(ctx, updateResponse.EnterLockdownModeRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultEnterLockdownModeRecurringTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *enterLockdownModeRecurringTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state enterLockdownModeRecurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RecurringTaskApi.DeleteRecurringTaskExecute(r.apiClient.RecurringTaskApi.DeleteRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Enter Lockdown Mode Recurring Task", err, httpResp)
		return
	}
}

func (r *enterLockdownModeRecurringTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importEnterLockdownModeRecurringTask(ctx, req, resp)
}

func (r *defaultEnterLockdownModeRecurringTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importEnterLockdownModeRecurringTask(ctx, req, resp)
}

func importEnterLockdownModeRecurringTask(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
