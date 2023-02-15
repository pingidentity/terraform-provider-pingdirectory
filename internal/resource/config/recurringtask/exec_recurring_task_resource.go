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
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &execRecurringTaskResource{}
	_ resource.ResourceWithConfigure   = &execRecurringTaskResource{}
	_ resource.ResourceWithImportState = &execRecurringTaskResource{}
)

// Create a Exec Recurring Task resource
func NewExecRecurringTaskResource() resource.Resource {
	return &execRecurringTaskResource{}
}

// execRecurringTaskResource is the resource implementation.
type execRecurringTaskResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *execRecurringTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exec_recurring_task"
}

// Configure adds the provider configured client to the resource.
func (r *execRecurringTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type execRecurringTaskResourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	LastUpdated                           types.String `tfsdk:"last_updated"`
	Notifications                         types.Set    `tfsdk:"notifications"`
	RequiredActions                       types.Set    `tfsdk:"required_actions"`
	CommandPath                           types.String `tfsdk:"command_path"`
	CommandArguments                      types.String `tfsdk:"command_arguments"`
	CommandOutputFileBaseName             types.String `tfsdk:"command_output_file_base_name"`
	RetainPreviousOutputFileCount         types.Int64  `tfsdk:"retain_previous_output_file_count"`
	RetainPreviousOutputFileAge           types.String `tfsdk:"retain_previous_output_file_age"`
	LogCommandOutput                      types.Bool   `tfsdk:"log_command_output"`
	TaskCompletionStateForNonzeroExitCode types.String `tfsdk:"task_completion_state_for_nonzero_exit_code"`
	WorkingDirectory                      types.String `tfsdk:"working_directory"`
	Description                           types.String `tfsdk:"description"`
	CancelOnTaskDependencyFailure         types.Bool   `tfsdk:"cancel_on_task_dependency_failure"`
	EmailOnStart                          types.Set    `tfsdk:"email_on_start"`
	EmailOnSuccess                        types.Set    `tfsdk:"email_on_success"`
	EmailOnFailure                        types.Set    `tfsdk:"email_on_failure"`
	AlertOnStart                          types.Bool   `tfsdk:"alert_on_start"`
	AlertOnSuccess                        types.Bool   `tfsdk:"alert_on_success"`
	AlertOnFailure                        types.Bool   `tfsdk:"alert_on_failure"`
}

