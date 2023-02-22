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
	_ resource.Resource                = &ldifExportRecurringTaskResource{}
	_ resource.ResourceWithConfigure   = &ldifExportRecurringTaskResource{}
	_ resource.ResourceWithImportState = &ldifExportRecurringTaskResource{}
)

// Create a Ldif Export Recurring Task resource
func NewLdifExportRecurringTaskResource() resource.Resource {
	return &ldifExportRecurringTaskResource{}
}

// ldifExportRecurringTaskResource is the resource implementation.
type ldifExportRecurringTaskResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *ldifExportRecurringTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldif_export_recurring_task"
}

// Configure adds the provider configured client to the resource.
func (r *ldifExportRecurringTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type ldifExportRecurringTaskResourceModel struct {
	Id                             types.String `tfsdk:"id"`
	LastUpdated                    types.String `tfsdk:"last_updated"`
	Notifications                  types.Set    `tfsdk:"notifications"`
	RequiredActions                types.Set    `tfsdk:"required_actions"`
	LdifDirectory                  types.String `tfsdk:"ldif_directory"`
	BackendID                      types.Set    `tfsdk:"backend_id"`
	ExcludeBackendID               types.Set    `tfsdk:"exclude_backend_id"`
	Compress                       types.Bool   `tfsdk:"compress"`
	Encrypt                        types.Bool   `tfsdk:"encrypt"`
	EncryptionSettingsDefinitionID types.String `tfsdk:"encryption_settings_definition_id"`
	Sign                           types.Bool   `tfsdk:"sign"`
	RetainPreviousLDIFExportCount  types.Int64  `tfsdk:"retain_previous_ldif_export_count"`
	RetainPreviousLDIFExportAge    types.String `tfsdk:"retain_previous_ldif_export_age"`
	MaxMegabytesPerSecond          types.Int64  `tfsdk:"max_megabytes_per_second"`
	Description                    types.String `tfsdk:"description"`
	CancelOnTaskDependencyFailure  types.Bool   `tfsdk:"cancel_on_task_dependency_failure"`
	EmailOnStart                   types.Set    `tfsdk:"email_on_start"`
	EmailOnSuccess                 types.Set    `tfsdk:"email_on_success"`
	EmailOnFailure                 types.Set    `tfsdk:"email_on_failure"`
	AlertOnStart                   types.Bool   `tfsdk:"alert_on_start"`
	AlertOnSuccess                 types.Bool   `tfsdk:"alert_on_success"`
	AlertOnFailure                 types.Bool   `tfsdk:"alert_on_failure"`
}

// GetSchema defines the schema for the resource.
func (r *ldifExportRecurringTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Ldif Export Recurring Task.",
		Attributes: map[string]schema.Attribute{
			"ldif_directory": schema.StringAttribute{
				Description: "The directory in which LDIF export files will be placed. The directory must already exist.",
				Required:    true,
			},
			"backend_id": schema.SetAttribute{
				Description: "The backend ID for a backend to be exported.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_backend_id": schema.SetAttribute{
				Description: "The backend ID for a backend to be excluded from the export.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"compress": schema.BoolAttribute{
				Description: "Indicates whether to compress the LDIF data as it is exported.",
				Optional:    true,
				Computed:    true,
			},
			"encrypt": schema.BoolAttribute{
				Description: "Indicates whether to encrypt the LDIF data as it exported.",
				Optional:    true,
				Computed:    true,
			},
			"encryption_settings_definition_id": schema.StringAttribute{
				Description: "The ID of an encryption settings definition to use to obtain the LDIF export encryption key.",
				Optional:    true,
			},
			"sign": schema.BoolAttribute{
				Description: "Indicates whether to cryptographically sign the exported data, which will make it possible to detect whether the LDIF data has been altered since it was exported.",
				Optional:    true,
				Computed:    true,
			},
			"retain_previous_ldif_export_count": schema.Int64Attribute{
				Description: "The minimum number of previous LDIF exports that should be preserved after a new export completes successfully.",
				Optional:    true,
			},
			"retain_previous_ldif_export_age": schema.StringAttribute{
				Description: "The minimum age of previous LDIF exports that should be preserved after a new export completes successfully.",
				Optional:    true,
			},
			"max_megabytes_per_second": schema.Int64Attribute{
				Description: "The maximum rate, in megabytes per second, at which LDIF exports should be written.",
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
func addOptionalLdifExportRecurringTaskFields(ctx context.Context, addRequest *client.AddLdifExportRecurringTaskRequest, plan ldifExportRecurringTaskResourceModel) {
	if internaltypes.IsDefined(plan.BackendID) {
		var slice []string
		plan.BackendID.ElementsAs(ctx, &slice, false)
		addRequest.BackendID = slice
	}
	if internaltypes.IsDefined(plan.ExcludeBackendID) {
		var slice []string
		plan.ExcludeBackendID.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeBackendID = slice
	}
	if internaltypes.IsDefined(plan.Compress) {
		boolVal := plan.Compress.ValueBool()
		addRequest.Compress = &boolVal
	}
	if internaltypes.IsDefined(plan.Encrypt) {
		boolVal := plan.Encrypt.ValueBool()
		addRequest.Encrypt = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		stringVal := plan.EncryptionSettingsDefinitionID.ValueString()
		addRequest.EncryptionSettingsDefinitionID = &stringVal
	}
	if internaltypes.IsDefined(plan.Sign) {
		boolVal := plan.Sign.ValueBool()
		addRequest.Sign = &boolVal
	}
	if internaltypes.IsDefined(plan.RetainPreviousLDIFExportCount) {
		intVal := int32(plan.RetainPreviousLDIFExportCount.ValueInt64())
		addRequest.RetainPreviousLDIFExportCount = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RetainPreviousLDIFExportAge) {
		stringVal := plan.RetainPreviousLDIFExportAge.ValueString()
		addRequest.RetainPreviousLDIFExportAge = &stringVal
	}
	if internaltypes.IsDefined(plan.MaxMegabytesPerSecond) {
		intVal := int32(plan.MaxMegabytesPerSecond.ValueInt64())
		addRequest.MaxMegabytesPerSecond = &intVal
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

// Read a LdifExportRecurringTaskResponse object into the model struct
func readLdifExportRecurringTaskResponse(ctx context.Context, r *client.LdifExportRecurringTaskResponse, state *ldifExportRecurringTaskResourceModel, expectedValues *ldifExportRecurringTaskResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.LdifDirectory = types.StringValue(r.LdifDirectory)
	state.BackendID = internaltypes.GetStringSet(r.BackendID)
	state.ExcludeBackendID = internaltypes.GetStringSet(r.ExcludeBackendID)
	state.Compress = internaltypes.BoolTypeOrNil(r.Compress)
	state.Encrypt = internaltypes.BoolTypeOrNil(r.Encrypt)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Sign = internaltypes.BoolTypeOrNil(r.Sign)
	state.RetainPreviousLDIFExportCount = internaltypes.Int64TypeOrNil(r.RetainPreviousLDIFExportCount)
	state.RetainPreviousLDIFExportAge = internaltypes.StringTypeOrNil(r.RetainPreviousLDIFExportAge, internaltypes.IsEmptyString(expectedValues.RetainPreviousLDIFExportAge))
	config.CheckMismatchedPDFormattedAttributes("retain_previous_ldif_export_age",
		expectedValues.RetainPreviousLDIFExportAge, state.RetainPreviousLDIFExportAge, diagnostics)
	state.MaxMegabytesPerSecond = internaltypes.Int64TypeOrNil(r.MaxMegabytesPerSecond)
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
func createLdifExportRecurringTaskOperations(plan ldifExportRecurringTaskResourceModel, state ldifExportRecurringTaskResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.LdifDirectory, state.LdifDirectory, "ldif-directory")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BackendID, state.BackendID, "backend-id")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeBackendID, state.ExcludeBackendID, "exclude-backend-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.Compress, state.Compress, "compress")
	operations.AddBoolOperationIfNecessary(&ops, plan.Encrypt, state.Encrypt, "encrypt")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionSettingsDefinitionID, state.EncryptionSettingsDefinitionID, "encryption-settings-definition-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.Sign, state.Sign, "sign")
	operations.AddInt64OperationIfNecessary(&ops, plan.RetainPreviousLDIFExportCount, state.RetainPreviousLDIFExportCount, "retain-previous-ldif-export-count")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainPreviousLDIFExportAge, state.RetainPreviousLDIFExportAge, "retain-previous-ldif-export-age")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxMegabytesPerSecond, state.MaxMegabytesPerSecond, "max-megabytes-per-second")
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
func (r *ldifExportRecurringTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ldifExportRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddLdifExportRecurringTaskRequest(plan.Id.ValueString(),
		[]client.EnumldifExportRecurringTaskSchemaUrn{client.ENUMLDIFEXPORTRECURRINGTASKSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0RECURRING_TASKLDIF_EXPORT},
		plan.LdifDirectory.ValueString())
	addOptionalLdifExportRecurringTaskFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RecurringTaskApi.AddRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRecurringTaskRequest(
		client.AddLdifExportRecurringTaskRequestAsAddRecurringTaskRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RecurringTaskApi.AddRecurringTaskExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Ldif Export Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state ldifExportRecurringTaskResourceModel
	readLdifExportRecurringTaskResponse(ctx, addResponse.LdifExportRecurringTaskResponse, &state, &plan, &resp.Diagnostics)

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
func (r *ldifExportRecurringTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ldifExportRecurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RecurringTaskApi.GetRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldif Export Recurring Task", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLdifExportRecurringTaskResponse(ctx, readResponse.LdifExportRecurringTaskResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *ldifExportRecurringTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ldifExportRecurringTaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state ldifExportRecurringTaskResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.RecurringTaskApi.UpdateRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createLdifExportRecurringTaskOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.RecurringTaskApi.UpdateRecurringTaskExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldif Export Recurring Task", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdifExportRecurringTaskResponse(ctx, updateResponse.LdifExportRecurringTaskResponse, &state, &plan, &resp.Diagnostics)
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
func (r *ldifExportRecurringTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ldifExportRecurringTaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RecurringTaskApi.DeleteRecurringTaskExecute(r.apiClient.RecurringTaskApi.DeleteRecurringTask(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Ldif Export Recurring Task", err, httpResp)
		return
	}
}

func (r *ldifExportRecurringTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
