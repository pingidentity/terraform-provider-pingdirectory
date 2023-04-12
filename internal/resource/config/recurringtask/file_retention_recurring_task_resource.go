package recurringtask

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &fileRetentionRecurringTaskResource{}
	_ resource.ResourceWithConfigure   = &fileRetentionRecurringTaskResource{}
	_ resource.ResourceWithImportState = &fileRetentionRecurringTaskResource{}
	_ resource.Resource                = &defaultFileRetentionRecurringTaskResource{}
	_ resource.ResourceWithConfigure   = &defaultFileRetentionRecurringTaskResource{}
	_ resource.ResourceWithImportState = &defaultFileRetentionRecurringTaskResource{}
)

// Create a File Retention Recurring Task resource
func NewFileRetentionRecurringTaskResource() resource.Resource {
	return &fileRetentionRecurringTaskResource{}
}

func NewDefaultFileRetentionRecurringTaskResource() resource.Resource {
	return &defaultFileRetentionRecurringTaskResource{}
}

// fileRetentionRecurringTaskResource is the resource implementation.
type fileRetentionRecurringTaskResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultFileRetentionRecurringTaskResource is the resource implementation.
type defaultFileRetentionRecurringTaskResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *fileRetentionRecurringTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_retention_recurring_task"
}

func (r *defaultFileRetentionRecurringTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_file_retention_recurring_task"
}