// GetSchema defines the schema for the resource.
func (r *execRecurringTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Exec Recurring Task.",
		Attributes: map[string]schema.Attribute{
			"command_path": schema.StringAttribute{
				Description: "The absolute path to the command to execute. It must be an absolute path, the corresponding file must exist, and it must be listed in the config/exec-command-whitelist.txt file.",
				Required:    true,
			},
			"command_arguments": schema.StringAttribute{
				Description: "A string containing the arguments to provide to the command. If the command should be run without arguments, this property should be left undefined. If there should be multiple arguments, then they should be separated with spaces.",
				Optional:    true,
			},
			"command_output_file_base_name": schema.StringAttribute{
				Description: "The path and base name for a file to which the command output (both standard output and standard error) should be written. This may be left undefined if the command output should not be recorded into a file.",
				Optional:    true,
			},
			"retain_previous_output_file_count": schema.Int64Attribute{
				Description: "The minimum number of previous command output files that should be preserved after a new instance of the command is invoked.",
				Optional:    true,
			},
			"retain_previous_output_file_age": schema.StringAttribute{
				Description: "The minimum age of previous command output files that should be preserved after a new instance of the command is invoked.",
				Optional:    true,
			},
			"log_command_output": schema.BoolAttribute{
				Description: "Indicates whether the command's output (both standard output and standard error) should be recorded in the server's error log.",
				Optional:    true,
				Computed:    true,
			},
			"task_completion_state_for_nonzero_exit_code": schema.StringAttribute{
				Description: "The final task state that a task instance should have if the task executes the specified command and that command completes with a nonzero exit code, which generally means that the command did not complete successfully.",
				Optional:    true,
				Computed:    true,
			},
			"working_directory": schema.StringAttribute{
				Description: "The absolute path to a working directory where the command should be executed. It must be an absolute path and the corresponding directory must exist.",
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
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalExecRecurringTaskFields(ctx context.Context, addRequest *client.AddExecRecurringTaskRequest, plan execRecurringTaskResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CommandArguments) {
		stringVal := plan.CommandArguments.ValueString()
		addRequest.CommandArguments = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CommandOutputFileBaseName) {
		stringVal := plan.CommandOutputFileBaseName.ValueString()
		addRequest.CommandOutputFileBaseName = &stringVal
	}
	if internaltypes.IsDefined(plan.RetainPreviousOutputFileCount) {
		intVal := int32(plan.RetainPreviousOutputFileCount.ValueInt64())
		addRequest.RetainPreviousOutputFileCount = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetainPreviousOutputFileAge) {
		stringVal := plan.RetainPreviousOutputFileAge.ValueString()
		addRequest.RetainPreviousOutputFileAge = &stringVal
	}
	if internaltypes.IsDefined(plan.LogCommandOutput) {
		boolVal := plan.LogCommandOutput.ValueBool()
		addRequest.LogCommandOutput = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TaskCompletionStateForNonzeroExitCode) {
		taskCompletionStateForNonzeroExitCode, err := client.NewEnumrecurringTaskTaskCompletionStateForNonzeroExitCodePropFromValue(plan.TaskCompletionStateForNonzeroExitCode.ValueString())
		if err != nil {
			return err
		}
		addRequest.TaskCompletionStateForNonzeroExitCode = taskCompletionStateForNonzeroExitCode
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.WorkingDirectory) {
		stringVal := plan.WorkingDirectory.ValueString()
		addRequest.WorkingDirectory = &stringVal
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
	return nil
}

// Read a ExecRecurringTaskResponse object into the model struct
func readExecRecurringTaskResponse(ctx context.Context, r *client.ExecRecurringTaskResponse, state *execRecurringTaskResourceModel, expectedValues *execRecurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.CommandPath = types.StringValue(r.CommandPath)
	state.CommandArguments = internaltypes.StringTypeOrNil(r.CommandArguments, internaltypes.IsEmptyString(expectedValues.CommandArguments))
	state.CommandOutputFileBaseName = internaltypes.StringTypeOrNil(r.CommandOutputFileBaseName, internaltypes.IsEmptyString(expectedValues.CommandOutputFileBaseName))
	state.RetainPreviousOutputFileCount = internaltypes.Int64TypeOrNil(r.RetainPreviousOutputFileCount)
	state.RetainPreviousOutputFileAge = internaltypes.StringTypeOrNil(r.RetainPreviousOutputFileAge, internaltypes.IsEmptyString(expectedValues.RetainPreviousOutputFileAge))
	config.CheckMismatchedPDFormattedAttributes("retain_previous_output_file_age",
		expectedValues.RetainPreviousOutputFileAge, state.RetainPreviousOutputFileAge, diagnostics)
	state.LogCommandOutput = internaltypes.BoolTypeOrNil(r.LogCommandOutput)
	state.TaskCompletionStateForNonzeroExitCode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrecurringTaskTaskCompletionStateForNonzeroExitCodeProp(r.TaskCompletionStateForNonzeroExitCode), internaltypes.IsEmptyString(expectedValues.TaskCompletionStateForNonzeroExitCode))
	state.WorkingDirectory = internaltypes.StringTypeOrNil(r.WorkingDirectory, internaltypes.IsEmptyString(expectedValues.WorkingDirectory))
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
func createExecRecurringTaskOperations(plan execRecurringTaskResourceModel, state execRecurringTaskResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.CommandPath, state.CommandPath, "command-path")
	operations.AddStringOperationIfNecessary(&ops, plan.CommandArguments, state.CommandArguments, "command-arguments")
	operations.AddStringOperationIfNecessary(&ops, plan.CommandOutputFileBaseName, state.CommandOutputFileBaseName, "command-output-file-base-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.RetainPreviousOutputFileCount, state.RetainPreviousOutputFileCount, "retain-previous-output-file-count")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainPreviousOutputFileAge, state.RetainPreviousOutputFileAge, "retain-previous-output-file-age")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogCommandOutput, state.LogCommandOutput, "log-command-output")
	operations.AddStringOperationIfNecessary(&ops, plan.TaskCompletionStateForNonzeroExitCode, state.TaskCompletionStateForNonzeroExitCode, "task-completion-state-for-nonzero-exit-code")
	operations.AddStringOperationIfNecessary(&ops, plan.WorkingDirectory, state.WorkingDirectory, "working-directory")
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
func (r *execRecurringTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan execRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddExecRecurringTaskRequest(plan.Id.ValueString(),
		[]client.EnumexecRecurringTaskSchemaUrn{client.ENUMEXECRECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKEXEC},
		plan.CommandPath.ValueString())
	err := addOptionalExecRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Exec Recurring Task", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddExecRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Exec Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state execRecurringTaskResourceModel
	readExecRecurringTaskResponse(ctx, addResponse.ExecRecurringTaskResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *execRecurringTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state execRecurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RecurringTaskApi.GetRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Exec Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readExecRecurringTaskResponse(ctx, readResponse.ExecRecurringTaskResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *execRecurringTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan execRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state execRecurringTaskResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.RecurringTaskApi.UpdateRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createExecRecurringTaskOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.RecurringTaskApi.UpdateRecurringTaskExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Exec Recurring Task", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readExecRecurringTaskResponse(ctx, updateResponse.ExecRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
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
func (r *execRecurringTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state execRecurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RecurringTaskApi.DeleteRecurringTaskExecute(r.apiClient.RecurringTaskApi.DeleteRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Exec Recurring Task", err, httpResp)
		return
	}
}

func (r *execRecurringTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
