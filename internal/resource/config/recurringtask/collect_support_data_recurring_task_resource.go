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
	_ resource.Resource                = &collectSupportDataRecurringTaskResource{}
	_ resource.ResourceWithConfigure   = &collectSupportDataRecurringTaskResource{}
	_ resource.ResourceWithImportState = &collectSupportDataRecurringTaskResource{}
)

// Create a Collect Support Data Recurring Task resource
func NewCollectSupportDataRecurringTaskResource() resource.Resource {
	return &collectSupportDataRecurringTaskResource{}
}

// collectSupportDataRecurringTaskResource is the resource implementation.
type collectSupportDataRecurringTaskResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *collectSupportDataRecurringTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collect_support_data_recurring_task"
}

// Configure adds the provider configured client to the resource.
func (r *collectSupportDataRecurringTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type collectSupportDataRecurringTaskResourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	LastUpdated                           types.String `tfsdk:"last_updated"`
	Notifications                         types.Set    `tfsdk:"notifications"`
	RequiredActions                       types.Set    `tfsdk:"required_actions"`
	OutputDirectory                       types.String `tfsdk:"output_directory"`
	EncryptionPassphraseFile              types.String `tfsdk:"encryption_passphrase_file"`
	IncludeExpensiveData                  types.Bool   `tfsdk:"include_expensive_data"`
	IncludeReplicationStateDump           types.Bool   `tfsdk:"include_replication_state_dump"`
	IncludeBinaryFiles                    types.Bool   `tfsdk:"include_binary_files"`
	IncludeExtensionSource                types.Bool   `tfsdk:"include_extension_source"`
	UseSequentialMode                     types.Bool   `tfsdk:"use_sequential_mode"`
	SecurityLevel                         types.String `tfsdk:"security_level"`
	JstackCount                           types.Int64  `tfsdk:"jstack_count"`
	ReportCount                           types.Int64  `tfsdk:"report_count"`
	ReportIntervalSeconds                 types.Int64  `tfsdk:"report_interval_seconds"`
	LogDuration                           types.String `tfsdk:"log_duration"`
	LogFileHeadCollectionSize             types.String `tfsdk:"log_file_head_collection_size"`
	LogFileTailCollectionSize             types.String `tfsdk:"log_file_tail_collection_size"`
	Comment                               types.String `tfsdk:"comment"`
	RetainPreviousSupportDataArchiveCount types.Int64  `tfsdk:"retain_previous_support_data_archive_count"`
	RetainPreviousSupportDataArchiveAge   types.String `tfsdk:"retain_previous_support_data_archive_age"`
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
func (r *collectSupportDataRecurringTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Collect Support Data Recurring Task.",
		Attributes: map[string]schema.Attribute{
			"output_directory": schema.StringAttribute{
				Description: "The directory in which the support data archive files will be placed. The path must be a directory, and that directory must already exist. Relative paths will be interpreted as relative to the server root.",
				Required:    true,
			},
			"encryption_passphrase_file": schema.StringAttribute{
				Description: "The path to a file that contains the passphrase to encrypt the contents of the support data archive.",
				Optional:    true,
			},
			"include_expensive_data": schema.BoolAttribute{
				Description: "Indicates whether the support data archive should include information that may be expensive to obtain, and that may temporarily affect the server's performance or responsiveness.",
				Optional:    true,
				Computed:    true,
			},
			"include_replication_state_dump": schema.BoolAttribute{
				Description: "Indicates whether the support data archive should include a replication state dump, which may be several megabytes in size.",
				Optional:    true,
				Computed:    true,
			},
			"include_binary_files": schema.BoolAttribute{
				Description: "Indicates whether the support data archive should include binary files that may not have otherwise been included. Note that it may not be possible to obscure or redact sensitive information in binary files.",
				Optional:    true,
				Computed:    true,
			},
			"include_extension_source": schema.BoolAttribute{
				Description: "Indicates whether the support data archive should include the source code (if available) for any third-party extensions that may be installed in the server.",
				Optional:    true,
				Computed:    true,
			},
			"use_sequential_mode": schema.BoolAttribute{
				Description: "Indicates whether to capture support data information sequentially rather than in parallel. Capturing data in sequential mode may reduce the amount of memory that the tool requires to operate, at the cost of taking longer to run.",
				Optional:    true,
				Computed:    true,
			},
			"security_level": schema.StringAttribute{
				Description: "The security level to use when deciding which information to include in or exclude from the support data archive, and which included data should be obscured or redacted.",
				Optional:    true,
				Computed:    true,
			},
			"jstack_count": schema.Int64Attribute{
				Description: "The number of times to invoke the jstack utility to obtain a stack trace of all threads running in the JVM. A value of zero indicates that the jstack utility should not be invoked.",
				Optional:    true,
				Computed:    true,
			},
			"report_count": schema.Int64Attribute{
				Description: "The number of intervals of data to collect from tools that use sample-based reporting, like vmstat, iostat, and mpstat. A value of zero indicates that these kinds of tools should not be used to collect any information.",
				Optional:    true,
				Computed:    true,
			},
			"report_interval_seconds": schema.Int64Attribute{
				Description: "The duration (in seconds) between each interval of data to collect from tools that use sample-based reporting, like vmstat, iostat, and mpstat.",
				Optional:    true,
				Computed:    true,
			},
			"log_duration": schema.StringAttribute{
				Description: "The maximum age (leading up to the time the collect-support-data tool was invoked) for log content to include in the support data archive.",
				Optional:    true,
			},
			"log_file_head_collection_size": schema.StringAttribute{
				Description: "The amount of data to collect from the beginning of each log file included in the support data archive.",
				Optional:    true,
			},
			"log_file_tail_collection_size": schema.StringAttribute{
				Description: "The amount of data to collect from the end of each log file included in the support data archive.",
				Optional:    true,
			},
			"comment": schema.StringAttribute{
				Description: "An optional comment to include in a README file within the support data archive.",
				Optional:    true,
			},
			"retain_previous_support_data_archive_count": schema.Int64Attribute{
				Description: "The minimum number of previous support data archives that should be preserved after a new archive is generated.",
				Optional:    true,
			},
			"retain_previous_support_data_archive_age": schema.StringAttribute{
				Description: "The minimum age of previous support data archives that should be preserved after a new archive is generated.",
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
func addOptionalCollectSupportDataRecurringTaskFields(ctx context.Context, addRequest *client.AddCollectSupportDataRecurringTaskRequest, plan collectSupportDataRecurringTaskResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionPassphraseFile) {
		stringVal := plan.EncryptionPassphraseFile.ValueString()
		addRequest.EncryptionPassphraseFile = &stringVal
	}
	if internaltypes.IsDefined(plan.IncludeExpensiveData) {
		boolVal := plan.IncludeExpensiveData.ValueBool()
		addRequest.IncludeExpensiveData = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeReplicationStateDump) {
		boolVal := plan.IncludeReplicationStateDump.ValueBool()
		addRequest.IncludeReplicationStateDump = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeBinaryFiles) {
		boolVal := plan.IncludeBinaryFiles.ValueBool()
		addRequest.IncludeBinaryFiles = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeExtensionSource) {
		boolVal := plan.IncludeExtensionSource.ValueBool()
		addRequest.IncludeExtensionSource = &boolVal
	}
	if internaltypes.IsDefined(plan.UseSequentialMode) {
		boolVal := plan.UseSequentialMode.ValueBool()
		addRequest.UseSequentialMode = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SecurityLevel) {
		securityLevel, err := client.NewEnumrecurringTaskSecurityLevelPropFromValue(plan.SecurityLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.SecurityLevel = securityLevel
	}
	if internaltypes.IsDefined(plan.JstackCount) {
		intVal := int32(plan.JstackCount.ValueInt64())
		addRequest.JstackCount = &intVal
	}
	if internaltypes.IsDefined(plan.ReportCount) {
		intVal := int32(plan.ReportCount.ValueInt64())
		addRequest.ReportCount = &intVal
	}
	if internaltypes.IsDefined(plan.ReportIntervalSeconds) {
		intVal := int32(plan.ReportIntervalSeconds.ValueInt64())
		addRequest.ReportIntervalSeconds = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogDuration) {
		stringVal := plan.LogDuration.ValueString()
		addRequest.LogDuration = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFileHeadCollectionSize) {
		stringVal := plan.LogFileHeadCollectionSize.ValueString()
		addRequest.LogFileHeadCollectionSize = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFileTailCollectionSize) {
		stringVal := plan.LogFileTailCollectionSize.ValueString()
		addRequest.LogFileTailCollectionSize = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Comment) {
		stringVal := plan.Comment.ValueString()
		addRequest.Comment = &stringVal
	}
	if internaltypes.IsDefined(plan.RetainPreviousSupportDataArchiveCount) {
		intVal := int32(plan.RetainPreviousSupportDataArchiveCount.ValueInt64())
		addRequest.RetainPreviousSupportDataArchiveCount = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetainPreviousSupportDataArchiveAge) {
		stringVal := plan.RetainPreviousSupportDataArchiveAge.ValueString()
		addRequest.RetainPreviousSupportDataArchiveAge = &stringVal
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

// Read a CollectSupportDataRecurringTaskResponse object into the model struct
func readCollectSupportDataRecurringTaskResponse(ctx context.Context, r *client.CollectSupportDataRecurringTaskResponse, state *collectSupportDataRecurringTaskResourceModel, expectedValues *collectSupportDataRecurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.OutputDirectory = types.StringValue(r.OutputDirectory)
	state.EncryptionPassphraseFile = internaltypes.StringTypeOrNil(r.EncryptionPassphraseFile, internaltypes.IsEmptyString(expectedValues.EncryptionPassphraseFile))
	state.IncludeExpensiveData = internaltypes.BoolTypeOrNil(r.IncludeExpensiveData)
	state.IncludeReplicationStateDump = internaltypes.BoolTypeOrNil(r.IncludeReplicationStateDump)
	state.IncludeBinaryFiles = internaltypes.BoolTypeOrNil(r.IncludeBinaryFiles)
	state.IncludeExtensionSource = internaltypes.BoolTypeOrNil(r.IncludeExtensionSource)
	state.UseSequentialMode = internaltypes.BoolTypeOrNil(r.UseSequentialMode)
	state.SecurityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumrecurringTaskSecurityLevelProp(r.SecurityLevel), internaltypes.IsEmptyString(expectedValues.SecurityLevel))
	state.JstackCount = internaltypes.Int64TypeOrNil(r.JstackCount)
	state.ReportCount = internaltypes.Int64TypeOrNil(r.ReportCount)
	state.ReportIntervalSeconds = internaltypes.Int64TypeOrNil(r.ReportIntervalSeconds)
	state.LogDuration = internaltypes.StringTypeOrNil(r.LogDuration, internaltypes.IsEmptyString(expectedValues.LogDuration))
	config.CheckMismatchedPDFormattedAttributes("log_duration",
		expectedValues.LogDuration, state.LogDuration, diagnostics)
	state.LogFileHeadCollectionSize = internaltypes.StringTypeOrNil(r.LogFileHeadCollectionSize, internaltypes.IsEmptyString(expectedValues.LogFileHeadCollectionSize))
	config.CheckMismatchedPDFormattedAttributes("log_file_head_collection_size",
		expectedValues.LogFileHeadCollectionSize, state.LogFileHeadCollectionSize, diagnostics)
	state.LogFileTailCollectionSize = internaltypes.StringTypeOrNil(r.LogFileTailCollectionSize, internaltypes.IsEmptyString(expectedValues.LogFileTailCollectionSize))
	config.CheckMismatchedPDFormattedAttributes("log_file_tail_collection_size",
		expectedValues.LogFileTailCollectionSize, state.LogFileTailCollectionSize, diagnostics)
	state.Comment = internaltypes.StringTypeOrNil(r.Comment, internaltypes.IsEmptyString(expectedValues.Comment))
	state.RetainPreviousSupportDataArchiveCount = internaltypes.Int64TypeOrNil(r.RetainPreviousSupportDataArchiveCount)
	state.RetainPreviousSupportDataArchiveAge = internaltypes.StringTypeOrNil(r.RetainPreviousSupportDataArchiveAge, internaltypes.IsEmptyString(expectedValues.RetainPreviousSupportDataArchiveAge))
	config.CheckMismatchedPDFormattedAttributes("retain_previous_support_data_archive_age",
		expectedValues.RetainPreviousSupportDataArchiveAge, state.RetainPreviousSupportDataArchiveAge, diagnostics)
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
func createCollectSupportDataRecurringTaskOperations(plan collectSupportDataRecurringTaskResourceModel, state collectSupportDataRecurringTaskResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.OutputDirectory, state.OutputDirectory, "output-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionPassphraseFile, state.EncryptionPassphraseFile, "encryption-passphrase-file")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeExpensiveData, state.IncludeExpensiveData, "include-expensive-data")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeReplicationStateDump, state.IncludeReplicationStateDump, "include-replication-state-dump")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeBinaryFiles, state.IncludeBinaryFiles, "include-binary-files")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeExtensionSource, state.IncludeExtensionSource, "include-extension-source")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseSequentialMode, state.UseSequentialMode, "use-sequential-mode")
	operations.AddStringOperationIfNecessary(&ops, plan.SecurityLevel, state.SecurityLevel, "security-level")
	operations.AddInt64OperationIfNecessary(&ops, plan.JstackCount, state.JstackCount, "jstack-count")
	operations.AddInt64OperationIfNecessary(&ops, plan.ReportCount, state.ReportCount, "report-count")
	operations.AddInt64OperationIfNecessary(&ops, plan.ReportIntervalSeconds, state.ReportIntervalSeconds, "report-interval-seconds")
	operations.AddStringOperationIfNecessary(&ops, plan.LogDuration, state.LogDuration, "log-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFileHeadCollectionSize, state.LogFileHeadCollectionSize, "log-file-head-collection-size")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFileTailCollectionSize, state.LogFileTailCollectionSize, "log-file-tail-collection-size")
	operations.AddStringOperationIfNecessary(&ops, plan.Comment, state.Comment, "comment")
	operations.AddInt64OperationIfNecessary(&ops, plan.RetainPreviousSupportDataArchiveCount, state.RetainPreviousSupportDataArchiveCount, "retain-previous-support-data-archive-count")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainPreviousSupportDataArchiveAge, state.RetainPreviousSupportDataArchiveAge, "retain-previous-support-data-archive-age")
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
func (r *collectSupportDataRecurringTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan collectSupportDataRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddCollectSupportDataRecurringTaskRequest(plan.Id.ValueString(),
		[]client.EnumcollectSupportDataRecurringTaskSchemaUrn{client.ENUMCOLLECTSUPPORTDATARECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKCOLLECT_SUPPORT_DATA},
		plan.OutputDirectory.ValueString())
	err := addOptionalCollectSupportDataRecurringTaskFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Collect Support Data Recurring Task", err.Error())
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
		client.AddCollectSupportDataRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Collect Support Data Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state collectSupportDataRecurringTaskResourceModel
	readCollectSupportDataRecurringTaskResponse(ctx, addResponse.CollectSupportDataRecurringTaskResponse, &state, &plan, &resp.Diagnostics)

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
func (r *collectSupportDataRecurringTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state collectSupportDataRecurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RecurringTaskApi.GetRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Collect Support Data Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCollectSupportDataRecurringTaskResponse(ctx, readResponse.CollectSupportDataRecurringTaskResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *collectSupportDataRecurringTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan collectSupportDataRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state collectSupportDataRecurringTaskResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.RecurringTaskApi.UpdateRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createCollectSupportDataRecurringTaskOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.RecurringTaskApi.UpdateRecurringTaskExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Collect Support Data Recurring Task", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCollectSupportDataRecurringTaskResponse(ctx, updateResponse.CollectSupportDataRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
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
func (r *collectSupportDataRecurringTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state collectSupportDataRecurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RecurringTaskApi.DeleteRecurringTaskExecute(r.apiClient.RecurringTaskApi.DeleteRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Collect Support Data Recurring Task", err, httpResp)
		return
	}
}

func (r *collectSupportDataRecurringTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
