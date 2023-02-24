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
	_ resource.Resource                = &generateServerProfileRecurringTaskResource{}
	_ resource.ResourceWithConfigure   = &generateServerProfileRecurringTaskResource{}
	_ resource.ResourceWithImportState = &generateServerProfileRecurringTaskResource{}
)

// Create a Generate Server Profile Recurring Task resource
func NewGenerateServerProfileRecurringTaskResource() resource.Resource {
	return &generateServerProfileRecurringTaskResource{}
}

// generateServerProfileRecurringTaskResource is the resource implementation.
type generateServerProfileRecurringTaskResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *generateServerProfileRecurringTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_generate_server_profile_recurring_task"
}

// Configure adds the provider configured client to the resource.
func (r *generateServerProfileRecurringTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type generateServerProfileRecurringTaskResourceModel struct {
	Id                            types.String `tfsdk:"id"`
	LastUpdated                   types.String `tfsdk:"last_updated"`
	Notifications                 types.Set    `tfsdk:"notifications"`
	RequiredActions               types.Set    `tfsdk:"required_actions"`
	ProfileDirectory              types.String `tfsdk:"profile_directory"`
	IncludePath                   types.Set    `tfsdk:"include_path"`
	RetainPreviousProfileCount    types.Int64  `tfsdk:"retain_previous_profile_count"`
	RetainPreviousProfileAge      types.String `tfsdk:"retain_previous_profile_age"`
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
func (r *generateServerProfileRecurringTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Generate Server Profile Recurring Task.",
		Attributes: map[string]schema.Attribute{
			"profile_directory": schema.StringAttribute{
				Description: "The directory in which the generated server profiles will be placed. The files will be named with the pattern \"server-profile-{timestamp}.zip\", where \"{timestamp}\" represents the time that the profile was generated.",
				Required:    true,
			},
			"include_path": schema.SetAttribute{
				Description: "An optional set of additional paths to files within the instance root that should be included in the generated server profile. All paths must be within the instance root, and relative paths will be relative to the instance root.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"retain_previous_profile_count": schema.Int64Attribute{
				Description: "The minimum number of previous server profile zip files that should be preserved after a new profile is generated.",
				Optional:    true,
				Computed:    true,
			},
			"retain_previous_profile_age": schema.StringAttribute{
				Description: "The minimum age of previous server profile zip files that should be preserved after a new profile is generated.",
				Optional:    true,
				Computed:    true,
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
func addOptionalGenerateServerProfileRecurringTaskFields(ctx context.Context, addRequest *client.AddGenerateServerProfileRecurringTaskRequest, plan generateServerProfileRecurringTaskResourceModel) {
	if internaltypes.IsDefined(plan.IncludePath) {
		var slice []string
		plan.IncludePath.ElementsAs(ctx, &slice, false)
		addRequest.IncludePath = slice
	}
	if internaltypes.IsDefined(plan.RetainPreviousProfileCount) {
		intVal := int32(plan.RetainPreviousProfileCount.ValueInt64())
		addRequest.RetainPreviousProfileCount = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetainPreviousProfileAge) {
		stringVal := plan.RetainPreviousProfileAge.ValueString()
		addRequest.RetainPreviousProfileAge = &stringVal
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

// Read a GenerateServerProfileRecurringTaskResponse object into the model struct
func readGenerateServerProfileRecurringTaskResponse(ctx context.Context, r *client.GenerateServerProfileRecurringTaskResponse, state *generateServerProfileRecurringTaskResourceModel, expectedValues *generateServerProfileRecurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ProfileDirectory = types.StringValue(r.ProfileDirectory)
	state.IncludePath = internaltypes.GetStringSet(r.IncludePath)
	state.RetainPreviousProfileCount = internaltypes.Int64TypeOrNil(r.RetainPreviousProfileCount)
	state.RetainPreviousProfileAge = internaltypes.StringTypeOrNil(r.RetainPreviousProfileAge, internaltypes.IsEmptyString(expectedValues.RetainPreviousProfileAge))
	config.CheckMismatchedPDFormattedAttributes("retain_previous_profile_age",
		expectedValues.RetainPreviousProfileAge, state.RetainPreviousProfileAge, diagnostics)
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
func createGenerateServerProfileRecurringTaskOperations(plan generateServerProfileRecurringTaskResourceModel, state generateServerProfileRecurringTaskResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ProfileDirectory, state.ProfileDirectory, "profile-directory")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludePath, state.IncludePath, "include-path")
	operations.AddInt64OperationIfNecessary(&ops, plan.RetainPreviousProfileCount, state.RetainPreviousProfileCount, "retain-previous-profile-count")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainPreviousProfileAge, state.RetainPreviousProfileAge, "retain-previous-profile-age")
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
func (r *generateServerProfileRecurringTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan generateServerProfileRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddGenerateServerProfileRecurringTaskRequest(plan.Id.ValueString(),
		[]client.EnumgenerateServerProfileRecurringTaskSchemaUrn{client.ENUMGENERATESERVERPROFILERECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKGENERATE_SERVER_PROFILE},
		plan.ProfileDirectory.ValueString())
	addOptionalGenerateServerProfileRecurringTaskFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddGenerateServerProfileRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Generate Server Profile Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state generateServerProfileRecurringTaskResourceModel
	readGenerateServerProfileRecurringTaskResponse(ctx, addResponse.GenerateServerProfileRecurringTaskResponse, &state, &plan, &resp.Diagnostics)

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
func (r *generateServerProfileRecurringTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state generateServerProfileRecurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RecurringTaskApi.GetRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Generate Server Profile Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readGenerateServerProfileRecurringTaskResponse(ctx, readResponse.GenerateServerProfileRecurringTaskResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *generateServerProfileRecurringTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan generateServerProfileRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state generateServerProfileRecurringTaskResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.RecurringTaskApi.UpdateRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createGenerateServerProfileRecurringTaskOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.RecurringTaskApi.UpdateRecurringTaskExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Generate Server Profile Recurring Task", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGenerateServerProfileRecurringTaskResponse(ctx, updateResponse.GenerateServerProfileRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
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
func (r *generateServerProfileRecurringTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state generateServerProfileRecurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RecurringTaskApi.DeleteRecurringTaskExecute(r.apiClient.RecurringTaskApi.DeleteRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Generate Server Profile Recurring Task", err, httpResp)
		return
	}
}

func (r *generateServerProfileRecurringTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