// Configure adds the provider configured client to the resource.
func (r *fileRetentionRecurringTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultFileRetentionRecurringTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type fileRetentionRecurringTaskResourceModel struct {
	Id                            types.String `tfsdk:"id"`
	LastUpdated                   types.String `tfsdk:"last_updated"`
	Notifications                 types.Set    `tfsdk:"notifications"`
	RequiredActions               types.Set    `tfsdk:"required_actions"`
	TargetDirectory               types.String `tfsdk:"target_directory"`
	FilenamePattern               types.String `tfsdk:"filename_pattern"`
	TimestampFormat               types.String `tfsdk:"timestamp_format"`
	RetainFileCount               types.Int64  `tfsdk:"retain_file_count"`
	RetainFileAge                 types.String `tfsdk:"retain_file_age"`
	RetainAggregateFileSize       types.String `tfsdk:"retain_aggregate_file_size"`
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
func (r *fileRetentionRecurringTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fileRetentionRecurringTaskSchema(ctx, req, resp, false)
}

func (r *defaultFileRetentionRecurringTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fileRetentionRecurringTaskSchema(ctx, req, resp, true)
}

func fileRetentionRecurringTaskSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a File Retention Recurring Task.",
		Attributes: map[string]schema.Attribute{
			"target_directory": schema.StringAttribute{
				Description: "The path to the directory containing the files to examine. The directory must exist.",
				Required:    true,
			},
			"filename_pattern": schema.StringAttribute{
				Description: "A pattern that specifies the names of the files to examine. The pattern may contain zero or more asterisks as wildcards, where each wildcard matches zero or more characters. It may also contain at most one occurrence of the special string \"${timestamp}\", which will match a timestamp with the format specified using the timestamp-format property. All other characters in the pattern will be treated literally.",
				Required:    true,
			},
			"timestamp_format": schema.StringAttribute{
				Description: "The format to use for the timestamp represented by the \"${timestamp}\" token in the filename pattern.",
				Required:    true,
			},
			"retain_file_count": schema.Int64Attribute{
				Description: "The minimum number of files matching the pattern that will be retained.",
				Optional:    true,
			},
			"retain_file_age": schema.StringAttribute{
				Description: "The minimum age of files matching the pattern that will be retained.",
				Optional:    true,
			},
			"retain_aggregate_file_size": schema.StringAttribute{
				Description: "The minimum aggregate size of files that will be retained. The size should be specified as an integer followed by a unit that is one of \"b\" or \"bytes\", \"kb\" or \"kilobytes\", \"mb\" or \"megabytes\", \"gb\" or \"gigabytes\", or \"tb\" or \"terabytes\". For example, a value of \"1 gb\" indicates that at least one gigabyte of files should be retained.",
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
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"email_on_start": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent whenever an instance of this Recurring Task starts running. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"email_on_success": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent whenever an instance of this Recurring Task completes successfully. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"email_on_failure": schema.SetAttribute{
				Description: "The email addresses to which a message should be sent if an instance of this Recurring Task fails to complete successfully. If this option is used, then at least one smtp-server must be configured in the global configuration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"alert_on_start": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task starts running.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"alert_on_success": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task completes successfully.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"alert_on_failure": schema.BoolAttribute{
				Description: "Indicates whether the server should generate an administrative alert whenever an instance of this Recurring Task fails to complete successfully.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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

// Add config validators
func (r fileRetentionRecurringTaskResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("retain_file_count"),
			path.MatchRoot("retain_aggregate_file_size"),
			path.MatchRoot("retain_file_age"),
		),
	}
}

// Add optional fields to create request
func addOptionalFileRetentionRecurringTaskFields(ctx context.Context, addRequest *client.AddFileRetentionRecurringTaskRequest, plan fileRetentionRecurringTaskResourceModel) {
	if internaltypes.IsDefined(plan.RetainFileCount) {
		intVal := int32(plan.RetainFileCount.ValueInt64())
		addRequest.RetainFileCount = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetainFileAge) {
		addRequest.RetainFileAge = plan.RetainFileAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetainAggregateFileSize) {
		addRequest.RetainAggregateFileSize = plan.RetainAggregateFileSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CancelOnTaskDependencyFailure) {
		addRequest.CancelOnTaskDependencyFailure = plan.CancelOnTaskDependencyFailure.ValueBoolPointer()
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
		addRequest.AlertOnStart = plan.AlertOnStart.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnSuccess) {
		addRequest.AlertOnSuccess = plan.AlertOnSuccess.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlertOnFailure) {
		addRequest.AlertOnFailure = plan.AlertOnFailure.ValueBoolPointer()
	}
}

// Read a FileRetentionRecurringTaskResponse object into the model struct
func readFileRetentionRecurringTaskResponse(ctx context.Context, r *client.FileRetentionRecurringTaskResponse, state *fileRetentionRecurringTaskResourceModel, expectedValues *fileRetentionRecurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.TargetDirectory = types.StringValue(r.TargetDirectory)
	state.FilenamePattern = types.StringValue(r.FilenamePattern)
	state.TimestampFormat = types.StringValue(r.TimestampFormat.String())
	state.RetainFileCount = internaltypes.Int64TypeOrNil(r.RetainFileCount)
	state.RetainFileAge = internaltypes.StringTypeOrNil(r.RetainFileAge, internaltypes.IsEmptyString(expectedValues.RetainFileAge))
	config.CheckMismatchedPDFormattedAttributes("retain_file_age",
		expectedValues.RetainFileAge, state.RetainFileAge, diagnostics)
	state.RetainAggregateFileSize = internaltypes.StringTypeOrNil(r.RetainAggregateFileSize, internaltypes.IsEmptyString(expectedValues.RetainAggregateFileSize))
	config.CheckMismatchedPDFormattedAttributes("retain_aggregate_file_size",
		expectedValues.RetainAggregateFileSize, state.RetainAggregateFileSize, diagnostics)
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
func createFileRetentionRecurringTaskOperations(plan fileRetentionRecurringTaskResourceModel, state fileRetentionRecurringTaskResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.TargetDirectory, state.TargetDirectory, "target-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.FilenamePattern, state.FilenamePattern, "filename-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.TimestampFormat, state.TimestampFormat, "timestamp-format")
	operations.AddInt64OperationIfNecessary(&ops, plan.RetainFileCount, state.RetainFileCount, "retain-file-count")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainFileAge, state.RetainFileAge, "retain-file-age")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainAggregateFileSize, state.RetainAggregateFileSize, "retain-aggregate-file-size")
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
func (r *fileRetentionRecurringTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fileRetentionRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timestampFormat, err := client.NewEnumrecurringTaskTimestampFormatPropFromValue(plan.TimestampFormat.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for TimestampFormat", err.Error())
		return
	}
	addRequest := client.NewAddFileRetentionRecurringTaskRequest(plan.Id.ValueString(),
		[]client.EnumfileRetentionRecurringTaskSchemaUrn{client.ENUMFILERETENTIONRECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKFILE_RETENTION},
		plan.TargetDirectory.ValueString(),
		plan.FilenamePattern.ValueString(),
		*timestampFormat)
	addOptionalFileRetentionRecurringTaskFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddFileRetentionRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the File Retention Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state fileRetentionRecurringTaskResourceModel
	readFileRetentionRecurringTaskResponse(ctx, addResponse.FileRetentionRecurringTaskResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultFileRetentionRecurringTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fileRetentionRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RecurringTaskApi.GetRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the File Retention Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state fileRetentionRecurringTaskResourceModel
	readFileRetentionRecurringTaskResponse(ctx, readResponse.FileRetentionRecurringTaskResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.RecurringTaskApi.UpdateRecurringTask(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createFileRetentionRecurringTaskOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.RecurringTaskApi.UpdateRecurringTaskExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the File Retention Recurring Task", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFileRetentionRecurringTaskResponse(ctx, updateResponse.FileRetentionRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
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
func (r *fileRetentionRecurringTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFileRetentionRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFileRetentionRecurringTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFileRetentionRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readFileRetentionRecurringTask(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state fileRetentionRecurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.RecurringTaskApi.GetRecurringTask(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the File Retention Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readFileRetentionRecurringTaskResponse(ctx, readResponse.FileRetentionRecurringTaskResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *fileRetentionRecurringTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFileRetentionRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFileRetentionRecurringTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFileRetentionRecurringTask(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateFileRetentionRecurringTask(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan fileRetentionRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state fileRetentionRecurringTaskResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.RecurringTaskApi.UpdateRecurringTask(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createFileRetentionRecurringTaskOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.RecurringTaskApi.UpdateRecurringTaskExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the File Retention Recurring Task", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFileRetentionRecurringTaskResponse(ctx, updateResponse.FileRetentionRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultFileRetentionRecurringTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *fileRetentionRecurringTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state fileRetentionRecurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RecurringTaskApi.DeleteRecurringTaskExecute(r.apiClient.RecurringTaskApi.DeleteRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the File Retention Recurring Task", err, httpResp)
		return
	}
}

func (r *fileRetentionRecurringTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFileRetentionRecurringTask(ctx, req, resp)
}

func (r *defaultFileRetentionRecurringTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFileRetentionRecurringTask(ctx, req, resp)
}

func importFileRetentionRecurringTask(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
